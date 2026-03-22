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

## Phase 10: Prometheus Monitoring

### 10.1 Prometheus Exporter
- [ ] Add /metrics endpoint
- [ ] Export metrics:
  - [ ] rtcm_frames_total{station}
  - [ ] rtcm_dropped_total{station}
  - [ ] rtcm_active_stations
  - [ ] rtcm_reconnect_total

---

### 10.2 Integration
- [ ] Test with Prometheus scrape
- [ ] Validate metric labels

---

### 10.3 Grafana (Optional)
- [ ] Create dashboard JSON
  - [ ] FPS per station
  - [ ] active stations
  - [ ] drop rate
  - [ ] reconnect count

---

## Phase 11: Reliability & Hardening

### 11.1 Reconnect Strategy
- [ ] Implement exponential backoff
  - [ ] min delay (1s)
  - [ ] max delay (60s)
- [ ] Reset backoff on success

---

### 11.2 Circuit Breaker (Per Relay)
- [ ] Track consecutive failures
- [ ] Stop reconnect after threshold
- [ ] Cooldown period before retry

---

### 11.3 Panic Recovery
- [ ] Add recover() in all goroutines
- [ ] Log panic with stack trace

---

### 11.4 Graceful Shutdown
- [ ] Handle SIGINT / SIGTERM
- [ ] Stop capture loop
- [ ] Close frame channel
- [ ] Drain dispatcher
- [ ] Close all relays
- [ ] Wait for goroutines to exit

---

## Phase 12: Performance & Scaling

### 12.1 Memory Optimization
- [ ] Reuse buffers (sync.Pool)
- [ ] Avoid unnecessary []byte copy
- [ ] Profile memory usage

---

### 12.2 Lock Optimization
- [ ] Review RWMutex usage
- [ ] Reduce lock contention
- [ ] Consider sharded map if needed

---

### 12.3 Benchmark
- [ ] Write benchmark for:
  - [ ] frame parsing
  - [ ] dispatcher routing
- [ ] Run go test -bench=.

---

### 12.4 Load Testing
- [ ] Simulate multiple stations
- [ ] Measure:
  - [ ] CPU
  - [ ] memory
  - [ ] drop rate

---

## Phase 13: Advanced Features (Optional)

### 13.1 RTCM Replay
- [ ] Save frames to file
- [ ] Replay tool for testing

---

### 13.2 Message Filtering
- [ ] Filter by RTCM message type
- [ ] Configurable allow/deny list

---

### 13.3 Multi-Caster Output
- [ ] One station → multiple casters
- [ ] Config support

---

### 13.4 Config Hot Reload
- [ ] Watch config.json
- [ ] Reload without restart


## Dependencies
- github.com/google/gopacket v1.1.19
- github.com/go-gnss/rtcm v0.0.8
