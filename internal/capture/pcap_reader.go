package capture

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type PcapReader struct {
	filename string
}

func NewPcapReader(filename string) *PcapReader {
	return &PcapReader{filename: filename}
}

type TCPPacket struct {
	SrcIP     string
	DstIP     string
	SrcPort   uint16
	DstPort   uint16
	Payload   []byte
	Timestamp int64
}

func (p *PcapReader) Read() ([]TCPPacket, error) {
	handle, err := pcap.OpenOffline(p.filename)
	if err != nil {
		return nil, err
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	var packets []TCPPacket
	for packet := range packetSource.Packets() {
		if ipLayer := packet.NetworkLayer(); ipLayer != nil {
			ip := ipLayer.NetworkFlow()
			srcIP := ip.Src().String()
			dstIP := ip.Dst().String()

			if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
				tcp, _ := tcpLayer.(*layers.TCP)
				if tcp != nil {
					pkt := TCPPacket{
						SrcIP:     srcIP,
						DstIP:     dstIP,
						SrcPort:   uint16(tcp.SrcPort),
						DstPort:   uint16(tcp.DstPort),
						Payload:   tcp.Payload,
						Timestamp: packet.Metadata().Timestamp.UnixNano(),
					}
					packets = append(packets, pkt)
				}
			}
		}
	}

	return packets, nil
}

func FilterByPort(packets []TCPPacket, port uint16) []TCPPacket {
	var result []TCPPacket
	for _, p := range packets {
		if p.SrcPort == port || p.DstPort == port {
			result = append(result, p)
		}
	}
	return result
}

func GetRTCMPayloads(packets []TCPPacket) [][]byte {
	var payloads [][]byte
	for _, p := range packets {
		if len(p.Payload) > 0 && p.Payload[0] == 0xD3 {
			payloads = append(payloads, p.Payload)
		}
	}
	return payloads
}
