package api

import (
	"net/http"
	"strconv"

	"github.com/Nguyen15idhue/rtcmv2/api/handlers"
	"github.com/Nguyen15idhue/rtcmv2/internal/relay"
)

type Server struct {
	httpServer *http.Server
	api        *handlers.API
	mux        *http.ServeMux
}

func NewServer(addr string, metrics *relay.Metrics, dispatcher *relay.Dispatcher) *Server {
	api := handlers.NewAPI(metrics, dispatcher)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", api.Health)
	mux.HandleFunc("/api/stations", api.GetStations)
	mux.HandleFunc("/api/system", api.GetSystem)
	mux.HandleFunc("/api/station/", api.StationAction)
	mux.HandleFunc("/", api.Root)

	return &Server{
		mux: mux,
		httpServer: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
		api: api,
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() error {
	return s.httpServer.Close()
}

func stationIDFromPath(path string) (uint16, string, bool) {
	if len(path) < 16 {
		return 0, "", false
	}

	parts := path[14:]
	for i, c := range parts {
		if c == '/' {
			id, err := strconv.ParseUint(parts[:i], 10, 16)
			if err != nil {
				return 0, "", false
			}
			action := parts[i+1:]
			return uint16(id), action, true
		}
	}
	return 0, "", false
}
