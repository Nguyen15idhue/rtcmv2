package capture

import (
	"io"
	"log"
	"runtime/debug"

	"github.com/Nguyen15idhue/rtcmv2/internal/buffer"
	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
)

type StreamFactory struct {
	frameChan chan []byte
}

func NewStreamFactory(frameChan chan []byte) *StreamFactory {
	return &StreamFactory{
		frameChan: frameChan,
	}
}

func (f *StreamFactory) New(netFlow, tcpFlow gopacket.Flow) tcpassembly.Stream {
	s := &tcpStream{
		frameChan: f.frameChan,
		buf:       buffer.New(),
	}
	go s.run()
	return &s.readerStream
}

type tcpStream struct {
	frameChan    chan []byte
	readerStream tcpreader.ReaderStream
	buf          *buffer.Buffer
}

func (t *tcpStream) run() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("tcpStream: panic recovered: %v\n%s", r, debug.Stack())
		}
	}()

	buf := make([]byte, 4096)
	for {
		n, err := t.readerStream.Read(buf)
		if n > 0 {
			frames := t.buf.Write(buf[:n])
			for _, frame := range frames {
				select {
				case t.frameChan <- frame:
				default:
					log.Printf("frame dropped: channel full")
				}
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("stream read error: %v", err)
			break
		}
	}
}

func (t *tcpStream) ReassemblyComplete() {
	t.buf.Reset()
}
