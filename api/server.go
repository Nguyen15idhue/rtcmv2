package api

import (
	"net/http"
	"strings"

	"github.com/Nguyen15idhue/rtcmv2/api/handlers"
	"github.com/Nguyen15idhue/rtcmv2/api/middleware"
	"github.com/Nguyen15idhue/rtcmv2/internal/relay"
)

type Server struct {
	httpServer      *http.Server
	stationsHandler *handlers.StationsHandler
	castersHandler  *handlers.CastersHandler
	configHandler   *handlers.ConfigHandler
	sseHandler      *handlers.SSEHandler
	authMiddleware  *middleware.AuthMiddleware
	dispatcher      *relay.Dispatcher
	metrics         *relay.Metrics
}

func NewServer(addr string, metrics *relay.Metrics, dispatcher *relay.Dispatcher) *Server {
	auth := middleware.NewAuthMiddleware()
	casters := handlers.NewCastersHandler()
	stations := handlers.NewStationsHandler(metrics, casters.GetConfig())
	config := handlers.NewConfigHandler(stations, casters)
	sse := handlers.NewSSEHandler(metrics)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("/api/stream", sse.StreamHandler)
	mux.HandleFunc("/api/system", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		global := metrics.GetGlobal()
		w.Write([]byte(`{"active_stations":` + itoa(global.ActiveStations) +
			`,"total_frames":` + i64toa(global.TotalFrames) +
			`,"total_drops":` + i64toa(global.TotalDrops) +
			`,"uptime_seconds":` + i64toa(global.UptimeSeconds) + `}`))
	})

	mux.HandleFunc("/api/stations", stations.GetStations)
	mux.HandleFunc("/api/station", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			stations.CreateStation(w, r)
		default:
			http.Error(w, `{"error":"method_not_allowed"}`, http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/station/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch r.Method {
		case http.MethodPost:
			if strings.HasSuffix(path, "/output") {
				stations.AddOutput(w, r)
				return
			}
			http.Error(w, `{"error":"not_found"}`, http.StatusNotFound)
		case http.MethodPut:
			if strings.HasSuffix(path, "/toggle") {
				stations.ToggleOutput(w, r)
				return
			}
			stations.UpdateStation(w, r)
		case http.MethodDelete:
			stations.DeleteStation(w, r)
		case http.MethodGet:
			stations.GetStations(w, r)
		default:
			http.Error(w, `{"error":"method_not_allowed"}`, http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/casters", casters.GetCasters)
	mux.HandleFunc("/api/caster", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			casters.CreateCaster(w, r)
		default:
			http.Error(w, `{"error":"method_not_allowed"}`, http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/caster/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			casters.DeleteCaster(w, r)
		default:
			http.Error(w, `{"error":"method_not_allowed"}`, http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/config", config.GetConfig)
	mux.HandleFunc("/api/reload", config.ReloadConfig)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || !strings.HasPrefix(r.URL.Path, "/api/") {
			http.ServeFile(w, r, "public/index.html")
			return
		}
		http.NotFound(w, r)
	})

	handler := auth.RequireAuth(mux)

	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
		stationsHandler: stations,
		castersHandler:  casters,
		configHandler:   config,
		sseHandler:      sse,
		authMiddleware:  auth,
		dispatcher:      dispatcher,
		metrics:         metrics,
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() error {
	s.sseHandler.Stop()
	return s.httpServer.Close()
}

func (s *Server) GetStationsHandler() *handlers.StationsHandler {
	return s.stationsHandler
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}

func i64toa(n int64) string {
	if n == 0 {
		return "0"
	}
	var buf [24]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
