package relay

import (
	"encoding/json"
	"os"
	"time"
)

type CaptureConfig struct {
	Interface string `json:"interface"`
	Filter    string `json:"filter"`
}

type CasterConfig struct {
	Name       string `json:"name"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Mountpoint string `json:"mountpoint"`
	Password   string `json:"password"`
	StationID  uint16 `json:"station_id"`
}

type RelayConfig struct {
	ReconnectInterval  time.Duration `json:"reconnect_interval"`
	WriteTimeout       time.Duration `json:"write_timeout"`
	Agent              string        `json:"agent"`
	MaxDynamicStations int           `json:"max_dynamic_stations"`
	IdleTimeout        time.Duration `json:"idle_timeout"`
}

type Config struct {
	Capture       CaptureConfig  `json:"capture"`
	Casters       []CasterConfig `json:"casters"`
	DefaultCaster *CasterConfig  `json:"default_caster"`
	Relay         RelayConfig    `json:"relay"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.Relay.ReconnectInterval == 0 {
		cfg.Relay.ReconnectInterval = 5 * time.Second
	}
	if cfg.Relay.WriteTimeout == 0 {
		cfg.Relay.WriteTimeout = 10 * time.Second
	}
	if cfg.Relay.Agent == "" {
		cfg.Relay.Agent = "rtcmv2-relay/1.0"
	}
	if cfg.Relay.MaxDynamicStations == 0 {
		cfg.Relay.MaxDynamicStations = 10
	}
	if cfg.Relay.IdleTimeout == 0 {
		cfg.Relay.IdleTimeout = 300 * time.Second
	}

	return &cfg, nil
}
