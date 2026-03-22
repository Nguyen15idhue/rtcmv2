package debug

import (
	"encoding/json"
	"net/http"
	"rtcmv2/internal/relay"
)

type Server struct {
	metrics *relay.Metrics
	addr    string
}

func NewServer(metrics *relay.Metrics, addr string) *Server {
	return &Server{
		metrics: metrics,
		addr:    addr,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.healthz)
	mux.HandleFunc("/debug/stations", s.debugStations)
	mux.HandleFunc("/debug/metrics", s.debugMetrics)

	return http.ListenAndServe(s.addr, mux)
}

func (s *Server) healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *Server) debugStations(w http.ResponseWriter, r *http.Request) {
	stations := s.metrics.GetStations()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stations)
}

func (s *Server) debugMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := s.metrics.GetGlobal()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}
