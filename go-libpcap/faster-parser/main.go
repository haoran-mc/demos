package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	device       string = "en0"
	snapshot_len int32  = 1024
	promiscuous  bool   = false
	err          error
	timeout      time.Duration = 30 * time.Second
	handle       *pcap.Handle
	// Will reuse these for each packet
	ethLayer layers.Ethernet
	ipLayer  layers.IPv4
	tcpLayer layers.TCP
)

func main() {
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {

		// 为指定的层创建解析器，这样就可以用已有的结构来存储 packet 信息，
		// 而不是为每个 packet 创建新的结构，既浪费内存又浪费时间
		// ↑ ethLayer
		// ↑ ipLayer
		// ↑ tcpLayer
		parser := gopacket.NewDecodingLayerParser(
			layers.LayerTypeEthernet,
			&ethLayer,
			&ipLayer,
			&tcpLayer,
		)

		foundLayerTypes := []gopacket.LayerType{}
		if err = parser.DecodeLayers(packet.Data(), &foundLayerTypes); err != nil {
			fmt.Println("Trouble decoding layers: ", err)
		}

		for _, layerType := range foundLayerTypes {
			switch layerType {
			case layers.LayerTypeEthernet:
				fmt.Println("Eth MAC:", ethLayer.SrcMAC, "-> ", ethLayer.DstMAC)
			case layers.LayerTypeIPv4:
				fmt.Println("IPv4:", ipLayer.SrcIP, "->", ipLayer.DstIP)
			case layers.LayerTypeTCP:
				fmt.Println("TCP Port:", tcpLayer.SrcPort, "->", tcpLayer.DstPort)
				fmt.Println("TCP SYN:", tcpLayer.SYN, " | ACK:", tcpLayer.ACK)
			}
		}
		fmt.Println()
	}
}
