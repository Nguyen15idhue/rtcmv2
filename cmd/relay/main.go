package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"rtcmv2/internal/capture"
	"rtcmv2/internal/debug"
	"rtcmv2/internal/relay"
)

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

	d := relay.NewDispatcher(*relayCfg, captureCfg.FrameChan, metrics)
	d.Start(ctx)

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		log.Println("received signal, stopping...")
		cancel()
	}()

	go func() {
		<-ctx.Done()
		d.Stop()
	}()

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
		if err := http.ListenAndServe(debugAddr, nil); err != nil {
			log.Printf("debug server error: %v", err)
		}
	}()
	_ = debugSrv

	log.Println("starting capture...")
	if err := c.Run(ctx); err != nil {
		log.Fatalf("capture error: %v", err)
	}
}
