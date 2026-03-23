package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/Nguyen15idhue/rtcmv2/internal/relay"
)

type CastersHandler struct {
	castersCfg *relay.CastersConfig
	mu         sync.RWMutex
}

func NewCastersHandler() *CastersHandler {
	cfg, err := relay.LoadCasters("")
	if err != nil {
		cfg = &relay.CastersConfig{
			Casters: []relay.Caster{},
		}
	}

	return &CastersHandler{
		castersCfg: cfg,
	}
}

func (h *CastersHandler) GetCasters(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")

	masked := make([]CasterResponse, len(h.castersCfg.Casters))
	for i, c := range h.castersCfg.Casters {
		masked[i] = CasterResponse{
			ID:   c.ID,
			Name: c.Name,
			Host: c.Host,
			Port: c.Port,
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"casters": masked,
		"total":   len(h.castersCfg.Casters),
	})
}

type CasterResponse struct {
	ID   uint16 `json:"id"`
	Name string `json:"name"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (h *CastersHandler) CreateCaster(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method_not_allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var c relay.Caster
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, `{"error":"invalid_json"}`, http.StatusBadRequest)
		return
	}

	if c.ID == 0 {
		http.Error(w, `{"error":"id_required"}`, http.StatusBadRequest)
		return
	}
	if c.Host == "" {
		http.Error(w, `{"error":"host_required"}`, http.StatusBadRequest)
		return
	}
	if c.Port == 0 {
		c.Port = 2101
	}
	if c.Name == "" {
		c.Name = "Caster-" + strconv.FormatUint(uint64(c.ID), 10)
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if err := relay.AddCaster(h.castersCfg, c); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusConflict)
		return
	}

	relay.SaveCasters(h.castersCfg, "")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := CasterResponse{ID: c.ID, Name: c.Name, Host: c.Host, Port: c.Port}
	json.NewEncoder(w).Encode(resp)
}

func (h *CastersHandler) DeleteCaster(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, `{"error":"method_not_allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	id := extractCasterID(r.URL.Path)
	if id == 0 {
		http.Error(w, `{"error":"invalid_caster_id"}`, http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if err := relay.DeleteCaster(h.castersCfg, id); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	relay.SaveCasters(h.castersCfg, "")

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

func (h *CastersHandler) GetConfig() *relay.CastersConfig {
	return h.castersCfg
}

func (h *CastersHandler) Reload() error {
	cfg, err := relay.LoadCasters("")
	if err != nil {
		return err
	}

	h.mu.Lock()
	h.castersCfg = cfg
	h.mu.Unlock()

	return nil
}

func extractCasterID(path string) uint16 {
	parts := splitPath(path)
	for i, p := range parts {
		if p == "caster" && i+1 < len(parts) {
			id, err := strconv.ParseUint(parts[i+1], 10, 16)
			if err == nil {
				return uint16(id)
			}
		}
	}
	return 0
}

func splitPath(path string) []string {
	var parts []string
	var current string
	for _, c := range path {
		if c == '/' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}
