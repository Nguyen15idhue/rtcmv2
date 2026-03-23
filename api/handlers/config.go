package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Nguyen15idhue/rtcmv2/internal/relay"
)

type ConfigHandler struct {
	stationsHandler *StationsHandler
	castersHandler  *CastersHandler
}

func NewConfigHandler(stations *StationsHandler, casters *CastersHandler) *ConfigHandler {
	return &ConfigHandler{
		stationsHandler: stations,
		castersHandler:  casters,
	}
}

func (h *ConfigHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"stations": h.stationsHandler.GetConfig(),
		"casters":  h.castersHandler.GetConfig(),
	})
}

func (h *ConfigHandler) ReloadConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method_not_allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	if err := h.stationsHandler.Reload(); err != nil {
		http.Error(w, `{"error":"reload_stations_failed","message":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	if err := h.castersHandler.Reload(); err != nil {
		http.Error(w, `{"error":"reload_casters_failed","message":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok","message":"config_reloaded"}`))
}

func (h *ConfigHandler) GetSystem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	global := relay.GlobalMetrics{
		ActiveStations: 0,
		TotalFrames:    0,
		TotalDrops:     0,
		UptimeSeconds:  0,
	}

	w.Write([]byte(`{`))
	w.Write([]byte(`"active_stations":` + strconv.Itoa(global.ActiveStations)))
	w.Write([]byte(`,"total_frames":` + strconv.FormatInt(global.TotalFrames, 10)))
	w.Write([]byte(`,"total_drops":` + strconv.FormatInt(global.TotalDrops, 10)))
	w.Write([]byte(`,"uptime_seconds":` + strconv.FormatInt(global.UptimeSeconds, 10)))
	w.Write([]byte(`}`))
}
