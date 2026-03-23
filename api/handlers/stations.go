package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/Nguyen15idhue/rtcmv2/internal/relay"
)

type StationsHandler struct {
	stationsCfg *relay.StationsConfig
	castersCfg  *relay.CastersConfig
	metrics     *relay.Metrics
	mu          sync.RWMutex
}

func NewStationsHandler(metrics *relay.Metrics, castersCfg *relay.CastersConfig) *StationsHandler {
	cfg, err := relay.LoadStations("")
	if err != nil {
		cfg = &relay.StationsConfig{
			Stations:   []relay.Station{},
			Unassigned: []relay.UnassignedStation{},
		}
	}

	return &StationsHandler{
		stationsCfg: cfg,
		castersCfg:  castersCfg,
		metrics:     metrics,
	}
}

func (h *StationsHandler) GetStations(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")

	stations := h.stationsCfg.Stations
	metricsStations := h.metrics.GetStations()

	merged := make([]StationWithMetrics, 0, len(stations))
	for _, s := range stations {
		sm := StationWithMetrics{
			Station: s,
		}
		for _, ms := range metricsStations {
			if ms.StationID == s.ID {
				sm.FPS = ms.FPS
				sm.FramesTotal = ms.FramesTotal
				sm.FramesDropped = ms.FramesDropped
				sm.Connected = ms.Connected
				sm.LastSeen = ms.LastSeen
				break
			}
		}
		merged = append(merged, sm)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"stations":   merged,
		"unassigned": h.stationsCfg.Unassigned,
		"total":      len(stations),
	})
}

type StationWithMetrics struct {
	relay.Station
	FPS           float64 `json:"fps"`
	FramesTotal   int64   `json:"frames_total"`
	FramesDropped int64   `json:"frames_dropped"`
	Connected     bool    `json:"connected"`
	LastSeen      int64   `json:"last_seen"`
}

func (h *StationsHandler) CreateStation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method_not_allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var s relay.Station
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, `{"error":"invalid_json"}`, http.StatusBadRequest)
		return
	}

	if s.ID == 0 {
		http.Error(w, `{"error":"id_required"}`, http.StatusBadRequest)
		return
	}
	if s.Name == "" {
		s.Name = "Station-" + strconv.FormatUint(uint64(s.ID), 10)
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if err := relay.AddStation(h.stationsCfg, s); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusConflict)
		return
	}

	relay.SaveStations(h.stationsCfg, "")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

func (h *StationsHandler) UpdateStation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, `{"error":"method_not_allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	id := extractStationID(r.URL.Path)
	if id == 0 {
		http.Error(w, `{"error":"invalid_station_id"}`, http.StatusBadRequest)
		return
	}

	var s relay.Station
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, `{"error":"invalid_json"}`, http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if err := relay.UpdateStation(h.stationsCfg, id, s); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	relay.SaveStations(h.stationsCfg, "")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func (h *StationsHandler) DeleteStation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, `{"error":"method_not_allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	id := extractStationID(r.URL.Path)
	if id == 0 {
		http.Error(w, `{"error":"invalid_station_id"}`, http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if err := relay.DeleteStation(h.stationsCfg, id); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	relay.SaveStations(h.stationsCfg, "")

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

func (h *StationsHandler) AddOutput(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method_not_allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	id := extractStationID(r.URL.Path)
	if id == 0 {
		http.Error(w, `{"error":"invalid_station_id"}`, http.StatusBadRequest)
		return
	}

	var req struct {
		CasterID   uint16 `json:"caster_id"`
		Mountpoint string `json:"mountpoint"`
		Enabled    bool   `json:"enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid_json"}`, http.StatusBadRequest)
		return
	}

	if req.Mountpoint == "" {
		http.Error(w, `{"error":"mountpoint_required"}`, http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	s := relay.GetStationByID(h.stationsCfg, id)
	if s == nil {
		http.Error(w, `{"error":"station_not_found"}`, http.StatusNotFound)
		return
	}

	caster := relay.GetCasterByID(h.castersCfg, req.CasterID)
	if caster == nil {
		http.Error(w, `{"error":"caster_not_found"}`, http.StatusNotFound)
		return
	}

	output := relay.Output{
		CasterID:   req.CasterID,
		Mountpoint: req.Mountpoint,
		Enabled:    true,
	}
	if !req.Enabled {
		output.Enabled = false
	}

	for i := range h.stationsCfg.Stations {
		if h.stationsCfg.Stations[i].ID == id {
			h.stationsCfg.Stations[i].Outputs = append(h.stationsCfg.Stations[i].Outputs, output)
			break
		}
	}

	relay.SaveStations(h.stationsCfg, "")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

func (h *StationsHandler) RemoveOutput(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, `{"error":"method_not_allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 7 {
		http.Error(w, `{"error":"invalid_path"}`, http.StatusBadRequest)
		return
	}

	stationID, _ := strconv.ParseUint(parts[3], 10, 16)
	casterID, _ := strconv.ParseUint(parts[5], 10, 16)

	h.mu.Lock()
	defer h.mu.Unlock()

	for i := range h.stationsCfg.Stations {
		if uint16(h.stationsCfg.Stations[i].ID) == uint16(stationID) {
			outputs := h.stationsCfg.Stations[i].Outputs
			for j := range outputs {
				if outputs[j].CasterID == uint16(casterID) {
					h.stationsCfg.Stations[i].Outputs = append(outputs[:j], outputs[j+1:]...)
					relay.SaveStations(h.stationsCfg, "")
					w.Header().Set("Content-Type", "application/json")
					w.Write([]byte(`{"status":"ok"}`))
					return
				}
			}
			break
		}
	}

	http.Error(w, `{"error":"output_not_found"}`, http.StatusNotFound)
}

func (h *StationsHandler) ToggleOutput(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, `{"error":"method_not_allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 7 {
		http.Error(w, `{"error":"invalid_path"}`, http.StatusBadRequest)
		return
	}

	stationID, _ := strconv.ParseUint(parts[3], 10, 16)
	casterID, _ := strconv.ParseUint(parts[5], 10, 16)

	h.mu.Lock()
	defer h.mu.Unlock()

	for i := range h.stationsCfg.Stations {
		if uint16(h.stationsCfg.Stations[i].ID) == uint16(stationID) {
			outputs := h.stationsCfg.Stations[i].Outputs
			for j := range outputs {
				if outputs[j].CasterID == uint16(casterID) {
					outputs[j].Enabled = !outputs[j].Enabled
					relay.SaveStations(h.stationsCfg, "")
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(outputs[j])
					return
				}
			}
			break
		}
	}

	http.Error(w, `{"error":"output_not_found"}`, http.StatusNotFound)
}

func (h *StationsHandler) GetConfig() *relay.StationsConfig {
	return h.stationsCfg
}

func (h *StationsHandler) Reload() error {
	cfg, err := relay.LoadStations("")
	if err != nil {
		return err
	}

	h.mu.Lock()
	h.stationsCfg = cfg
	h.mu.Unlock()

	return nil
}

func extractStationID(path string) uint16 {
	parts := strings.Split(path, "/")
	for i, p := range parts {
		if p == "station" && i+1 < len(parts) {
			id, err := strconv.ParseUint(parts[i+1], 10, 16)
			if err == nil {
				return uint16(id)
			}
		}
	}
	return 0
}
