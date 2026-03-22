# RTCMv2 Relay - Tài liệu

Hệ thống relay luồng TCP cho dữ liệu RTCM. Bắt lưu lượng mạng một cách thụ động (không ảnh hưởng đến dịch vụ gốc), phân tích luồng RTCM, tách theo trạm, và chuyển tiếp đến nhiều NTRIP caster.

---

## Tính năng

### ✅ Đã hoàn thành

| Tính năng | Trạng thái |
|-----------|------------|
| Bắt TCP thụ động (gopacket) | ✅ |
| Tái tạo luồng TCP (tcpassembly) | ✅ |
| Trích xuất khung RTCM (0xD3 sync) | ✅ |
| Phân tích station ID RTCM (MSM messages) | ✅ |
| Relay đa caster | ✅ |
| Phát hiện trạm động | ✅ |
| Tự động kết nối lại | ✅ |
| Write timeout | ✅ |
| Kênh relay theo trạm | ✅ |
| Structured logging (slog JSON) | ✅ |
| Metrics in-memory | ✅ |
| Tính FPS (sliding window) | ✅ |
| Debug HTTP server | ✅ |

---

## Kiến trúc

```
┌─────────────────────────────────────────────────────────────────┐
│                        Giao diện mạng                           │
└─────────────────────────────────┬───────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                     gopacket (pcap)                             │
│                   Bắt gói tin với BPF filter                   │
└─────────────────────────────────┬───────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                    tcpassembly                                   │
│              Tái tạo luồng TCP                                  │
└─────────────────────────────────┬───────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                     tcpStream                                    │
│           Xử lý luồng cho mỗi kết nối                        │
└─────────────────────────────────┬───────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                        Buffer                                    │
│              Trích xuất khung RTCM hoàn chỉnh (0xD3)          │
└─────────────────────────────────┬───────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                      FrameChan                                    │
│                    ([]chan []byte)                              │
└─────────────────────────────────┬───────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Dispatcher                                   │
│           Định tuyến khung theo station ID                     │
│         Tạo/dọn dẹp relay động                                │
└─────────────────────────────────┬───────────────────────────────┘
                                  │
                    ┌─────────────┴─────────────┐
                    ▼                           ▼
          ┌─────────────────┐         ┌─────────────────┐
          │   Relay (1234)   │         │   Relay (5678)   │
          │   /STATION_1234 │         │   /STATION_5678 │
          └────────┬────────┘         └────────┬────────┘
                   │                           │
                   ▼                           ▼
          ┌─────────────────┐         ┌─────────────────┐
          │   Caster A       │         │   Caster B       │
          └─────────────────┘         └─────────────────┘
```

---

## Cấu trúc dự án

```
rtcmv2/
├── cmd/relay/
│   └── main.go              # Điểm bắt đầu
├── internal/
│   ├── buffer/
│   │   ├── buffer.go        # Trích xuất khung RTCM
│   │   └── buffer_test.go   # 10 unit tests
│   ├── capture/
│   │   ├── capture.go       # gopacket + tcpassembly
│   │   └── stream.go        # StreamFactory + tcpStream
│   ├── relay/
│   │   ├── config.go        # Load config JSON
│   │   ├── relay.go         # NTRIP client
│   │   ├── dispatcher.go    # Định tuyến khung
│   │   ├── station.go       # Phân tích station ID RTCM
│   │   ├── metrics.go       # Metrics in-memory
│   │   └── logger.go        # Structured logging
│   └── debug/
│       └── http.go          # Debug HTTP server
├── docs/
│   ├── AGENTS.md            # Hướng dẫn cho agent
│   ├── SESSION.md           # Phiên hiện tại
│   ├── TASKS.md            # Theo dõi task
│   └── README.md           # Tài liệu này
├── config.json              # Cấu hình mẫu
├── go.mod
└── go.sum
```

---

## Cấu hình

### config.json

```json
{
  "capture": {
    "interface": "eth0",
    "filter": "tcp port 12101"
  },
  "casters": [
    {
      "name": "caster-a",
      "host": "caster1.example.com",
      "port": 2101,
      "mountpoint": "/STATION1",
      "password": "secret1",
      "station_id": 1234
    }
  ],
  "default_caster": {
    "name": "default",
    "host": "caster.example.com",
    "port": 2101,
    "mountpoint": "/STATION",
    "password": "default_pass"
  },
  "relay": {
    "reconnect_interval": 5,
    "write_timeout": 10,
    "agent": "rtcmv2-relay/1.0",
    "max_dynamic_stations": 10,
    "idle_timeout": 300
  }
}
```

