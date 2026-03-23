package relay

import (
	"net"
	"testing"
	"time"
)

func TestE2ERelayDirect(t *testing.T) {
	caster := NewMockCaster(":12115")
	if err := caster.Start(); err != nil {
		t.Fatalf("failed to start caster: %v", err)
	}
	defer caster.Close()

	time.Sleep(100 * time.Millisecond)

	// Connect to caster directly
	conn, err := net.Dial("tcp", "localhost:12115")
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	// Send SOURCE command
	req := "SOURCE testpass /TEST\r\nSource-Agent: test/1.0\r\n\r\n"
	if _, err := conn.Write([]byte(req)); err != nil {
		t.Fatalf("failed to write: %v", err)
	}

	resp := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(resp)
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}

	response := string(resp[:n])
	if !responseContains(response, "ICY 200") {
		t.Errorf("expected ICY 200, got: %s", response)
	}

	t.Log("Connected to caster successfully")

	// Send frames
	for i := 0; i < 10; i++ {
		frame := []byte{0xD3, 0x00, 0x03, 0x11, 0x22, 0x33, 0xAA, 0xBB, 0xCC}
		if _, err := conn.Write(frame); err != nil {
			t.Fatalf("failed to write frame: %v", err)
		}
		time.Sleep(10 * time.Millisecond)
	}

	time.Sleep(100 * time.Millisecond)

	frameCount := caster.FramesCount()
	t.Logf("Received %d frames at caster", frameCount)

	if frameCount == 0 {
		t.Error("expected frames received, got 0")
	}
}

func responseContains(resp, substr string) bool {
	for i := 0; i <= len(resp)-len(substr); i++ {
		if resp[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
