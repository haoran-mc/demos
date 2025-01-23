package main

import (
	"fmt"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
)

var (
	deviceName  string = "en0"
	snapshotLen uint32 = 1024
	promiscuous bool   = false
	err         error
	timeout     time.Duration = -1 * time.Second
	handle      *pcap.Handle
	packetCount int = 0
)

func main() {
	f, _ := os.Create("test.pcap")

	// w 是 pcapgo.Writer 实例，用于写入 pcap 格式的数据
	w := pcapgo.NewWriter(f)

	// 写入 pcap 文件头，指定最大数据包长度 1024 字节，指定链路层类型为以太网
	w.WriteFileHeader(snapshotLen, layers.LinkTypeEthernet)
	defer f.Close()

	// Open the device for capturing
	handle, err = pcap.OpenLive(deviceName, int32(snapshotLen), promiscuous, timeout)
	if err != nil {
		fmt.Printf("Error opening device %s: %v", deviceName, err)
		os.Exit(1)
	}
	defer handle.Close()

	// 创建数据包源，返回一个可以读取网络数据包的源
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	// packetSource.Packets 返回一个阻塞式的 channel，持续提供新捕获的数据包
	for packet := range packetSource.Packets() {
		// Process packet here
		fmt.Println(packet)

		w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
		packetCount++

		// 捕获 100 个数据包后停止
		if packetCount > 100 {
			break
		}
	}
}
