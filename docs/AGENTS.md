# Agent Instructions for RTCMv2 Relay (Sniffer Architecture)

Passive TCP relay system for RTCM data.
Captures network traffic WITHOUT connecting to source server,
reassembles TCP streams, extracts RTCM frames, and relays to NTRIP caster.

**Stack**: Go, gopacket, tcpassembly, NTRIP protocol

---

## 🎯 Project Goal

- Capture RTCM data from an existing TCP port WITHOUT affecting original service
- Do NOT create any connection to source server
- Reassemble TCP stream correctly
- Extract RTCM frames
- Relay frames to another NTRIP caster reliably

---

## ⚠️ Core Principles (MANDATORY)

- Correctness FIRST (especially TCP reassembly)
- Do NOT over-engineer
- Keep modules minimal and isolated
- Build step-by-step (phase-based)
- NEVER break existing working modules
- NEVER implement multiple phases at once

---

## 🧠 Development Workflow (CRITICAL)

1. Explain approach first (NO code)
2. Break into smallest possible step
3. Implement ONLY that step
4. WAIT for approval before continuing

NEVER:
- Write full system at once
- Jump across phases
- Modify unrelated code
- Add features not requested

---

## 🔒 Scope Control

- Do NOT redesign architecture
- Do NOT introduce new features
- Only implement EXACT requested functionality
- If unclear → ASK before coding

---

## 🏗️ Project Structure

rtcmv2/
├── cmd/relay/main.go
├── internal/
│ ├── capture/       # packet capture (gopacket)
│ ├── reassembly/    # tcpassembly logic
│ ├── buffer/        # RTCM frame extraction (DONE)
│ ├── relay/         # NTRIP client
│ └── config/

---

## 🔄 System Architecture (CRITICAL)

Modules:

- capture: sniff packets (pcap)
- reassembly: rebuild TCP stream (tcpassembly)
- buffer: extract RTCM frames (already implemented)
- relay: send frames to caster

---

## 🔁 Data Flow (CRITICAL)

Network Packets
    ↓
gopacket capture
    ↓
tcpassembly (reorder TCP)
    ↓
byte stream (io.Reader)
    ↓
buffer.Write()
    ↓
RTCM frames
    ↓
relay to caster

---

## 🚫 IMPORTANT DIFFERENCE

This is NOT a TCP client system.

- Do NOT use net.Dial to source server
- Do NOT create connections to source
- ONLY sniff traffic via network interface

---

## ⚡ TCP Reassembly Rules (CRITICAL)

- Use tcpassembly (DO NOT implement manually)
- Each TCP flow must have its own stream
- Use goroutine per stream
- Read stream via tcpreader.ReaderStream

NEVER:
- Append raw packets directly
- Assume packet order is correct
- Ignore retransmissions

---

## ⚡ RTCM Buffer Rules (CRITICAL)

(ALREADY IMPLEMENTED — DO NOT MODIFY)

- Input: []byte stream
- Output: [][]byte (complete RTCM frames)

Rules:
- Sync byte = 0xD3
- Length = 10-bit field
- Frame size = 6 + length

NEVER:
- Modify buffer logic unless explicitly asked
- Parse RTCM here (buffer ONLY extracts frames)

---

## 🛰️ RTCM Handling (DEFERRED)

- DO NOT implement full RTCM parsing yet
- DO NOT integrate go-gnss/rtcm in early phases

Only:
- treat RTCM as binary frames
- forward as-is

---

## 🌐 Relay (NTRIP)

- Connect ONLY to destination caster
- Send:
  SOURCE <password> <mountpoint>

- Send frames immediately (low latency)

- Must handle:
  - reconnect
  - timeout
  - broken pipe

---

## 🔌 Network Rules

- Capture via interface (pcap)
- Apply BPF filter:
  tcp and port XXXX

- Do NOT capture unnecessary traffic

---

## 🔁 Failure Handling (CRITICAL)

- System MUST NOT crash
- Auto-reconnect relay
- Drop broken flows safely
- Continue processing other streams

---

## ⚙️ Concurrency Rules

- Goroutine per TCP stream
- Do NOT block packet capture loop
- Avoid memory leaks
- Clean up inactive streams

---

## 🧪 Testing Strategy

- Test buffer separately ✅
- Test reassembly with sample packets
- Test relay with mock caster

---

## 🧹 Code Style

- Minimal code only
- No unnecessary abstraction
- Small functions
- Clear responsibilities

---

## 🧠 Implementation Rules (VERY IMPORTANT)

- ONE function at a time
- ONE phase at a time
- NEVER jump ahead
- NEVER combine capture + relay in one step

---

## 🐛 Debugging Rules

- Identify ROOT CAUSE
- Fix minimal code only
- Do NOT rewrite modules

---

## 🚀 Development Priority (UPDATED)

1. Packet capture (gopacket)
2. TCP reassembly (tcpassembly) 🔥
3. Integrate buffer (stream → frames)
4. Relay (NTRIP)
5. Reconnect logic
6. Stability & cleanup
7. (Optional) RTCM parsing later

---

## 🚨 Anti-Overengineering Rule

If solution includes:
- unnecessary abstraction
- complex architecture
- unused components

→ STOP and simplify

---

## 🧪 Running

```bash
go run ./cmd/relay/main.go
🧪 Testing
go test -v ./...
go test -cover ./...