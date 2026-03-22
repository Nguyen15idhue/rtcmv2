# Current Session

## Task
Phase 9: Observability - COMPLETED

## Status
- [x] Phase 2: Buffer (RTCM frame extraction)
- [x] Phase 3: TCP Capture + Reassembly
- [x] Phase 4a: Basic NTRIP Relay
- [x] Phase 4b: RTCM Parsing + Station Routing
- [x] Phase 5: Frame Dispatcher (Dynamic Station Discovery)
- [x] Phase 9: Observability
  - [x] 9.1 Logging (slog structured)
  - [x] 9.2 Metrics (in-memory)
  - [x] 9.3 Debug HTTP Server
  - [x] 9.4 FPS Calculation (sliding window)

## Architecture
```
FrameChan → Dispatcher → Relay → Caster
                    ↓
              Metrics (in-memory)
                    ↓
              Debug HTTP Server (:8080)
```

## Files

### internal/relay/logger.go (NEW)
```go
var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func logInfo(msg string, fields LogFields)
func logError(msg string, fields LogFields)
func logWarn(msg string, fields LogFields)
```

### internal/relay/metrics.go (NEW)
```go
type Metrics struct {
    stations map[uint16]*stationMetrics
    totalFrames, totalDrops atomic.Int64
}

type StationMetrics struct {
    StationID, FramesTotal, FramesDropped int64
    LastSeen int64, Connected bool
    FPS float64
}

func NewMetrics() *Metrics
func (m *Metrics) RecordFrame(stationID uint16, name string)
func (m *Metrics) GetStations() []StationMetrics
func (m *Metrics) GetGlobal() GlobalMetrics
```

### internal/debug/http.go (NEW)
```go
Endpoints:
  GET /healthz          → 200 OK
  GET /debug/stations   → JSON station metrics
  GET /debug/metrics    → JSON global metrics
```

### internal/relay/relay.go (UPDATED)
- Replace log.Printf with structured logger
- Record frame counts, drops, reconnects

### internal/relay/dispatcher.go (UPDATED)
- Replace log.Printf with structured logger
- Record station events

## Test Results
| Test | Status |
|------|--------|
| Buffer tests (10) | PASS |
| All packages build | PASS |

## Dependencies
- github.com/google/gopacket v1.1.19
- github.com/go-gnss/rtcm v0.0.8
