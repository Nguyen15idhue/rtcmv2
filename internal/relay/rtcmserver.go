package relay

import (
	"net"
	"sync"
	"time"
)

type RTCMServer struct {
	Addr      string
	Port      int
	Frames    [][]byte
	StationID uint16
	Listener  net.Listener
	Connected int
	mu        sync.Mutex
	wg        sync.WaitGroup
	stopped   bool
}

func NewRTCMServer(port int, stationID uint16) *RTCMServer {
	return &RTCMServer{
		Port:      port,
		StationID: stationID,
		Frames:    make([][]byte, 0),
	}
}

func (s *RTCMServer) Start() error {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	s.Listener = ln
	s.Addr = ln.Addr().String()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			if s.stopped {
				return
			}
			conn, err := ln.Accept()
			if err != nil {
				if s.stopped {
					return
				}
				time.Sleep(100 * time.Millisecond)
				continue
			}
			s.mu.Lock()
			s.Connected++
			s.mu.Unlock()
			go s.handleConnection(conn)
		}
	}()

	time.Sleep(100 * time.Millisecond)
	return nil
}

func (s *RTCMServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	for _, frame := range s.Frames {
		if s.stopped {
			break
		}
		conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
		if _, err := conn.Write(frame); err != nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (s *RTCMServer) Stop() {
	s.stopped = true
	if s.Listener != nil {
		s.Listener.Close()
	}
	s.wg.Wait()
}

func (s *RTCMServer) GetConnectedCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Connected
}

func GenerateRTCMFrames(stationID uint16, count int) [][]byte {
	frames := make([][]byte, count)
	for i := 0; i < count; i++ {
		frame := make([]byte, 24)
		frame[0] = 0xD3
		frame[1] = 0x00
		frame[2] = 0x12
		frame[3] = 0x43
		frame[4] = 0x20
		frame[5] = byte(stationID & 0xFF)
		frame[6] = byte((stationID >> 8) & 0x0F)
		for j := 7; j < 21; j++ {
			frame[j] = byte(i*10 + j)
		}
		frame[21] = 0xAA
		frame[22] = 0xBB
		frame[23] = 0xCC
		frames[i] = frame
	}
	return frames
}
