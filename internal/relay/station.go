package relay

import (
	"github.com/go-gnss/rtcm/rtcm3"
)

type StationParser struct {
}

func NewStationParser() *StationParser {
	return &StationParser{}
}

func (p *StationParser) ExtractStationID(frame []byte) (uint16, error) {
	if len(frame) < 6 {
		return 0, nil
	}

	length := int(frame[1]&0x03)<<8 | int(frame[2])
	if len(frame) < 3+length {
		return 0, nil
	}

	payload := frame[3 : 3+length]
	if len(payload) < 4 {
		return 0, nil
	}

	msgType, err := rtcm3.MessageNumber(payload)
	if err != nil {
		return 0, nil
	}

	switch {
	case msgType >= 1074 && msgType <= 1077:
		return p.extractFromMSM(payload)
	case msgType >= 1084 && msgType <= 1087:
		return p.extractFromMSM(payload)
	case msgType >= 1094 && msgType <= 1097:
		return p.extractFromMSM(payload)
	case msgType >= 1124 && msgType <= 1127:
		return p.extractFromMSM(payload)
	default:
		return 0, nil
	}
}

func (p *StationParser) extractFromMSM(payload []byte) (uint16, error) {
	if len(payload) < 4 {
		return 0, nil
	}

	stationID := uint16(payload[2]) | (uint16(payload[3]&0x0F) << 8)
	return stationID, nil
}
