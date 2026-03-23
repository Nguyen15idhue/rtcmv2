package capture

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
)

type Config struct {
	Interface   string
	BPFFilter   string
	FrameChan   chan []byte
	SnapshotLen int32
	Promiscuous bool
	Timeout     time.Duration
}

func DefaultConfig() Config {
	return Config{
		SnapshotLen: 1600,
		Promiscuous: true,
		Timeout:     30 * time.Second,
	}
}

type Capture struct {
	config    Config
	factory   *StreamFactory
	assembler *tcpassembly.Assembler
}

func NewCapture(cfg Config) *Capture {
	factory := NewStreamFactory(cfg.FrameChan)
	pool := tcpassembly.NewStreamPool(factory)
	assembler := tcpassembly.NewAssembler(pool)

	return &Capture{
		config:    cfg,
		factory:   factory,
		assembler: assembler,
	}
}

func (c *Capture) Run(ctx context.Context) error {
	handle, err := pcap.OpenLive(
		c.config.Interface,
		c.config.SnapshotLen,
		c.config.Promiscuous,
		c.config.Timeout,
	)
	if err != nil {
		return fmt.Errorf("open pcap: %w", err)
	}
	defer handle.Close()

	if c.config.BPFFilter != "" {
		if err := handle.SetBPFFilter(c.config.BPFFilter); err != nil {
			return fmt.Errorf("set bpf filter: %w", err)
		}
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packetSource.Lazy = true
	packetSource.NoCopy = true

	log.Printf("capture: started on interface %s, filter: %s", c.config.Interface, c.config.BPFFilter)

	defer func() {
		if r := recover(); r != nil {
			log.Printf("capture: panic recovered: %v\n%s", r, debug.Stack())
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Printf("capture: stopped")
			return nil
		default:
			packet, err := packetSource.NextPacket()
			if err != nil {
				continue
			}
			c.processPacket(packet)
		}
	}
}

func (c *Capture) processPacket(packet gopacket.Packet) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("capture: processPacket panic recovered: %v", r)
		}
	}()

	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer == nil {
		return
	}

	tcp, _ := tcpLayer.(*layers.TCP)
	netFlow := packet.NetworkLayer().NetworkFlow()

	c.assembler.AssembleWithTimestamp(netFlow, tcp, packet.Metadata().Timestamp)
}

func (c *Capture) Flush() {
	c.assembler.FlushOlderThan(time.Now().Add(time.Minute))
}
