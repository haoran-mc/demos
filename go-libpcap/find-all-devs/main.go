package main

import (
	"fmt"
	"log"
	"net"

	"github.com/google/gopacket/pcap"
)

func main() {
	// 得到所有的(网络)设备
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	// 打印设备信息
	fmt.Println("Devices found:")
	for _, device := range devices {
		fmt.Println("\nName: ", device.Name)
		fmt.Println("Description: ", device.Description)
		fmt.Println("Devices addresses: ", device.Description)
		for _, address := range device.Addresses {
			fmt.Println("- IP address: ", address.IP)
			fmt.Println("- Subnet mask: ", address.Netmask)
		}
	}

	// 使用内置 net 包可以查看网卡 UP DOWN 状态
	nics, err := net.Interfaces()
	if err != nil {
		log.Fatal("failed to find all network interfaces: " + err.Error())
	}
	for _, nic := range nics {
		fmt.Printf("\nName: %s\n", nic.Name)
		fmt.Printf("Hardware Address: %s\n", nic.HardwareAddr)
		fmt.Printf("MTU: %d\n", nic.MTU)
		fmt.Printf("Flags: %s\n", nic.Flags)
		if nic.Flags&net.FlagUp != 0 {
			fmt.Println("Status: UP")
		} else {
			fmt.Println("Status: DOWN")
		}
	}
}
