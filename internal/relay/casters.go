package relay

import (
	"encoding/json"
	"os"
)

type Caster struct {
	ID       uint16 `json:"id"`
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
}

type CastersConfig struct {
	Casters []Caster `json:"casters"`
}

var defaultCastersPath = "casters.json"

func GetCastersPath() string {
	path := os.Getenv("CASTERS_PATH")
	if path != "" {
		return path
	}
	return defaultCastersPath
}

func LoadCasters(path string) (*CastersConfig, error) {
	if path == "" {
		path = GetCastersPath()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &CastersConfig{
				Casters: []Caster{},
			}, nil
		}
		return nil, err
	}

	var cfg CastersConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func SaveCasters(cfg *CastersConfig, path string) error {
	if path == "" {
		path = GetCastersPath()
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func AddCaster(cfg *CastersConfig, c Caster) error {
	if GetCasterByID(cfg, c.ID) != nil {
		return ErrCasterExists
	}
	cfg.Casters = append(cfg.Casters, c)
	return nil
}

func DeleteCaster(cfg *CastersConfig, id uint16) error {
	for i, c := range cfg.Casters {
		if c.ID == id {
			cfg.Casters = append(cfg.Casters[:i], cfg.Casters[i+1:]...)
			return nil
		}
	}
	return ErrCasterNotFound
}

func GetCasterByID(cfg *CastersConfig, id uint16) *Caster {
	for _, c := range cfg.Casters {
		if c.ID == id {
			return &c
		}
	}
	return nil
}

var (
	ErrCasterNotFound = &CasterError{"caster not found"}
	ErrCasterExists   = &CasterError{"caster already exists"}
)

type CasterError struct {
	msg string
}

func (e *CasterError) Error() string {
	return e.msg
}
