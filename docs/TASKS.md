# TASKS

## Phase 1: Project Setup (Skeleton)
- [ ] Init Go module
- [ ] Create project structure
- [ ] Create empty entry point (main.go)

## Phase 2: Buffer (CRITICAL)
- [ ] Design buffer logic (NO CODE)
- [ ] Implement byte buffer
- [ ] Detect RTCM frame (0xD3)
- [ ] Extract full message

## Phase 3: TCP Input
- [ ] Connect to TCP source
- [ ] Read stream into buffer

## Phase 4: Parser
- [ ] Integrate rtcm library
- [ ] Parse messages

## Phase 5: Routing
- [ ] Extract station ID
- [ ] Map station → output

## Phase 6: Relay (NTRIP)
- [ ] Implement NTRIP client
- [ ] Send data

## Phase 7: Stability
- [ ] Reconnect logic
- [ ] Error handling