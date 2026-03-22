package buffer

import (
	"reflect"
	"testing"
)

func TestWrite_SingleCompleteFrame(t *testing.T) {
	// Frame: sync(0xD3) + length(3 bytes) + payload(3 bytes) + crc(3 bytes) = 10 bytes
	// Length is 10-bit big-endian: byte[1] bits 0-1 = high 2 bits, byte[2] = low 8 bits
	frame := []byte{0xD3, 0x00, 0x03, 0x11, 0x22, 0x33, 0xAA, 0xBB, 0xCC}
	buf := New()

	frames := buf.Write(frame)

	if len(frames) != 1 {
		t.Fatalf("expected 1 frame, got %d", len(frames))
	}
	if !reflect.DeepEqual(frames[0], frame) {
		t.Errorf("frame mismatch\ngot:  %x\nwant: %x", frames[0], frame)
	}
}

func TestWrite_TwoFramesConcatenated(t *testing.T) {
	frame1 := []byte{0xD3, 0x00, 0x03, 0x11, 0x22, 0x33, 0xAA, 0xBB, 0xCC}
	frame2 := []byte{0xD3, 0x00, 0x02, 0x44, 0x55, 0xDD, 0xEE, 0xFF}
	buf := New()

	frames := buf.Write(append(frame1, frame2...))

	if len(frames) != 2 {
		t.Fatalf("expected 2 frames, got %d", len(frames))
	}
	if !reflect.DeepEqual(frames[0], frame1) {
		t.Errorf("frame1 mismatch\ngot:  %x\nwant: %x", frames[0], frame1)
	}
	if !reflect.DeepEqual(frames[1], frame2) {
		t.Errorf("frame2 mismatch\ngot:  %x\nwant: %x", frames[1], frame2)
	}
}

func TestWrite_GarbageBeforeSync(t *testing.T) {
	frame := []byte{0xD3, 0x00, 0x02, 0x11, 0x22, 0xAA, 0xBB, 0xCC}
	garbage := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE}
	buf := New()

	frames := buf.Write(append(garbage, frame...))

	if len(frames) != 1 {
		t.Fatalf("expected 1 frame, got %d", len(frames))
	}
	if !reflect.DeepEqual(frames[0], frame) {
		t.Errorf("frame mismatch\ngot:  %x\nwant: %x", frames[0], frame)
	}
}

func TestWrite_IncompleteFrame(t *testing.T) {
	// Chỉ có sync byte + 1 byte length
	data := []byte{0xD3, 0x00}
	buf := New()

	frames := buf.Write(data)

	if len(frames) != 0 {
		t.Fatalf("expected 0 frames, got %d", len(frames))
	}
	if buf.Len() != 2 {
		t.Errorf("expected buffer len 2, got %d", buf.Len())
	}
}

func TestWrite_NoSyncByte(t *testing.T) {
	data := []byte{0x00, 0x01, 0x02, 0x03, 0x04}
	buf := New()

	frames := buf.Write(data)

	if len(frames) != 0 {
		t.Fatalf("expected 0 frames, got %d", len(frames))
	}
	// Giữ lại 1 byte cuối
	if buf.Len() != 1 {
		t.Errorf("expected buffer len 1, got %d", buf.Len())
	}
}

func TestWrite_EmptyInput(t *testing.T) {
	buf := New()
	buf.data = []byte{0xD3, 0x00}

	frames := buf.Write([]byte{})

	if len(frames) != 0 {
		t.Fatalf("expected 0 frames, got %d", len(frames))
	}
}

func TestWrite_OnlySyncByte(t *testing.T) {
	data := []byte{0xD3}
	buf := New()

	frames := buf.Write(data)

	if len(frames) != 0 {
		t.Fatalf("expected 0 frames, got %d", len(frames))
	}
	if buf.Len() != 1 {
		t.Errorf("expected buffer len 1, got %d", buf.Len())
	}
}

func TestWrite_PartialInMultipleWrites(t *testing.T) {
	frame := []byte{0xD3, 0x00, 0x03, 0x11, 0x22, 0x33, 0xAA, 0xBB, 0xCC}
	buf := New()

	// Write first 5 bytes (sync + length + 2 payload)
	frames1 := buf.Write(frame[:5])
	if len(frames1) != 0 {
		t.Errorf("expected 0 frames in first write, got %d", len(frames1))
	}

	// Write remaining bytes
	frames2 := buf.Write(frame[5:])
	if len(frames2) != 1 {
		t.Errorf("expected 1 frame in second write, got %d", len(frames2))
	}
	if !reflect.DeepEqual(frames2[0], frame) {
		t.Errorf("frame mismatch\ngot:  %x\nwant: %x", frames2[0], frame)
	}
}

func TestWrite_LengthParsing(t *testing.T) {
	// length = 0x0002 = 2 bytes
	frame := []byte{0xD3, 0x00, 0x02, 0x11, 0x22, 0xAA, 0xBB, 0xCC}
	buf := New()

	frames := buf.Write(frame)

	if len(frames) != 1 {
		t.Fatalf("expected 1 frame, got %d", len(frames))
	}
	// Frame size = 6 + 2 = 8 bytes
	if len(frames[0]) != 8 {
		t.Errorf("expected frame len 8, got %d", len(frames[0]))
	}
}

func TestWrite_Reset(t *testing.T) {
	buf := New()
	buf.data = []byte{0xD3, 0x00, 0x02, 0x11, 0x22}

	buf.Reset()

	if buf.Len() != 0 {
		t.Errorf("expected buffer len 0 after reset, got %d", buf.Len())
	}
}
