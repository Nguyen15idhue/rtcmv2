package buffer

const (
	syncByte     = 0xD3
	headerSize   = 3 // sync + 2 length bytes
	crcSize      = 3
	minFrameSize = headerSize + crcSize // 6 bytes minimum
)

type Buffer struct {
	data []byte
}

func New() *Buffer {
	return &Buffer{}
}

func (b *Buffer) Write(newData []byte) [][]byte {
	b.data = append(b.data, newData...)
	var frames [][]byte

	for {
		// a. Tìm index của 0xD3
		syncIdx := -1
		for i := 0; i < len(b.data); i++ {
			if b.data[i] == syncByte {
				syncIdx = i
				break
			}
		}

		// Không tìm thấy sync byte
		if syncIdx == -1 {
			// Giữ lại 1 byte cuối để xử lý trường hợp sync bị chia cắt
			if len(b.data) > 1 {
				b.data = b.data[len(b.data)-1:]
			}
			break
		}

		// b. Nếu index > 0, discard data trước sync byte
		if syncIdx > 0 {
			b.data = b.data[syncIdx:]
		}

		// c. Buffer < 3 byte thì break
		if len(b.data) < headerSize {
			break
		}

		// d. Đọc length (10 bits: big-endian)
		// byte[1] bits 0-1 = high 2 bits, byte[2] = low 8 bits
		length := int(b.data[1]&0x03)<<8 | int(b.data[2])

		// e. Tính frameSize = 1 (sync) + 2 (length) + length + 3 (CRC)
		frameSize := headerSize + length + crcSize

		// f. Nếu chưa đủ frame thì break
		if len(b.data) < frameSize {
			break
		}

		// g. Extract frame
		frame := make([]byte, frameSize)
		copy(frame, b.data[:frameSize])
		frames = append(frames, frame)

		// h. Remove frame khỏi buffer
		b.data = b.data[frameSize:]
	}

	return frames
}

func (b *Buffer) Reset() {
	b.data = b.data[:0]
}

func (b *Buffer) Len() int {
	return len(b.data)
}
