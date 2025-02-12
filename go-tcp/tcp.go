package main

import (
	"log"
	"net"
	"syscall"
)

func main() {
	payload := []byte("GET /hello HTTP/1.1\r\nHost: tcp.net\r\n\r\n")
	srcIP := net.ParseIP("192.168.1.1").To4()
	dstIP := net.ParseIP("192.168.1.2").To4()
	srcPort, dstPort := 12345, 9022
	seq, ack := 1000, 2000

	tcpHeader := []byte{
		byte(srcPort >> 8), byte(srcPort),
		byte(dstPort >> 8), byte(dstPort),
		byte(seq >> 24), byte(seq >> 16), byte(seq >> 8), byte(seq),
		byte(ack >> 24), byte(ack >> 16), byte(ack >> 8), byte(ack),
		0x50,       // Data Offset and Deversed
		0x18,       // 00011000 ACK: true, PSH: true
		0xff, 0xff, // Window Size
		0x00, 0x00, // Checksum set to zero at start
		0x00, 0x00, // Urgent Pointer
	}
	// 伪头部
	pseudoHeader := []byte{
		srcIP[0], srcIP[1], srcIP[2], srcIP[3],
		dstIP[0], dstIP[1], dstIP[2], dstIP[3],
		0x00, 0x06,
		byte((len(tcpHeader) + len(payload)) >> 8),
		byte(len(tcpHeader) + len(payload)),
	}
	tcpChecksumData := append(pseudoHeader, append(tcpHeader, payload...)...)
	tcpHeader[16], tcpHeader[17] = checksum(tcpChecksumData)
	// fmt.Printf("%x %x\n", tcpHeader[16], tcpHeader[17])

	packet := append(tcpHeader, payload...)

	err := sendRawPacket(packet, "127.0.0.1", 80)
	if err != nil {
		log.Fatalf("Failed to send packet: %v", err)
	}
}

func checksum(data []byte) (byte, byte) {
	sum := uint32(0)
	for i := 0; i+1 < len(data); i += 2 {
		sum += uint32(data[i])<<8 | uint32(data[i+1])
	}
	if len(data)%2 == 1 {
		sum += uint32(data[len(data)-1]) << 8
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)
	return byte(^sum >> 8), byte(^sum)
}

func sendRawPacket(packet []byte, dstIP string, dstPort int) error {
	// create a raw socket
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_TCP)
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	addr := &syscall.SockaddrInet4{}
	addr.Port = dstPort
	copy(addr.Addr[:], net.ParseIP(dstIP).To4())

	// send
	err = syscall.Sendto(fd, packet, 0, addr)
	if err != nil {
		return err
	}

	log.Println("Packet sent successfully")
	return nil
}
