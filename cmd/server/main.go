package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Nguyen15idhue/rtcmv2/api"
	"github.com/Nguyen15idhue/rtcmv2/internal/capture"
	"github.com/Nguyen15idhue/rtcmv2/internal/relay"
)

var (
	flagConfig   = flag.String("config", "config.json", "Path to config file")
	flagAPI      = flag.String("api", ":1507", "API server address")
	flagDemo     = flag.Bool("demo", false, "Run in demo mode with simulated data")
	flagDemoPort = flag.Int("demo-port", 2101, "Demo data port")
)

const shutdownTimeout = 10 * time.Second

func main() {
	flag.Parse()

	relayCfg, err := relay.LoadConfig(*flagConfig)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	metrics := relay.NewMetrics()

	var dispatcher *relay.Dispatcher

	if *flagDemo {
		log.Println("Starting in demo mode...")
		dispatcher = startDemoMode(ctx, metrics, *flagDemoPort)
	} else {
		captureCfg := capture.DefaultConfig()
		captureCfg.Interface = relayCfg.Capture.Interface
		captureCfg.BPFFilter = relayCfg.Capture.Filter
		captureCfg.FrameChan = make(chan []byte, 100)

		dispatcher = relay.NewDispatcher(*relayCfg, captureCfg.FrameChan, metrics)
		dispatcher.Start(ctx)

		c := capture.NewCapture(captureCfg)
		go func() {
			if err := c.Run(ctx); err != nil {
				log.Printf("capture error: %v", err)
			}
		}()
	}

	apiServer := api.NewServer(*flagAPI, metrics, dispatcher)
	go func() {
		log.Printf("Dashboard: http://localhost%s/", *flagAPI)
		if err := apiServer.Start(); err != nil {
			log.Printf("api server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("RTCM Relay running. Press Ctrl+C to exit.")

	select {
	case sig := <-sigChan:
		log.Printf("received %s, shutting down...", sig)
		cancel()

		done := make(chan struct{})
		go func() {
			if dispatcher != nil {
				dispatcher.Stop()
			}
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

func startDemoMode(ctx context.Context, metrics *relay.Metrics, port int) *relay.Dispatcher {
	cfg := relay.Config{
		Casters: []relay.CasterConfig{},
	}

	frameChan := make(chan []byte, 100)
	dispatcher := relay.NewDispatcher(cfg, frameChan, metrics)
	dispatcher.Start(ctx)

	go func() {
		stationIDs := []uint16{1000, 2000, 3000}
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(500 * time.Millisecond):
			}
			for _, stationID := range stationIDs {
				metrics.RecordFrame(stationID, "demo-"+itoa(int(stationID)))
			}
		}
	}()

	return dispatcher
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