### Các trường cấu hình

| Trường | Kiểu | Mô tả | Mặc định |
|--------|------|--------|-----------|
| `capture.interface` | string | Giao diện mạng để bắt | - |
| `capture.filter` | string | BPF filter (vd: `tcp port 12101`) | - |
| `casters[]` | array | Cấu hình caster định sẵn | - |
| `casters[].station_id` | uint16 | Station ID RTCM để match | - |
| `default_caster` | object | Caster dự phòng cho trạm động | - |
| `relay.max_dynamic_stations` | int | Số relay động tối đa | 10 |
| `relay.idle_timeout` | int | Giây trước khi xóa relay nhàn rỗi | 300 |
| `relay.write_timeout` | int | Giây trước khi write timeout | 10 |

---

## Sử dụng

### Build
```bash
go build -o rtcmv2 ./cmd/relay
```

### Chạy
```bash
# Với config.json mặc định
./rtcmv2

# Với config tùy chỉnh
./rtcmv2 /path/to/config.json
```

### Debug Server
```bash
# Đặt địa chỉ debug
DEBUG_ADDR=:9090 ./rtcmv2
```

---

## Debug Endpoints

| Endpoint | Method | Mô tả |
|----------|--------|--------|
| `/healthz` | GET | Health check (trả về 200 OK) |
| `/debug/stations` | GET | Danh sách các trạm đang hoạt động với metrics |
| `/debug/metrics` | GET | Snapshot metrics toàn cục |

### Ví dụ Output

**GET /debug/stations**
```json
[
  {
    "station_id": 1234,
    "name": "dynamic-1234",
    "frames_total": 15420,
    "frames_dropped": 0,
    "last_seen": 1711000000000000000,
    "connected": true,
    "reconnects": 2,
    "fps": 15.4
  }
]
```

**GET /debug/metrics**
```json
{
  "active_stations": 3,
  "total_frames": 45230,
  "total_drops": 0,
  "reconnect_total": 0,
  "uptime_seconds": 3600
}
```

---

## Logging

Structured JSON logs với các fields:

| Event | Fields |
|-------|--------|
| `relay_connected` | station_id, name, caster, mountpoint |
| `relay_write_error` | station_id, name, error |
| `relay_stopped` | name, error |
| `reconnect_attempt` | name, error |
| `new_station_detected` | station_id, name, mountpoint, caster |
| `station_removed` | station_id, event |
| `station_rejected` | station_id, event |
| `frame_dropped` | station_id, event |
| `unknown_frame` | event |

---

## Dependencies

| Package | Version | Mục đích |
|---------|---------|-----------|
| github.com/google/gopacket | v1.1.19 | Bắt gói tin |
| github.com/go-gnss/rtcm | v0.0.8 | Phân tích message RTCM |

---

## Luồng dữ liệu

1. **Capture**: gopacket bắt TCP packets từ giao diện mạng
2. **Reassemble**: tcpassembly tái tạo luồng TCP theo thứ tự đúng
3. **Extract**: Buffer trích xuất các khung RTCM hoàn chỉnh (sync: 0xD3)
4. **Route**: Dispatcher trích xuất station ID và định tuyến đến relay phù hợp
5. **Relay**: NTRIP client gửi khung đến caster đích
6. **Monitor**: Metrics và logs theo dõi mọi hoạt động

---

## NTRIP Protocol

Relay implements NTRIP v1 protocol:

1. TCP connect đến caster
2. Gửi authentication:
   ```
   SOURCE <password> <mountpoint>\r\n
   Source-Agent: rtcmv2-relay/1.0\r\n\r\n
   ```
3. Chờ response `ICY 200`
4. Chuyển tiếp RTCM frames
5. Xử lý reconnect khi thất bại

---

## An toàn bộ nhớ

- RWMutex bảo vệ shared maps
- Channel buffers (100 frames)
- Idle timeout cleanup cho các relay động
- Giới hạn số dynamic stations tối đa
- Non-blocking sends với drop khi đầy

---

## Testing

```bash
# Chạy tất cả tests
go test -v ./...

# Chạy với coverage
go test -cover ./...

# Chạy test cụ thể
go test -v -run TestWrite_SingleCompleteFrame ./internal/buffer
```

---

## Các cải tiến trong tương lai

- Prometheus metrics export
- Exponential backoff
- Circuit breaker
- Config hot reload
- Message filtering
