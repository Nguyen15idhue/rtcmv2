package relay

import (
	"encoding/json"
	"os"
	"time"
)

type Output struct {
	CasterID   uint16 `json:"caster_id"`
	Mountpoint string `json:"mountpoint"`
	Enabled    bool   `json:"enabled"`
}

type Station struct {
	ID        uint16    `json:"id"`
	Name      string    `json:"name"`
	Outputs   []Output  `json:"outputs"`
	CreatedAt time.Time `json:"created_at"`
}

type UnassignedStation struct {
	ID         uint16    `json:"id"`
	DetectedAt time.Time `json:"detected_at"`
}

type StationsConfig struct {
	Stations   []Station           `json:"stations"`
	Unassigned []UnassignedStation `json:"unassigned"`
}

var defaultStationsPath = "stations.json"

func GetStationsPath() string {
	path := os.Getenv("STATIONS_PATH")
	if path != "" {
		return path
	}
	return defaultStationsPath
}

func LoadStations(path string) (*StationsConfig, error) {
	if path == "" {
		path = GetStationsPath()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &StationsConfig{
				Stations:   []Station{},
				Unassigned: []UnassignedStation{},
			}, nil
		}
		return nil, err
	}

	var cfg StationsConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func SaveStations(cfg *StationsConfig, path string) error {
	if path == "" {
		path = GetStationsPath()
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func AddStation(cfg *StationsConfig, s Station) error {
	if GetStationByID(cfg, s.ID) != nil {
		return ErrStationExists
	}

	if s.CreatedAt.IsZero() {
		s.CreatedAt = time.Now()
	}
	if s.Outputs == nil {
		s.Outputs = []Output{}
	}

	cfg.Stations = append(cfg.Stations, s)
	return nil
}

func UpdateStation(cfg *StationsConfig, id uint16, s Station) error {
	for i, st := range cfg.Stations {
		if st.ID == id {
			s.ID = id
			s.CreatedAt = st.CreatedAt
			if s.Outputs == nil {
				s.Outputs = st.Outputs
			}
			cfg.Stations[i] = s
			return nil
		}
	}
	return ErrStationNotFound
}

func DeleteStation(cfg *StationsConfig, id uint16) error {
	for i, st := range cfg.Stations {
		if st.ID == id {
			cfg.Stations = append(cfg.Stations[:i], cfg.Stations[i+1:]...)
			return nil
		}
	}
	return ErrStationNotFound
}

func GetStationByID(cfg *StationsConfig, id uint16) *Station {
	for _, st := range cfg.Stations {
		if st.ID == id {
			return &st
		}
	}
	return nil
}

func AddUnassigned(cfg *StationsConfig, id uint16) {
	if IsUnassigned(cfg, id) {
		return
	}
	cfg.Unassigned = append(cfg.Unassigned, UnassignedStation{
		ID:         id,
		DetectedAt: time.Now(),
	})
}

func RemoveUnassigned(cfg *StationsConfig, id uint16) {
	for i, st := range cfg.Unassigned {
		if st.ID == id {
			cfg.Unassigned = append(cfg.Unassigned[:i], cfg.Unassigned[i+1:]...)
			return
		}
	}
}

func IsUnassigned(cfg *StationsConfig, id uint16) bool {
	for _, st := range cfg.Unassigned {
		if st.ID == id {
			return true
		}
	}
	return false
}

var (
	ErrStationNotFound = &StationError{"station not found"}
	ErrStationExists   = &StationError{"station already exists"}
)

type StationError struct {
	msg string
}

func (e *StationError) Error() string {
	return e.msg
}
