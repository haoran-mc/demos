package main

import (
	"fmt"
	"log"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func main() {
	payload := []byte{2, 4, 6}
	options := gopacket.SerializeOptions{}
	buffer := gopacket.NewSerializeBuffer()
	gopacket.SerializeLayers(buffer, options,
		&layers.Ethernet{},
		&layers.IPv4{},
		&layers.TCP{},
		gopacket.Payload(payload),
	)
	rawBytes := buffer.Bytes()

	ethPacket :=
		gopacket.NewPacket(
			rawBytes,
			layers.LayerTypeEthernet,
			gopacket.Default,
		)

	ipPacket :=
		gopacket.NewPacket(
			rawBytes,
			layers.LayerTypeIPv4,
			gopacket.Lazy,
		)

	tcpPacket :=
		gopacket.NewPacket(
			rawBytes,
			layers.LayerTypeTCP,
			gopacket.NoCopy,
		)

	fmt.Println(ethPacket)
	fmt.Println(ipPacket)
	fmt.Println(tcpPacket)

	pkt := createFullPacket() // 构造一个完整的 Packet
	fmt.Println(pkt)
	fmt.Println(forgeReset(pkt)) // 操作 Packet
}

func createFullPacket() gopacket.Packet {
	// 构造 Ethernet 层（数据链路层，网络七层模型中的第二层）
	ethLayer := &layers.Ethernet{
		SrcMAC:       net.HardwareAddr{0x00, 0x0c, 0x29, 0x2d, 0x8d, 0x21}, // 源网卡 MAC 地址
		DstMAC:       net.HardwareAddr{0x00, 0x50, 0x56, 0xfc, 0x00, 0x01}, // 目标网卡 MAC 地址
		EthernetType: layers.EthernetTypeIPv4,
	}

	// 构造 IPv4 层（第三层）
	ipLayer := &layers.IPv4{
		Version:  4,
		IHL:      5,
		TOS:      0,
		TTL:      64,
		Protocol: layers.IPProtocolTCP,
		SrcIP:    net.IPv4(192, 168, 1, 100),
		DstIP:    net.IPv4(192, 168, 1, 101),
	}

	// 构造 TCP 层（第四层）
	tcpLayer := &layers.TCP{
		SrcPort: layers.TCPPort(12345),
		DstPort: layers.TCPPort(80),
		Seq:     1105024978,
		Ack:     0,
		Window:  1500,
		Options: []layers.TCPOption{
			{OptionType: layers.TCPOptionKindMSS, OptionLength: 4, OptionData: []byte{0x05, 0xb4}},
		},
	}
	tcpLayer.SetNetworkLayerForChecksum(ipLayer)

	// 构造数据包
	buffer := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	if err := gopacket.SerializeLayers(buffer, options,
		ethLayer,
		ipLayer,
		tcpLayer,
		gopacket.Payload([]byte("Hello, TCP!")),
	); err != nil {
		log.Printf("序列化失败: %v", err)
	}

	newPacket := gopacket.NewPacket(buffer.Bytes(), layers.LayerTypeEthernet, gopacket.Default)
	return newPacket
}

func forgeReset(packet gopacket.Packet) gopacket.Packet {
	ethLayer := packet.Layer(layers.LayerTypeEthernet)
	eth, _ := ethLayer.(*layers.Ethernet)

	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	ip, _ := ipLayer.(*layers.IPv4)

	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	tcp, _ := tcpLayer.(*layers.TCP)

	options := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}

	tcp.RST = true
	tcp.NS = false
	tcp.CWR = false
	tcp.ECE = false
	tcp.URG = false
	tcp.PSH = false
	tcp.SYN = false
	tcp.FIN = false

	if tcp.ACK {
		tcp.ACK = false
		tcp.Seq = tcp.Ack
		tcp.Ack = 0
	} else {
		tcp.ACK = true
		tcp.Ack = tcp.Seq + uint32(len(packet.Data()))
		tcp.Seq = 0
	}

	tcp.Window = 0
	tcp.Urgent = 0
	tcp.Options = tcp.Options[:0]
	tcp.Payload = tcp.Payload[:0]

	tcp.SrcPort, tcp.DstPort = tcp.DstPort, tcp.SrcPort
	ip.SrcIP, ip.DstIP = ip.DstIP, ip.SrcIP
	eth.SrcMAC, eth.DstMAC = eth.DstMAC, eth.SrcMAC

	tcp.SetNetworkLayerForChecksum(ip)

	resetPacketBuffer := gopacket.NewSerializeBuffer()
	err := gopacket.SerializePacket(resetPacketBuffer, options, packet)
	if err != nil {
		panic(err)
	}
	resetPacket := gopacket.NewPacket(resetPacketBuffer.Bytes(), layers.LayerTypeEthernet, gopacket.Default)
	return resetPacket
}
