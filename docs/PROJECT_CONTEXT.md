# Project: RTCMv2 Relay

## Goal
Capture TCP stream from existing port without affecting service,
parse RTCM data, split by station, relay to another NTRIP caster.

## Key Requirements
- Real-time processing
- No data loss
- No interference with original service
- Multi-station support
- Auto reconnect

## Tech Stack
- Go
- go-gnss/rtcm
- TCP + NTRIP

## Constraints
- TCP is continuous stream (no boundaries)
- RTCM is binary protocol
- Must buffer and reconstruct messages