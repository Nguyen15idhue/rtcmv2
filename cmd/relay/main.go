package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"rtcmv2/internal/capture"
	"rtcmv2/internal/debug"
	"rtcmv2/internal/relay"
)

const shutdownTimeout = 10 * time.Second

func main() {
	configPath := "config.json"
	if len(os.Args) > 1 {
		configPath = os.Args[2]
	}

	relayCfg, err := relay.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	captureCfg := capture.DefaultConfig()
	captureCfg.Interface = relayCfg.Capture.Interface
	captureCfg.BPFFilter = relayCfg.Capture.Filter
	captureCfg.FrameChan = make(chan []byte, 100)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	metrics := relay.NewMetrics()

	c := capture.NewCapture(captureCfg)

	_ = c

	d := relay.NewDispatcher(*relayCfg, captureCfg.FrameChan, metrics)
	d.Start(ctx)

	debugAddr := ":8080"
	if os.Getenv("DEBUG_ADDR") != "" {
		debugAddr = os.Getenv("DEBUG_ADDR")
	}

	debugSrv := debug.NewServer(metrics, debugAddr)
	go func() {
		log.Printf("debug server listening on %s", debugAddr)
		log.Printf("endpoints:")
		log.Printf("  GET /healthz")
		log.Printf("  GET /debug/stations")
		log.Printf("  GET /debug/metrics")
		if err := http.ListenAndServe(debugAddr, nil); err != nil && err != http.ErrServerClosed {
			log.Printf("debug server error: %v", err)
		}
	}()
	_ = debugSrv

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("starting capture...")

	select {
	case sig := <-sigChan:
		log.Printf("received %s, shutting down...", sig)
		cancel()

		done := make(chan struct{})
		go func() {
			close(captureCfg.FrameChan)
			d.Stop()
			close(done)
		}()

		select {
		case <-done:
			log.Println("shutdown complete")
		case <-time.After(shutdownTimeout):
			log.Printf("shutdown timeout (%v), forcing exit", shutdownTimeout)
		}
	}
}
