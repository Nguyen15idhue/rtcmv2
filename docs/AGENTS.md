# Agent Instructions for RTCMv2 Relay

TCP stream relay for RTCM data. Captures TCP from a port, parses RTCM stream, splits by station, relays to NTRIP caster. Built with Go.

**Stack**: Go, go-gnss/rtcm, NTRIP protocol

**Project Goal**:
Capture data from an existing TCP port WITHOUT affecting the original service,
parse RTCM stream into individual stations, and relay to another caster reliably.

---

## ⚠️ Core Principles (MANDATORY)

- Correctness and stability FIRST
- Do NOT over-engineer
- Keep modules simple and isolated
- Build skeleton first, then improve
- Never break existing structure
- Never implement full system at once

---

## 🧠 Development Workflow (CRITICAL)

1. Explain approach first (NO code)
2. Break into small steps
3. Implement ONE step only
4. Wait for approval before continuing

NEVER:
- Write full system at once
- Modify unrelated code
- Add extra features

---

## 🔒 Scope Control

- Do NOT redesign architecture
- Do NOT introduce new features unless asked
- Only implement requested module/function
- Warn if request may break existing structure

---

## 🏗️ Project Structure


rtcmv2/
├── cmd/relay/main.go
├── internal/
│ ├── capture/ # TCP capture
│ ├── buffer/ # byte stream reassembly
│ ├── parser/ # RTCM decoding
│ ├── router/ # station split
│ ├── relay/ # NTRIP client
│ └── config/


---

## 🔄 System Architecture

Modules:

- capture: reads raw TCP stream
- buffer: reassembles byte stream
- parser: decodes RTCM messages
- router: splits by station
- relay: sends to NTRIP caster

---

## 🔁 Data Flow (CRITICAL)

TCP Stream → Buffer → RTCM Decode → Station Router → Relay

IMPORTANT:

- TCP is a continuous byte stream
- No packet/message boundaries guaranteed
- RTCM messages must be reconstructed from bytes
- Multiple stations may exist in one stream

---

## ⚡ RTCM Buffer Strategy (CRITICAL)

- Maintain continuous byte buffer per connection
- Search for sync byte: 0xD3
- Next 2 bytes define message length
- Wait until FULL message is available
- Then parse message

After parsing:
- Remove processed bytes from buffer

NEVER assume:
- One TCP read = one RTCM message
- Packets are aligned

---

## 🛰️ RTCM Handling

- Use go-gnss/rtcm
- Extract station ID from messages (MSM)
- Support multiple stations in same stream
- Handle partial messages safely

---

## 🔀 Station Routing Strategy

- Extract stationID from RTCM message
- Maintain map: stationID → relay client
- Each station has its own mountpoint
- Route messages by stationID

NEVER:
- Mix data between stations

---

## 🌐 NTRIP Protocol

- Use NTRIP v1/v2
- Authenticate with:
  SOURCE <password> <mountpoint>
- Send:
  Source-Agent: rtcmv2-relay

- Handle:
  - connection drops
  - reconnect
  - basic responses

---

## 🔌 Network I/O Rules

- Use timeouts for all connections
- Handle partial reads
- Use buffer size: 8KB–64KB
- Implement reconnect with backoff (1–5s)

---

## 🔁 Failure Handling (CRITICAL)

- Do NOT crash on connection errors
- Reconnect automatically
- Use backoff retry
- Log all errors but continue running

---

## ⚙️ Concurrency

- Use goroutines for:
  - capture
  - relay per station

- Use:
  - channels for communication
  - sync.Mutex / RWMutex for shared state

- Avoid goroutine leaks

---

## 🧪 Testing

- Use table-driven tests
- Mock:
  - TCP input
  - NTRIP caster

- Test:
  - buffer correctness
  - RTCM parsing
  - routing logic

---

## 🧹 Code Style

- Use gofmt / goimports
- Wrap errors with context
- Keep functions small
- Avoid unnecessary abstractions

---

## 🧠 Implementation Rules (VERY IMPORTANT)

- Implement ONE function at a time
- Do NOT write entire module in one step
- Do NOT refactor unrelated code
- Always explain before coding

---

## 🐛 Debugging Rules

- Find ROOT CAUSE first
- Do NOT rewrite large code blocks
- Suggest minimal fix

---

## 🚀 Development Priority

1. Capture TCP stream
2. Implement buffer correctly
3. Parse RTCM messages
4. Split by station
5. Relay to caster
6. Add reconnect logic
7. Improve performance/logging

---

## 🧪 Running

```bash
go run ./cmd/relay/main.go --config=config.yaml
📊 Debugging
go test -v ./...
go test -cover ./...