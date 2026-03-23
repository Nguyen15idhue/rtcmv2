# TASKS - RTCMv2 Relay Production Dashboard

## Hoàn thành ✅

### Phase 1-10: Core System
- [x] RTCM Buffer
- [x] TCP Capture + Reassembly
- [x] NTRIP Relay
- [x] Frame Dispatcher
- [x] Logging & Metrics
- [x] Graceful Shutdown

### Phase 11.1: Backend Foundation (Tasks 1-14) ✅
- [x] Task 1: API Key middleware (`api/middleware/auth.go`)
- [x] Task 2: Integrate middleware into server
- [x] Task 3: Station structs (`internal/relay/stations.go`)
- [x] Task 4: Load/Save stations
- [x] Task 5: CRUD stations functions
- [x] Task 6: Unassigned functions
- [x] Task 7: Caster structs (`internal/relay/casters.go`)
- [x] Task 8: Load/Save casters
- [x] Task 9: CRUD casters functions
- [x] Task 10: Stations CRUD handlers (`api/handlers/stations.go`)
- [x] Task 11: Outputs handlers (AddOutput, RemoveOutput, ToggleOutput)
- [x] Task 12: Casters handlers (`api/handlers/casters.go`)
- [x] Task 13: Config handlers (`api/handlers/config.go`)
- [x] Task 14: Update server routes

### Phase 11.2: Real-time Updates (SSE) (Tasks 15-18) ✅
- [x] Task 15: Broadcaster struct (`api/handlers/sse.go`)
- [x] Task 16: SSE Handler (GET /api/stream)
- [x] Task 17: Update dispatcher (not fully integrated)
- [x] Task 18: Notification events

---

## Phase 11.3: Vue Frontend Setup (Tasks 19-25)

### 11.3.1 Vue Router

#### Task 19: Install dependencies
- [ ] Thêm vào frontend/package.json:
  ```json
  "vue-router": "^4.2.0",
  "pinia": "^2.1.0",
  "axios": "^1.6.0",
  "apexcharts": "^3.45.0",
  "vue3-apexcharts": "^1.5.0",
  "lucide-vue-next": "^0.300.0"
  ```
- [ ] Chạy `npm install`

#### Task 20: Create router
- [ ] Tạo file `frontend/src/router/index.js`
- [ ] Định nghĩa routes:
  ```js
  const routes = [
    { path: '/', component: () => import('../views/Dashboard.vue') },
    { path: '/stations', component: () => import('../views/Stations.vue') },
    { path: '/casters', component: () => import('../views/Casters.vue') },
    { path: '/settings', component: () => import('../views/Settings.vue') },
  ]
  ```

#### Task 21: Update main.js
- [ ] Import router và store
- [ ] Register plugin

### 11.3.2 Pinia Stores

#### Task 22: Stream store
- [ ] Tạo file `frontend/src/stores/stream.js`
- [ ] State: `connected`, `lastData`, `lastUpdate`
- [ ] Actions: `connect()`, `disconnect()`
- [ ] Connect SSE và update state liên tục

#### Task 23: Stations store
- [ ] Tạo file `frontend/src/stores/stations.js`
- [ ] State: `stations`, `loading`, `error`
- [ ] Actions: `fetchStations()`, `createStation()`, `updateStation()`, `deleteStation()`
- [ ] Actions: `addOutput()`, `removeOutput()`, `toggleOutput()`

#### Task 24: Casters store
- [ ] Tạo file `frontend/src/stores/casters.js`
- [ ] State: `casters`, `loading`, `error`
- [ ] Actions: `fetchCasters()`, `createCaster()`, `deleteCaster()`

### 11.3.3 API Client

#### Task 25: API client
- [ ] Tạo file `frontend/src/api/index.js`
- [ ] Tạo axios instance
- [ ] Interceptor thêm `X-API-Key` header

---

## Phase 11.4: Dashboard Page (Tasks 26-33)

### 11.4.1 Layout

#### Task 26: Header component
- [ ] Tạo file `frontend/src/components/Header.vue`
- [ ] Logo "RTCMv2 Relay"
- [ ] Navigation links

#### Task 27: App layout
- [ ] Update `frontend/src/App.vue`
- [ ] Thêm Header
- [ ] Basic CSS grid layout

### 11.4.2 Components

#### Task 28: StatsCard component
- [ ] Tạo file `frontend/src/components/StatsCard.vue`
- [ ] Props: `label`, `value`, `color`

