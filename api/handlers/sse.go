package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Nguyen15idhue/rtcmv2/internal/relay"
)

type SSEHandler struct {
	broadcaster *Broadcaster
	metrics     *relay.Metrics
}

type Broadcaster struct {
	clients map[chan []byte]bool
	mutex   sync.RWMutex
	metrics *relay.Metrics
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewBroadcaster(metrics *relay.Metrics) *Broadcaster {
	ctx, cancel := context.WithCancel(context.Background())
	return &Broadcaster{
		clients: make(map[chan []byte]bool),
		metrics: metrics,
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (b *Broadcaster) Start() {
	go b.broadcastLoop()
}

func (b *Broadcaster) Stop() {
	b.cancel()
	b.mutex.Lock()
	defer b.mutex.Unlock()
	for ch := range b.clients {
		close(ch)
		delete(b.clients, ch)
	}
}

func (b *Broadcaster) AddClient() chan []byte {
	ch := make(chan []byte, 100)
	b.mutex.Lock()
	b.clients[ch] = true
	b.mutex.Unlock()
	return ch
}

func (b *Broadcaster) RemoveClient(ch chan []byte) {
	b.mutex.Lock()
	if _, ok := b.clients[ch]; ok {
		close(ch)
		delete(b.clients, ch)
	}
	b.mutex.Unlock()
}

func (b *Broadcaster) Broadcast(data []byte) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	for ch := range b.clients {
		select {
		case ch <- data:
		default:
			go b.RemoveClient(ch)
		}
	}
}

func (b *Broadcaster) broadcastLoop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-b.ctx.Done():
			return
		case <-ticker.C:
			b.sendMetrics()
		}
	}
}

func (b *Broadcaster) sendMetrics() {
	stations := b.metrics.GetStations()
	global := b.metrics.GetGlobal()

	data := map[string]interface{}{
		"stations": stations,
		"system": map[string]interface{}{
			"active_stations": global.ActiveStations,
			"total_frames":    global.TotalFrames,
			"total_drops":     global.TotalDrops,
			"uptime_seconds":  global.UptimeSeconds,
		},
		"timestamp": time.Now().Unix(),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	b.Broadcast(jsonData)
}

func NewSSEHandler(metrics *relay.Metrics) *SSEHandler {
	broadcaster := NewBroadcaster(metrics)
	broadcaster.Start()

	return &SSEHandler{
		broadcaster: broadcaster,
		metrics:     metrics,
	}
}

func (h *SSEHandler) StreamHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ch := h.broadcaster.AddClient()
	defer h.broadcaster.RemoveClient(ch)

	clientGone := r.Context().Done()
	for {
		select {
		case <-clientGone:
			return
		case data, ok := <-ch:
			if !ok {
				return
			}
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}

func (h *SSEHandler) Stop() {
	h.broadcaster.Stop()
}

func (h *SSEHandler) NotifyNewStation(stationID uint16) {
	data := map[string]interface{}{
		"type": "new_station",
		"data": map[string]interface{}{
			"station_id": stationID,
		},
	}
	jsonData, _ := json.Marshal(data)
	h.broadcaster.Broadcast(jsonData)
}
