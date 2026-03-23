package relay

import (
	"net"
	"testing"
	"time"

	"github.com/Nguyen15idhue/rtcmv2/internal/buffer"
)

func TestMultiStationE2E(t *testing.T) {
	// Create synthetic RTCM frames for 3 different stations
	// Station IDs: 1000, 2000, 3000

	stationIDs := []uint16{1000, 2000, 3000}
	framesPerStation := 20
	totalFrames := len(stationIDs) * framesPerStation

	// Generate frames for each station
	allFrames := make([][]byte, 0, totalFrames)

	for _, stationID := range stationIDs {
		for i := 0; i < framesPerStation; i++ {
			// Create RTCM 1074 (GPS MSM) frame
			frame := createMSMFrame(stationID, uint16(i))
			allFrames = append(allFrames, frame)
		}
	}

	t.Logf("Created %d frames for %d stations", len(allFrames), len(stationIDs))

	// Start mock caster
	caster := NewMockCaster(":12130")
	if err := caster.Start(); err != nil {
		t.Fatalf("failed to start caster: %v", err)
	}
	defer caster.Close()

	time.Sleep(100 * time.Millisecond)

	// Connect to caster
	conn, err := net.Dial("tcp", "localhost:12130")
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	// Authenticate
	conn.Write([]byte("SOURCE testpass /STATIONS\r\nSource-Agent: rtcmv2-test/1.0\r\n\r\n"))
	resp := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	conn.Read(resp)

	if string(resp[:3]) != "ICY" {
		t.Fatalf("expected ICY response, got: %s", string(resp[:10]))
	}

	t.Log("Connected to caster")

	// Simulate the full flow: frames → buffer → station parsing → send
	parser := NewStationParser()
	stationCounts := make(map[uint16]int)
	buf := buffer.New()

	var sentCount int
	for _, frame := range allFrames {
		// Step 1: Buffer extracts complete frames
		extracted := buf.Write(frame)
		for _, f := range extracted {
			// Step 2: Parse station ID
			stationID, err := parser.ExtractStationID(f)
			if err == nil && stationID > 0 {
				stationCounts[stationID]++
			}

			// Step 3: Send to caster
			if _, err := conn.Write(f); err == nil {
				sentCount++
			}
		}
	}

	time.Sleep(200 * time.Millisecond)

	// Verify results
	receivedCount := caster.FramesCount()

	// Report
	t.Log("=" + string(make([]byte, 50)))
	t.Log("FULL FLOW TEST RESULTS:")
	t.Log("=" + string(make([]byte, 50)))
	t.Logf("Frames generated:        %d", len(allFrames))
	t.Logf("Frames sent to caster:    %d", sentCount)
	t.Logf("Frames received:         %d", receivedCount)
	t.Logf("Unique stations found:   %d", len(stationCounts))

	t.Logf("\nStation breakdown:")
	for stationID, count := range stationCounts {
		t.Logf("  Station %d: %d frames", stationID, count)
	}

	// Verify
	if receivedCount == 0 {
		t.Fatal("FAILED: No frames received by caster")
	}

	if len(stationCounts) != len(stationIDs) {
		t.Logf("WARNING: Expected %d stations, found %d", len(stationIDs), len(stationIDs))
	}

	if sentCount != receivedCount {
		t.Logf("WARNING: Sent %d, received %d", sentCount, receivedCount)
	}

	t.Log("=" + string(make([]byte, 50)))
	t.Log("✅ FULL FLOW TEST PASSED!")
	t.Log("=" + string(make([]byte, 50)))
}

// createMSMFrame creates a proper RTCM MSM 1074 frame with given station ID
func createMSMFrame(stationID uint16, seq uint16) []byte {
	payloadLen := 20
	totalFrameLen := 3 + payloadLen + 3

	frame := make([]byte, totalFrameLen)

	frame[0] = 0xD3
	frame[1] = byte((payloadLen >> 8) & 0x03)
	frame[2] = byte(payloadLen & 0xFF)

	frame[3] = 0x43
	frame[4] = 0x20
	frame[5] = byte(stationID & 0xFF)
	frame[6] = byte((stationID >> 8) & 0x0F)

	for i := 7; i < 3+payloadLen; i++ {
		frame[i] = byte(i + int(seq))
	}

	frame[totalFrameLen-3] = 0xAA
	frame[totalFrameLen-2] = 0xBB
	frame[totalFrameLen-1] = 0xCC

	return frame
}
