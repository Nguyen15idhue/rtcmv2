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
- [x] Task 17: SSE broadcasts metrics every 1s
- [x] Task 18: Notification events support

### Phase 11.3: Vue Frontend Setup (Tasks 19-25) ✅
- [x] Task 19: Install dependencies (vue-router, pinia, axios, apexcharts)
- [x] Task 20: Create router (`frontend/src/router/index.js`)
- [x] Task 21: Update main.js with router
- [x] Task 22: Stream store (`frontend/src/stores/stream.js`)
- [x] Task 23: Stations store (`frontend/src/stores/stations.js`)
- [x] Task 24: Casters store (`frontend/src/stores/casters.js`)
- [x] Task 25: API client (`frontend/src/api/index.js`)

---

## Phase 11.4: Dashboard Components (Tasks 26-33)

### 11.4.1 Layout

#### Task 26: Header component
- [x] Navigation bar với router-links
- [ ] Sidebar (optional)

#### Task 27: App layout
- [x] Navbar với brand, links, status
- [x] Router-view

### 11.4.2 Components

#### Task 28: StatsCard component
- [x] Stats grid với 4 cards

#### Task 29: StationTable component
- [x] Table với stations data từ SSE

#### Task 30: Search/Filter
- [ ] Thêm search input

### 11.4.3 Station Actions

#### Task 31: StationModal component
- [ ] Form thêm/sửa station

#### Task 32: ConfirmDialog component
- [ ] Dialog xác nhận xóa

#### Task 33: Dashboard view
- [x] Stats + Station table
- [ ] Thêm search/filter

---

## Phase 11.5: Station Detail & Outputs (Tasks 34-38)

#### Task 34: StationDetailModal
- [ ] Chi tiết station

#### Task 35: OutputTable
- [ ] Bảng outputs

#### Task 36: OutputModal
- [ ] Form thêm output

#### Task 37: LineChart
- [ ] FPS chart

#### Task 38: Integrate chart
- [ ] Thêm chart vào Dashboard

---

## Phase 11.6: Casters & Settings (Tasks 39-41)

#### Task 39: Casters view
- [ ] Table casters

#### Task 40: CasterModal
- [ ] Form thêm caster

#### Task 41: Settings view
- [ ] API Key, theme toggle

---

## Phase 11.7: Notifications & Polish (Tasks 42-45)

#### Task 42: Toast components
- [ ] Toast notifications

#### Task 43: ToastContainer
- [ ] Container cho toasts

#### Task 44: Notification store
- [ ] Store cho notifications

#### Task 45: UI Polish
- [ ] Dark mode, responsive

---

## Phase 11.8: Docker (Tasks 46-48)

#### Task 46: Dockerfile
- [ ] Go alpine image

#### Task 47: docker-compose.yml
- [ ] Service + volumes

#### Task 48: .dockerignore
- [ ] Ignore files

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
| DELETE | `/api/station/:id/output/:cid` | No | Remove output |
| PUT | `/api/station/:id/output/:cid/toggle` | No | Toggle output |
| GET | `/api/casters` | No | List casters |
| POST | `/api/caster` | No | Create caster |
| DELETE | `/api/caster/:id` | No | Delete caster |
| GET | `/api/stream` | No | SSE stream |
| GET | `/api/config` | Yes | Get full config |
| POST | `/api/reload` | Yes | Reload config |

## Frontend Files đã tạo

```
frontend/src/
├── main.js              # [UPDATE] Thêm router
├── App.vue              # [UPDATE] Navbar + router-view
├── router/
│   └── index.js         # [NEW] Vue Router
├── api/
│   └── index.js         # [NEW] Axios client
├── stores/
│   ├── stream.js        # [NEW] SSE connection
│   ├── stations.js       # [NEW] Stations CRUD
│   └── casters.js       # [NEW] Casters CRUD
└── views/
    ├── Dashboard.vue    # [NEW] Dashboard page
    ├── Stations.vue      # [NEW] Placeholder
    ├── Casters.vue       # [NEW] Placeholder
    └── Settings.vue      # [NEW] Placeholder
```

## Dependencies
- github.com/google/gopacket v1.1.19
- github.com/go-gnss/rtcm v0.0.8

## Frontend Dependencies
- vue ^3.4.0
- vue-router ^4.2.0
- pinia ^2.1.0
- axios ^1.6.0
- apexcharts ^3.45.0