#### Task 29: StationTable component
- [ ] Tạo file `frontend/src/components/StationTable.vue`
- [ ] Columns: ID, Name, Status, FPS, Actions

#### Task 30: Search/Filter
- [ ] Thêm search input vào StationTable

### 11.4.3 Station Actions

#### Task 31: StationModal component
- [ ] Tạo file `frontend/src/components/StationModal.vue`
- [ ] Form: Name input

#### Task 32: ConfirmDialog component
- [ ] Tạo file `frontend/src/components/ConfirmDialog.vue`

#### Task 33: Dashboard view
- [ ] Tạo file `frontend/src/views/Dashboard.vue`
- [ ] Connect SSE để update real-time

---

## Phase 11.5: Station Detail & Outputs (Tasks 34-38)

#### Task 34: StationDetailModal component
- [ ] Props: `station`
- [ ] Display: ID, Name, Status, FPS, Frames, Drops

#### Task 35: OutputTable component
- [ ] Columns: Caster, Mountpoint, Enabled, Actions

#### Task 36: OutputModal component
- [ ] Select caster dropdown
- [ ] Mountpoint input

#### Task 37: LineChart component
- [ ] Use ApexCharts
- [ ] FPS over time

#### Task 38: Integrate chart vào Dashboard
- [ ] Update data từ SSE stream

---

## Phase 11.6: Casters & Settings (Tasks 39-41)

#### Task 39: Casters view
- [ ] Table: Name, Host, Port

#### Task 40: CasterModal component
- [ ] Form: Name, Host, Port, Password

#### Task 41: Settings view
- [ ] API Key section
- [ ] Theme toggle

---

## Phase 11.7: Notifications & Polish (Tasks 42-45)

#### Task 42: Toast components
- [ ] Auto-dismiss sau 5s

#### Task 43: ToastContainer component
- [ ] Position: top-right fixed

#### Task 44: Notification store
- [ ] Actions: `success()`, `error()`, `warning()`, `info()`

#### Task 45: UI Polish
- [ ] CSS variables for colors
- [ ] Dark mode support

---

## Phase 11.8: Docker Deployment (Tasks 46-48)

#### Task 46: Dockerfile
```dockerfile
FROM golang:1.24-alpine
WORKDIR /app
COPY . .
RUN go build -o relay ./cmd/server
EXPOSE 1507
CMD ["./relay"]
```

#### Task 47: docker-compose.yml
- [ ] Service: app
- [ ] Volume cho config

#### Task 48: .dockerignore
- [ ] Ignore: node_modules, .git, frontend/

---

## API Endpoints đã implement

| Method | Endpoint | Auth | Mô tả |
|--------|----------|------|--------|
| GET | `/api/health` | No | Health check |
| GET | `/api/system` | No | System stats |
| GET | `/api/stations` | No | List stations |
| POST | `/api/station` | No | Create station |
| PUT | `/api/station/:id` | No | Update station |
| DELETE | `/api/station/:id` | No | Delete station |
| POST | `/api/station/:id/output` | No | Add output |
| DELETE | `/api/station/:id/output/:caster_id` | No | Remove output |
| PUT | `/api/station/:id/output/:caster_id/toggle` | No | Toggle output |
| GET | `/api/casters` | No | List casters |
| POST | `/api/caster` | No | Create caster |
| DELETE | `/api/caster/:id` | No | Delete caster |
| GET | `/api/stream` | No | SSE stream |
| GET | `/api/config` | Yes | Get full config |
| POST | `/api/reload` | Yes | Reload config |

## Files đã tạo

```
api/
├── middleware/
│   └── auth.go              # [NEW]
├── handlers/
│   ├── stations.go          # [NEW]
│   ├── casters.go          # [NEW]
│   ├── config.go           # [NEW]
│   └── sse.go              # [NEW]
└── server.go               # [UPDATE]

internal/relay/
├── stations.go              # [NEW]
├── casters.go              # [NEW]
└── dispatcher.go           # [UPDATE]
```

## Dependencies
- github.com/google/gopacket v1.1.19
- github.com/go-gnss/rtcm v0.0.8

## Frontend Dependencies (cần cài)
- vue ^3.4.0
- vue-router ^4.2.0
- pinia ^2.1.0
- axios ^1.6.0
- apexcharts ^3.45.0
- vue3-apexcharts ^1.5.0
- lucide-vue-next ^0.300.0
