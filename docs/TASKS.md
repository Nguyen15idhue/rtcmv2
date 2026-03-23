# TASKS

## Phase 1: Project Setup ✅
- [x] Init Go module
- [x] Define folder structure
- [x] Create main.go (entry point)

## Phase 2: RTCM Buffer ✅
- [x] Design buffer logic
- [x] Implement byte buffer
- [x] Detect RTCM frame (0xD3)
- [x] Extract full message frames
- [x] Handle edge cases (garbage, partial, multi-frame)
- [x] Unit tests (10/10 PASS)

## Phase 3: TCP Capture + Reassembly ✅
- [x] Config struct with pcap settings
- [x] StreamFactory struct
- [x] tcpStream struct with tcpreader.ReaderStream
- [x] Stream reading loop
- [x] Buffer integration (buffer.Write())
- [x] Capture loop (pcap reading)
- [x] Frame output via channel
- [x] Graceful shutdown

## Phase 4a: Basic NTRIP Relay ✅
- [x] JSON config loading (multi-caster support)
- [x] NTRIP client (TCP connect)
- [x] SOURCE authentication
- [x] Channel input + goroutine
- [x] Reconnect loop
- [x] Write timeout
- [x] Response parsing (ICY 200 / ERROR)
- [x] Multi-caster dispatcher

## Phase 4b: RTCM Parsing + Station Routing ✅
- [x] RTCM message parsing (go-gnss/rtcm)
- [x] Extract Station ID from MSM header
- [x] Map Station ID → caster config
- [x] Per-station relay channels

## Phase 5: Frame Dispatcher ✅
- [x] Group frames by StationID
- [x] Create per-station channel
- [x] Route frames to corresponding relay
- [x] Dynamic station discovery
  - [x] RWMutex for thread safety
  - [x] MaxDynamicStations limit
  - [x] IdleTimeout cleanup
  - [x] Dynamic mountpoint (/STATION_<id>)

## Phase 6: NTRIP Relay (DONE in Phase 4a) ✅
- [x] One connection per station
- [x] Mountpoint mapping
- [x] Send RTCM frames
- [x] Handle caster response

## Phase 7: Connection Management (DONE in Phase 4a) ✅
- [x] Auto-reconnect
- [x] Backpressure handling
- [x] Non-blocking writes

## Phase 8: Flow & Memory Control (DONE in Phase 5) ✅
- [x] Limit buffer size
- [x] Cleanup idle stations
- [x] Drop policy if overload

## Phase 9: Observability ✅

### 9.1 Logging (Structured) ✅
- [x] Replace fmt.Println with structured logger (slog)
- [x] Log relay lifecycle events
  - [x] relay_connected (stationID, caster)
  - [x] relay_write_error
  - [x] reconnect_attempt
- [x] Log dispatcher events
  - [x] new_station_detected
  - [x] station_removed (idle timeout)
- [x] Log drop events
- [x] Add log levels (INFO, WARN, ERROR)

### 9.2 Metrics (In-Memory) ✅
- [x] Create metrics struct (thread-safe)
- [x] Track per-station metrics
  - [x] frames_total
  - [x] frames_dropped
  - [x] last_seen_timestamp
  - [x] fps (frames/sec)
  - [x] connected status
- [x] Track global metrics
  - [x] active_stations
  - [x] total_frames
  - [x] total_drops
  - [x] uptime_seconds
- [x] Update metrics inside dispatcher & relay

### 9.3 Debug HTTP Server ✅
- [x] Create internal/debug/http.go
- [x] Start HTTP server (separate goroutine)

Endpoints:
- [x] GET /healthz
- [x] GET /debug/stations
- [x] GET /debug/metrics

### 9.4 FPS Calculation ✅
- [x] Implement sliding window (5s)
- [x] Compute frames/sec per station
- [x] Avoid heavy locking (per-station mutex)

---

## Phase 10: Production Hardening ✅

### 10.1 Better Reconnect ✅
- [x] Implement exponential backoff
  - [x] start: 1s
  - [x] max: 60s
- [x] Reset backoff on successful connect

### 10.2 Panic Safety ✅
- [x] Add recover() in all goroutines
- [x] Log panic errors
- [x] safeGo() helper function

### 10.3 Graceful Shutdown ✅
- [x] Handle SIGINT / SIGTERM
- [x] Stop capture loop
- [x] Close channels safely
- [x] Close all relay connections
- [x] WaitGroup for goroutines
- [x] Timeout (10s)

### 10.4 Basic Performance (Optional) ⏸
- [ ] Avoid unnecessary []byte copy
- [ ] Optional: use sync.Pool for buffers


## Phase 11: Mini Web Dashboard

### 11.1 HTTP Server (reuse debug server)
- [ ] Use net/http
- [ ] Run on :1507 (configurable)

---

### 11.2 API Endpoints

- [ ] GET /api/stations
  - [ ] stationID
  - [ ] fps
  - [ ] total_frames
  - [ ] dropped_frames
  - [ ] connected
  - [ ] last_seen

- [ ] GET /api/system
  - [ ] active_stations
  - [ ] total_frames
  - [ ] total_drops
  - [ ] uptime

- [ ] POST /api/station/:id/reconnect
- [ ] POST /api/station/:id/disable (optional)

---

### 11.3 Web UI (static HTML)



---

### 11.4 UI Features

- [ ] Table: list stations
- [ ] Show:
  - [ ] station ID
  - [ ] FPS
  - [ ] drops
  - [ ] status (🟢/🔴)
- [ ] Auto refresh (2s via JS)
- [ ] Button: reconnect

---

### 11.5 Basic Styling
- [ ] Simple CSS (no framework)
- [ ] Or use CDN (optional: Tailwind)

---

### 11.6 Optional Enhancements
- [ ] Chart (FPS over time)
- [ ] Filter/search station

## Dependencies
- github.com/google/gopacket v1.1.19
- github.com/go-gnss/rtcm v0.0.8
