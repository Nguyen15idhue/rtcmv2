package relay

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

type Relay struct {
	config   CasterConfig
	relayCfg RelayConfig
	input    <-chan []byte
	metrics  *Metrics
	conn     net.Conn
	stopped  bool
}

func NewRelay(caster CasterConfig, relayCfg RelayConfig, input <-chan []byte, metrics *Metrics) *Relay {
	return &Relay{
		config:   caster,
		relayCfg: relayCfg,
		input:    input,
		metrics:  metrics,
	}
}

func (r *Relay) Run(ctx context.Context) error {
	for !r.stopped {
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := r.connect(ctx); err != nil {
				logWarn("reconnect_attempt", LogFields{
					Name:  r.config.Name,
					Error: err.Error(),
				})
				time.Sleep(r.relayCfg.ReconnectInterval)
				continue
			}
			r.sendLoop(ctx)
		}
	}
	return nil
}

func (r *Relay) Stop() {
	r.stopped = true
	if r.conn != nil {
		r.conn.Close()
	}
}

func (r *Relay) connect(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", r.config.Host, r.config.Port)

	dialer := net.Dialer{Timeout: 10 * time.Second}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return err
	}

	req := fmt.Sprintf("SOURCE %s %s\r\nSource-Agent: %s\r\n\r\n",
		r.config.Password,
		r.config.Mountpoint,
		r.relayCfg.Agent,
	)

	if _, err := conn.Write([]byte(req)); err != nil {
		conn.Close()
		return fmt.Errorf("send request: %w", err)
	}

	resp, err := r.readResponse(conn)
	if err != nil {
		conn.Close()
		return fmt.Errorf("read response: %w", err)
	}

	if !strings.HasPrefix(resp, "ICY 200") {
		conn.Close()
		return fmt.Errorf("caster rejected: %s", resp)
	}

	r.conn = conn

	if r.metrics != nil {
		s := r.metrics.getOrCreateStation(r.config.StationID, r.config.Name)
		s.SetConnected(true)
	}

	logInfo("relay_connected", LogFields{
		StationID:  r.config.StationID,
		Name:       r.config.Name,
		Caster:     r.config.Host,
		Mountpoint: r.config.Mountpoint,
	})

	return nil
}

func (r *Relay) readResponse(conn net.Conn) (string, error) {
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	reader := bufio.NewReader(conn)
	resp, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	return strings.TrimSpace(resp), nil
}

func (r *Relay) sendLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case frame, ok := <-r.input:
			if !ok {
				return
			}
			if r.conn == nil {
				continue
			}
			r.conn.SetWriteDeadline(time.Now().Add(r.relayCfg.WriteTimeout))
			if _, err := r.conn.Write(frame); err != nil {
				logError("relay_write_error", LogFields{
					Name:      r.config.Name,
					StationID: r.config.StationID,
					Error:     err.Error(),
				})
				r.conn.Close()
				r.conn = nil
				if r.metrics != nil {
					r.metrics.stations[r.config.StationID].SetConnected(false)
				}
				return
			}
			if r.metrics != nil {
				r.metrics.RecordFrame(r.config.StationID, r.config.Name)
			}
		case <-time.After(1 * time.Second):
		}
	}
}
