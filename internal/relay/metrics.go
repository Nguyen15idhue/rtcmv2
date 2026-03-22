package relay

import (
	"sync"
	"sync/atomic"
	"time"
)

type StationMetrics struct {
	StationID     uint16  `json:"station_id"`
	Name          string  `json:"name"`
	FramesTotal   int64   `json:"frames_total"`
	FramesDropped int64   `json:"frames_dropped"`
	LastSeen      int64   `json:"last_seen"`
	Connected     bool    `json:"connected"`
	Reconnects    int64   `json:"reconnects"`
	FPS           float64 `json:"fps"`
}

type GlobalMetrics struct {
	ActiveStations int   `json:"active_stations"`
	TotalFrames    int64 `json:"total_frames"`
	TotalDrops     int64 `json:"total_drops"`
	ReconnectTotal int64 `json:"reconnect_total"`
	UptimeSeconds  int64 `json:"uptime_seconds"`
}

type Metrics struct {
	mu          sync.RWMutex
	stations    map[uint16]*stationMetrics
	totalFrames atomic.Int64
	totalDrops  atomic.Int64
	startTime   time.Time
}

type stationMetrics struct {
	mu          sync.RWMutex
	StationID   uint16
	Name        string
	frames      atomic.Int64
	dropped     atomic.Int64
	lastSeen    atomic.Int64
	connected   atomic.Bool
	reconnects  atomic.Int64
	fpsMu       sync.Mutex
	fpsWindow   []int64
	fpsIndex    int
	fpsLastCalc atomic.Int64
}

const fpsWindowSize = 60

func NewMetrics() *Metrics {
	return &Metrics{
		stations:  make(map[uint16]*stationMetrics),
		startTime: time.Now(),
	}
}

func (m *Metrics) getOrCreateStation(stationID uint16, name string) *stationMetrics {
	m.mu.Lock()
	defer m.mu.Unlock()

	if s, ok := m.stations[stationID]; ok {
		return s
	}

	s := &stationMetrics{
		StationID: stationID,
		Name:      name,
		fpsWindow: make([]int64, fpsWindowSize),
	}
	m.stations[stationID] = s
	return s
}

func (m *Metrics) RecordFrame(stationID uint16, name string) {
	s := m.getOrCreateStation(stationID, name)
	s.recordFrame()
	m.totalFrames.Add(1)
}

func (s *stationMetrics) recordFrame() {
	now := time.Now().UnixNano()

	s.fpsMu.Lock()
	s.fpsWindow[s.fpsIndex] = now
	s.fpsIndex = (s.fpsIndex + 1) % fpsWindowSize
	s.fpsMu.Unlock()

	s.frames.Add(1)
	s.lastSeen.Store(now)
	s.connected.Store(true)
}

func (s *stationMetrics) RecordDrop() {
	s.dropped.Add(1)
}

func (s *stationMetrics) RecordReconnect() {
	s.reconnects.Add(1)
}

func (s *stationMetrics) SetConnected(connected bool) {
	s.connected.Store(connected)
}

func (s *stationMetrics) GetFPS() float64 {
	s.fpsMu.Lock()
	defer s.fpsMu.Unlock()

	cutoff := time.Now().Add(-5 * time.Second).UnixNano()
	count := 0
	for _, ts := range s.fpsWindow {
		if ts > cutoff {
			count++
		}
	}
	return float64(count) / 5.0
}

func (m *Metrics) GetStations() []StationMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]StationMetrics, 0, len(m.stations))
	for _, s := range m.stations {
		result = append(result, StationMetrics{
			StationID:     s.StationID,
			Name:          s.Name,
			FramesTotal:   s.frames.Load(),
			FramesDropped: s.dropped.Load(),
			LastSeen:      s.lastSeen.Load(),
			Connected:     s.connected.Load(),
			Reconnects:    s.reconnects.Load(),
			FPS:           s.GetFPS(),
		})
	}
	return result
}

func (m *Metrics) GetGlobal() GlobalMetrics {
	return GlobalMetrics{
		ActiveStations: len(m.stations),
		TotalFrames:    m.totalFrames.Load(),
		TotalDrops:     m.totalDrops.Load(),
		ReconnectTotal: 0,
		UptimeSeconds:  int64(time.Since(m.startTime).Seconds()),
	}
}

func (m *Metrics) RemoveStation(stationID uint16) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.stations, stationID)
}
