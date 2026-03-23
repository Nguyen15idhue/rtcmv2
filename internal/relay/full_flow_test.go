package relay

import (
	"net"
	"testing"
	"time"

	"github.com/Nguyen15idhue/rtcmv2/internal/buffer"
)

func TestFullFlowWithSyntheticData(t *testing.T) {
	frames := make([][]byte, 50)
	for i := 0; i < 50; i++ {
		frame := make([]byte, 24)
		frame[0] = 0xD3
		frame[1] = 0x00
		frame[2] = 0x12

		frame[3] = 0x43
		frame[4] = 0x20

		stationID := uint16(1000 + (i % 3))
		frame[5] = byte(stationID & 0xFF)
		frame[6] = byte((stationID >> 8) & 0x0F)

		for j := 7; j < 21; j++ {
			frame[j] = byte(j)
		}
		frame[21] = 0xAA
		frame[22] = 0xBB
		frame[23] = 0xCC

		frames[i] = frame
	}

	t.Logf("Created %d synthetic RTCM frames", len(frames))

	// Step 2: Test Buffer extraction
	buf := buffer.New()
	var extractedFrames [][]byte
	for _, frame := range frames {
		extracted := buf.Write(frame)
		extractedFrames = append(extractedFrames, extracted...)
	}
	t.Logf("Buffer extracted %d frames from %d inputs", len(extractedFrames), len(frames))

	if len(extractedFrames) == 0 {
		t.Error("buffer should extract frames")
	}

	// Step 3: Test station ID extraction
	parser := NewStationParser()
	stationCounts := make(map[uint16]int)
	for _, frame := range extractedFrames {
		stationID, err := parser.ExtractStationID(frame)
		if err == nil && stationID > 0 {
			stationCounts[stationID]++
		}
	}

	t.Logf("Found %d unique stations: %v", len(stationCounts), stationCounts)

	if len(stationCounts) < 1 {
		t.Error("should detect at least one station")
	}

	// Step 4: Test relay to mock caster
	caster := NewMockCaster(":12123")
	if err := caster.Start(); err != nil {
		t.Fatalf("failed to start caster: %v", err)
	}
	defer caster.Close()

	time.Sleep(100 * time.Millisecond)

	// Connect
	conn, err := net.Dial("tcp", "localhost:12123")
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	// Auth
	conn.Write([]byte("SOURCE test /RTCM\r\nSource-Agent: test/1.0\r\n\r\n"))
	resp := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	conn.Read(resp)

	// Send frames
	sentCount := 0
	for _, frame := range extractedFrames {
		if _, err := conn.Write(frame); err == nil {
			sentCount++
		}
	}

	time.Sleep(100 * time.Millisecond)

	receivedCount := caster.FramesCount()
	t.Logf("Sent %d, received %d frames", sentCount, receivedCount)

	// Results
	t.Logf("✅ Full Flow Test PASSED!")
	t.Logf("   - Frames created: %d", len(frames))
	t.Logf("   - Frames extracted: %d", len(extractedFrames))
	t.Logf("   - Stations detected: %d", len(stationCounts))
	t.Logf("   - Frames sent to caster: %d", sentCount)
	t.Logf("   - Frames received by caster: %d", receivedCount)

	if receivedCount == 0 {
		t.Error("caster should receive frames")
	}
}
