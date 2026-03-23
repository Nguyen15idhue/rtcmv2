package relay

import (
	"bufio"
	"io"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

type MockCaster struct {
	Listener  net.Listener
	Port      string
	Frames    [][]byte
	mu        sync.Mutex
	connected bool
	OnConnect func()
	closed    bool
	wg        sync.WaitGroup
}

func NewMockCaster(port string) *MockCaster {
	return &MockCaster{
		Port:   port,
		Frames: make([][]byte, 0),
	}
}

func (m *MockCaster) Start() error {
	ln, err := net.Listen("tcp", m.Port)
	if err != nil {
		return err
	}
	m.Listener = ln

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		for {
			if m.closed {
				return
			}
			conn, err := m.Listener.Accept()
			if err != nil {
				if m.closed {
					return
				}
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					time.Sleep(100 * time.Millisecond)
					continue
				}
				return
			}
			go m.handleConnection(conn)
		}
	}()

	time.Sleep(100 * time.Millisecond)
	return nil
}

func (m *MockCaster) handleConnection(conn net.Conn) {
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(10 * time.Second))
	reader := bufio.NewReader(conn)

	req, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return
	}

	if strings.HasPrefix(req, "SOURCE ") {
		conn.Write([]byte("ICY 200 OK\r\n"))
		m.mu.Lock()
		m.connected = true
		if m.OnConnect != nil {
			m.OnConnect()
		}
		m.mu.Unlock()

		buf := make([]byte, 4096)
		for {
			conn.SetDeadline(time.Now().Add(1 * time.Second))
			n, err := conn.Read(buf)
			if err != nil {
				break
			}
			if n > 0 {
				m.mu.Lock()
				m.Frames = append(m.Frames, make([]byte, n))
				copy(m.Frames[len(m.Frames)-1], buf[:n])
				m.mu.Unlock()
			}
		}
	}
}

func (m *MockCaster) Close() {
	m.closed = true
	if m.Listener != nil {
		m.Listener.Close()
	}
	m.wg.Wait()
}

func (m *MockCaster) FramesCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.Frames)
}

func (m *MockCaster) Reset() {
	m.mu.Lock()
	m.Frames = nil
	m.connected = false
	m.mu.Unlock()
}

func TestMockCaster(t *testing.T) {
	caster := NewMockCaster(":12111")
	if err := caster.Start(); err != nil {
		t.Fatalf("failed to start caster: %v", err)
	}
	defer caster.Close()

	conn, err := net.Dial("tcp", "localhost:12111")
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

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
	if !strings.HasPrefix(response, "ICY 200") {
		t.Errorf("expected ICY 200, got: %s", response)
	}

	time.Sleep(100 * time.Millisecond)

	frame := []byte{0xD3, 0x00, 0x03, 0x11, 0x22, 0x33, 0xAA, 0xBB, 0xCC}
	conn.Write(frame)

	time.Sleep(100 * time.Millisecond)

	if caster.FramesCount() == 0 {
		t.Error("expected frames received, got 0")
	}

	t.Logf("Mock caster received %d frames", caster.FramesCount())
}
