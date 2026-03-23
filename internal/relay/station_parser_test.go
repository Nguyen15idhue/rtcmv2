package relay

import (
	"testing"

	"rtcmv2/internal/buffer"
)

func TestStationParser(t *testing.T) {
	stationIDs := []uint16{1000, 2000, 3000}

	buf := buffer.New()
	parser := NewStationParser()

	t.Log("Testing Station ID extraction from RTCM frames")

	for _, stationID := range stationIDs {
		frame := generateTestMSMFrame(stationID)

		extracted := buf.Write(frame)
		if len(extracted) == 0 {
			t.Errorf("Buffer should extract frame for station %d", stationID)
			continue
		}

		parsedID, err := parser.ExtractStationID(extracted[0])

		if err != nil {
			t.Errorf("Station parser error for station %d: %v", stationID, err)
			continue
		}

		if parsedID != stationID {
			t.Errorf("Station ID mismatch: expected %d, got %d", stationID, parsedID)
		} else {
			t.Logf("Station ID %d: OK", stationID)
		}
	}
}

func TestStationParserWithDebug(t *testing.T) {
	testCases := []uint16{1000, 2000, 3000}

	for _, expected := range testCases {
		frame := generateTestMSMFrame(expected)
		payload := frame[3:]

		t.Logf("=== Station ID: %d ===", expected)
		t.Logf("frame[3:6] = %x", frame[3:6])
		t.Logf("payload[0:4] = %x", payload[0:4])
		t.Logf("payload[0] = 0x%02x, payload[1] = 0x%02x, payload[2] = 0x%02x, payload[3] = 0x%02x",
			payload[0], payload[1], payload[2], payload[3])

		manualOld := uint16(payload[3]) | (uint16(payload[4]&0x03) << 8)
		manualNew := uint16(payload[2]) | (uint16(payload[3]&0x03) << 8)
		t.Logf("Manual (payload[3]|payload[4]&0x03<<8): %d", manualOld)
		t.Logf("Manual (payload[2]|payload[3]&0x03<<8): %d", manualNew)

		parser := NewStationParser()
		result, _ := parser.ExtractStationID(frame)
		t.Logf("Parser result: %d", result)
	}
}

func generateTestMSMFrame(stationID uint16) []byte {
	frame := make([]byte, 24)
	frame[0] = 0xD3
	frame[1] = 0x00
	frame[2] = 0x12

	frame[3] = 0x43
	frame[4] = 0x20

	frame[5] = byte(stationID & 0xFF)
	frame[6] = byte((stationID >> 8) & 0x0F)

	for i := 7; i < 21; i++ {
		frame[i] = byte(i)
	}
	frame[21] = 0xAA
	frame[22] = 0xBB
	frame[23] = 0xCC

	return frame
}
