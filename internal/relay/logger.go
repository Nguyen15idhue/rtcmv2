package relay

import (
	"log/slog"
	"os"
	"runtime/debug"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelInfo,
}))

type LogFields struct {
	StationID  uint16
	Caster     string
	Mountpoint string
	Name       string
	Event      string
	Error      string
	Frames     int64
	Duration   string
}

func logInfo(msg string, fields LogFields) {
	args := []any{}
	if fields.StationID > 0 {
		args = append(args, "station_id", fields.StationID)
	}
	if fields.Caster != "" {
		args = append(args, "caster", fields.Caster)
	}
	if fields.Mountpoint != "" {
		args = append(args, "mountpoint", fields.Mountpoint)
	}
	if fields.Name != "" {
		args = append(args, "name", fields.Name)
	}
	if fields.Event != "" {
		args = append(args, "event", fields.Event)
	}
	if fields.Error != "" {
		args = append(args, "error", fields.Error)
	}
	if fields.Frames > 0 {
		args = append(args, "frames", fields.Frames)
	}
	logger.Info(msg, args...)
}

func logError(msg string, fields LogFields) {
	args := []any{}
	if fields.StationID > 0 {
		args = append(args, "station_id", fields.StationID)
	}
	if fields.Caster != "" {
		args = append(args, "caster", fields.Caster)
	}
	if fields.Mountpoint != "" {
		args = append(args, "mountpoint", fields.Mountpoint)
	}
	if fields.Name != "" {
		args = append(args, "name", fields.Name)
	}
	if fields.Event != "" {
		args = append(args, "event", fields.Event)
	}
	if fields.Error != "" {
		args = append(args, "error", fields.Error)
	}
	logger.Error(msg, args...)
}

func logWarn(msg string, fields LogFields) {
	args := []any{}
	if fields.StationID > 0 {
		args = append(args, "station_id", fields.StationID)
	}
	if fields.Name != "" {
		args = append(args, "name", fields.Name)
	}
	if fields.Event != "" {
		args = append(args, "event", fields.Event)
	}
	if fields.Error != "" {
		args = append(args, "error", fields.Error)
	}
	logger.Warn(msg, args...)
}

func safeGo(name string, fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logError("panic_recovered", LogFields{
					Name:  name,
					Error: string(debug.Stack()),
				})
			}
		}()
		fn()
	}()
}
