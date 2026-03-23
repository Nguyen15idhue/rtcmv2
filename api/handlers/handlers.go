package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Nguyen15idhue/rtcmv2/internal/relay"
)

type API struct {
	metrics    *relay.Metrics
	dispatcher *relay.Dispatcher
}

func NewAPI(metrics *relay.Metrics, dispatcher *relay.Dispatcher) *API {
	return &API{
		metrics:    metrics,
		dispatcher: dispatcher,
	}
}

func (h *API) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

func (h *API) Root(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/index.html")
}

func (h *API) GetStations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stations := h.metrics.GetStations()

	w.Write([]byte(`{"stations":[`))
	for i, s := range stations {
		if i > 0 {
			w.Write([]byte(","))
		}
		stationJSON := formatStation(s)
		w.Write([]byte(stationJSON))
	}
	w.Write([]byte(`],"total":`))
	writeInt(w, len(stations))
	w.Write([]byte(`}`))
}

func (h *API) GetSystem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	global := h.metrics.GetGlobal()

	w.Write([]byte(`{`))
	writeField(w, "active_stations", global.ActiveStations)
	w.Write([]byte(","))
	writeField(w, "total_frames", global.TotalFrames)
	w.Write([]byte(","))
	writeField(w, "total_drops", global.TotalDrops)
	w.Write([]byte(","))
	writeField(w, "uptime_seconds", global.UptimeSeconds)
	w.Write([]byte("}"))
}

func (h *API) StationAction(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if len(path) < 15 || path[:14] != "/api/station/" {
		http.NotFound(w, r)
		return
	}

	remaining := path[14:]

	var stationID uint64
	var action string
	var found bool

	for i, c := range remaining {
		if c == '/' {
			id, err := strconv.ParseUint(remaining[:i], 10, 16)
			if err != nil {
				http.Error(w, `{"error":"invalid_station_id"}`, 400)
				return
			}
			stationID = id
			action = remaining[i+1:]
			found = true
			break
		}
	}

	if !found {
		http.Error(w, `{"error":"invalid_path"}`, 400)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch action {
	case "reconnect":
		if h.dispatcher == nil {
			http.Error(w, `{"error":"dispatcher_not_available"}`, 500)
			return
		}
		relay := h.dispatcher.FindRelayByStationID(uint16(stationID))
		if relay == nil {
			http.Error(w, `{"error":"station_not_found"}`, 404)
			return
		}
		relay.Reconnect()
		w.Write([]byte(`{"status":"ok","action":"reconnect","station_id":` + strconv.FormatUint(stationID, 10) + `}`))

	case "disable":
		w.Write([]byte(`{"status":"ok","action":"disable","station_id":` + strconv.FormatUint(stationID, 10) + `}`))

	default:
		http.Error(w, `{"error":"unknown_action"}`, 400)
	}
}

func formatStation(s relay.StationMetrics) string {
	return `{"station_id":` + strconv.Itoa(int(s.StationID)) +
		`,"name":"` + s.Name + `"` +
		`,"frames_total":` + i64toa(s.FramesTotal) +
		`,"frames_dropped":` + i64toa(s.FramesDropped) +
		`,"last_seen":` + i64toa(s.LastSeen) +
		`,"connected":` + btoa(s.Connected) +
		`,"reconnects":` + i64toa(s.Reconnects) +
		`,"fps":` + ftoa(s.FPS) + `}`
}

func writeField(w http.ResponseWriter, name string, value interface{}) {
	w.Write([]byte(`"` + name + `":`))
	switch v := value.(type) {
	case int:
		writeInt(w, v)
	case int64:
		writeInt(w, int(v))
	case string:
		w.Write([]byte(`"` + v + `"`))
	}
}

func writeInt(w http.ResponseWriter, n int) {
	if n == 0 {
		w.Write([]byte("0"))
		return
	}
	var buf [20]byte
	i := len(buf)
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	w.Write(buf[i:])
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

func ftoa(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

func btoa(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
