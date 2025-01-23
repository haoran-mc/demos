package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

var (
	device       string = "en0" // 网络接口名称
	snapshot_len int32  = 1024  // 每个数据包捕获的最大字节数
	promiscuous  bool   = false // 是否设置混杂模式
	err          error
	timeout      time.Duration = 30 * time.Second // 数据包捕获超时时间
	handle       *pcap.Handle                     // TODO pcap 句柄 ?
)

func main() {
	// 打开网络接口
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// 创建一个数据包源（packet source），进入一个循环，持续从数据包源读取数据包
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// 只是简单地打印每个捕获到的数据包
		fmt.Println(packet)
	}
}
