package relay

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Dispatcher struct {
	mu             sync.RWMutex
	cfg            Config
	frameChan      <-chan []byte
	relays         []*Relay
	inputChans     map[uint16]chan []byte
	parser         *StationParser
	dynamicRelays  map[uint16]*time.Time
	relayByStation map[*Relay]uint16
	metrics        *Metrics
}

func NewDispatcher(cfg Config, frameChan <-chan []byte, metrics *Metrics) *Dispatcher {
	relays := make([]*Relay, 0, len(cfg.Casters))
	inputChans := make(map[uint16]chan []byte)

	for _, caster := range cfg.Casters {
		if caster.StationID == 0 {
			continue
		}
		ch := make(chan []byte, 100)
		relay := NewRelay(caster, cfg.Relay, ch, metrics)
		relays = append(relays, relay)
		inputChans[caster.StationID] = ch
	}

	return &Dispatcher{
		cfg:            cfg,
		frameChan:      frameChan,
		relays:         relays,
		inputChans:     inputChans,
		parser:         NewStationParser(),
		dynamicRelays:  make(map[uint16]*time.Time),
		relayByStation: make(map[*Relay]uint16),
		metrics:        metrics,
	}
}

func (d *Dispatcher) Start(ctx context.Context) {
	for _, relay := range d.relays {
		go func(r *Relay) {
			if err := r.Run(ctx); err != nil {
				logWarn("relay_stopped", LogFields{
					Name:  r.config.Name,
					Error: err.Error(),
				})
			}
		}(relay)
	}

	go d.route()
	go d.cleanupLoop()
}

func (d *Dispatcher) route() {
	for frame := range d.frameChan {
		stationID, err := d.parser.ExtractStationID(frame)
		if err != nil {
			logWarn("unknown_frame", LogFields{
				Event: "parse_error",
			})
			continue
		}
		if stationID == 0 {
			continue
		}

		d.mu.Lock()
		ch, ok := d.inputChans[stationID]
		if !ok {
			if d.cfg.DefaultCaster == nil {
				d.mu.Unlock()
				continue
			}
			if len(d.dynamicRelays) >= d.cfg.Relay.MaxDynamicStations {
				d.mu.Unlock()
				logWarn("station_rejected", LogFields{
					StationID: stationID,
					Event:     "max_dynamic_reached",
				})
				continue
			}
			ch = d.createRelay(stationID)
			now := time.Now()
			d.dynamicRelays[stationID] = &now
		} else {
			if ts := d.dynamicRelays[stationID]; ts != nil {
				*ts = time.Now()
			}
		}
		d.mu.Unlock()

		select {
		case ch <- frame:
		default:
			if d.metrics != nil {
				s := d.metrics.getOrCreateStation(stationID, "unknown")
				s.RecordDrop()
			}
			logWarn("frame_dropped", LogFields{
				StationID: stationID,
				Event:     "channel_full",
			})
		}
	}
}

func (d *Dispatcher) createRelay(stationID uint16) chan []byte {
	caster := *d.cfg.DefaultCaster
	caster.StationID = stationID
	caster.Mountpoint = fmt.Sprintf("/STATION_%d", stationID)
	caster.Name = fmt.Sprintf("dynamic-%d", stationID)

	ch := make(chan []byte, 100)
	relay := NewRelay(caster, d.cfg.Relay, ch, d.metrics)

	d.inputChans[stationID] = ch
	d.relays = append(d.relays, relay)
	d.relayByStation[relay] = stationID

	ctx := context.Background()
	go func() {
		if err := relay.Run(ctx); err != nil {
			logWarn("dynamic_relay_stopped", LogFields{
				StationID: stationID,
				Error:     err.Error(),
			})
		}
	}()

	logInfo("new_station_detected", LogFields{
		StationID:  stationID,
		Name:       caster.Name,
		Mountpoint: caster.Mountpoint,
		Caster:     caster.Host,
	})

	return ch
}

func (d *Dispatcher) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		d.cleanupIdle()
	}
}

func (d *Dispatcher) cleanupIdle() {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now()
	for stationID, lastSeen := range d.dynamicRelays {
		if now.Sub(*lastSeen) > d.cfg.Relay.IdleTimeout {
			d.removeRelayLocked(stationID)
		}
	}
}

func (d *Dispatcher) removeRelayLocked(stationID uint16) {
	ch, ok := d.inputChans[stationID]
	if !ok {
		return
	}

	delete(d.dynamicRelays, stationID)
	delete(d.inputChans, stationID)
	close(ch)

	for i, relay := range d.relays {
		if d.relayByStation[relay] == stationID {
			relay.Stop()
			d.relays = append(d.relays[:i], d.relays[i+1:]...)
			delete(d.relayByStation, relay)
			break
		}
	}

	if d.metrics != nil {
		d.metrics.RemoveStation(stationID)
	}

	logInfo("station_removed", LogFields{
		StationID: stationID,
		Event:     "idle_timeout",
	})
}

func (d *Dispatcher) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, relay := range d.relays {
		relay.Stop()
	}
	for _, ch := range d.inputChans {
		close(ch)
	}
}

func (d *Dispatcher) GetMetrics() *Metrics {
	return d.metrics
}
