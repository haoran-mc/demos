package main

import (
	"log"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	device       string = "en0"
	snapshot_len int32  = 1024
	promiscuous  bool   = false
	err          error
	handle       *pcap.Handle
	srcIP        net.IP         = net.IPv4(192, 168, 1, 100) // 替换为你的源 IP
	dstIP        net.IP         = net.IPv4(192, 168, 1, 101) // 替换为目标 IP
	srcPort      layers.TCPPort = layers.TCPPort(12345)      // 源端口
	dstPort      layers.TCPPort = layers.TCPPort(80)         // 目标端口
	payload      []byte         = []byte("Hello, TCP!")      // TCP 数据负载
)

func main() {
	// 获取设备的网络信息
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, pcap.BlockForever)
	if err != nil {
		log.Fatalf("打开设备失败: %v", err)
	}
	defer handle.Close()

	sendRawBytes() // DecodeFailure: Packet decoding error: Ethernet packet too small
	sendEmptyTCP() // DecodeFailure: Packet decoding error: Invalid TCP data offset 0 < 5
	sendTCP()      // successfully send!
}

func sendRawBytes() {
	rawBytes := []byte{10, 20, 30}
	if err := handle.WritePacketData(rawBytes); err != nil {
		log.Printf("send raw bytes fail: %v", err)
	} else {
		log.Printf("send raw bytes succ.")
	}
}

func sendEmptyTCP() {
	buffer := gopacket.NewSerializeBuffer()
	gopacket.SerializeLayers(buffer, gopacket.SerializeOptions{},
		&layers.Ethernet{},
		&layers.IPv4{},
		&layers.TCP{},
		gopacket.Payload([]byte("Empty TCP.")),
	)

	outgoingPacket := buffer.Bytes()
	if err := handle.WritePacketData(outgoingPacket); err != nil {
		log.Printf("send empty tcp fail: %v", err)
	} else {
		log.Printf("send empty tcp succ.")
	}
}

func sendTCP() {
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
		SrcIP:    srcIP,
		DstIP:    dstIP,
	}

	// 构造 TCP 层（第四层）
	tcpLayer := &layers.TCP{
		SrcPort: srcPort,
		DstPort: dstPort,
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

	// 发送数据包
	outgoingPacket := buffer.Bytes()
	if err := handle.WritePacketData(outgoingPacket); err != nil {
		log.Printf("send TCP fail: %v", err)
	} else {
		log.Printf("send TCP succ.")
	}
}
