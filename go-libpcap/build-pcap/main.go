package main

import (
	"net"
	"os"
)

func main() {
	pcap := buildPcap()

	file, _ := os.Create("output.pcap")
	defer file.Close()

	file.Write(pcap)
}

func buildPcap() []byte {
	pcapHeader := []byte{
		0xd4, 0xc3, 0xb2, 0xa1, // Magic Number (little-endian)
		0x02, 0x00, 0x04, 0x00, // Major Version 2, Minor Version 4
		0x00, 0x00, 0x00, 0x00, // ThisZone = 0
		0x00, 0x00, 0x00, 0x00, // SigFigs = 0
		0xff, 0xff, 0x00, 0x00, // SnapLen = 65535
		0x01, 0x00, 0x00, 0x00, // LinkType = Ethernet
	}

	// request
	reqPayload := []byte("GET /hello HTTP/1.1\r\nHost: pcap.ran.net\r\n\r\n")
	reqPacket := buildTCPPacketIPv4("192.168.1.1", "192.168.1.2", 12345, 80, 1, 0, reqPayload)
	reqPacketLen := uint32(len(reqPacket))

	reqPcapPacketHeader := []byte{
		0x00, 0x00, 0x00, 0x00, // Timestamp (Seconds)
		0x00, 0x00, 0x00, 0x00, // Timestamp (Microseconds or nanoseconds)
		byte(reqPacketLen), byte(reqPacketLen >> 8), byte(reqPacketLen >> 16), byte(reqPacketLen >> 24), // Captured Packet Length
		byte(reqPacketLen), byte(reqPacketLen >> 8), byte(reqPacketLen >> 16), byte(reqPacketLen >> 24), // Original Packet Length
	}

	// response
	respPayload := []byte("HTTP/1.1 200 OK\r\n" + "Content-Length: 13\r\n" + "\r\n" + "hello, world!")
	respPacket := buildTCPPacketIPv4("192.168.1.2", "192.168.1.1", 80, 12345, 1, 1, respPayload)
	respPacketLen := uint32(len(respPacket))

	respPcapPacketHeader := []byte{
		0x00, 0x00, 0x00, 0x00, // Timestamp (Seconds)
		0x00, 0x00, 0x00, 0x00, // Timestamp (Microseconds or nanoseconds)
		byte(respPacketLen), byte(respPacketLen >> 8), byte(respPacketLen >> 16), byte(respPacketLen >> 24), // Captured Packet Length
		byte(respPacketLen), byte(respPacketLen >> 8), byte(respPacketLen >> 16), byte(respPacketLen >> 24), // Original Packet Length
	}

	// ipv6 req
	ipv6Payload := []byte("GET /hello HTTP/1.1\r\nHost: pcap.ran.net\r\n\r\n")
	ipv6Packet := buildTCPPacketIPv6("2001:db8::2", "2001:db8::1", 12345, 80, 2, 1, ipv6Payload)
	ipv6PacketLen := uint32(len(ipv6Packet))

	ipv6PcapPacketHeader := []byte{
		0x00, 0x00, 0x00, 0x00, // Timestamp (Seconds)
		0x00, 0x00, 0x00, 0x00, // Timestamp (Microseconds or nanoseconds)
		byte(ipv6PacketLen), byte(ipv6PacketLen >> 8), byte(ipv6PacketLen >> 16), byte(ipv6PacketLen >> 24), // Captured Packet Length
		byte(ipv6PacketLen), byte(ipv6PacketLen >> 8), byte(ipv6PacketLen >> 16), byte(ipv6PacketLen >> 24), // Original Packet Length
	}
	return append(pcapHeader,
		append(reqPcapPacketHeader,
			append(reqPacket,
				append(respPcapPacketHeader,
					append(respPacket,
						append(ipv6PcapPacketHeader, ipv6Packet...)...)...)...)...)...)
}

// big-endian
func buildTCPPacketIPv4(srcIP, dstIP string, srcPort, dstPort int, seq, ack int, payload []byte) []byte {
	ethHeader := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // DstMac
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // SrcMac
		0x08, 0x00, // IPv4
	}

	totalLen := uint16(20 + 20 + len(payload)) // ipHeaderLen + tcpHeaderLen + payloadLen

	ipHeader := []byte{
		0x45,                                // Version and IHL
		0x00,                                // TOS
		byte(totalLen >> 8), byte(totalLen), // TL
		0x00, 0x00, // Identification
		0x00, 0x00, // Flags and FO
		0x40,       // TTL = 64
		0x06,       // PROT = TCP
		0x00, 0x00, // Checksum ↓
	}
	ipHeader = append(ipHeader, net.ParseIP(srcIP).To4()...)
	ipHeader = append(ipHeader, net.ParseIP(dstIP).To4()...)
	ipHeader[10], ipHeader[11] = checksum(ipHeader)

	tcpHeader := []byte{
		byte(srcPort >> 8), byte(srcPort),
		byte(dstPort >> 8), byte(dstPort),
		byte(seq >> 24), byte(seq >> 16), byte(seq >> 8), byte(seq),
		byte(ack >> 24), byte(ack >> 16), byte(ack >> 8), byte(ack),
		0x50,       // Data Offset and Deversed
		0x18,       // 00011000 ACK: true, PSH: true
		0xff, 0xff, // Window Size
		0x00, 0x00, // Checksum ↓
		0x00, 0x00, // Urgent Pointer
	}
	// 伪头部
	pseudoHeader := []byte{
		// srcIP
		ipHeader[12], ipHeader[13], ipHeader[14], ipHeader[15],
		// dstIP
		ipHeader[16], ipHeader[17], ipHeader[18], ipHeader[19],
		0x00, 0x06,
		byte((len(tcpHeader) + len(payload)) >> 8),
		byte(len(tcpHeader) + len(payload)),
	}
	tcpChecksumData := append(pseudoHeader, append(tcpHeader, payload...)...)
	tcpHeader[16], tcpHeader[17] = checksum(tcpChecksumData)

	packet := append(ethHeader, append(ipHeader, append(tcpHeader, payload...)...)...)
	return packet
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

func buildTCPPacketIPv6(srcIP, dstIP string, srcPort, dstPort int, seq, ack int, payload []byte) []byte {
	ethHeader := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // DstMac
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // SrcMac
		0x86, 0xdd, // IPv6
	}

	payloadLen := uint16(20 + len(payload)) // TCP Header + Payload

	ipHeader := []byte{
		0x60,       // Version(6)
		0x00,       // Traffic Class
		0x00, 0x00, // Flow Label
		byte(payloadLen >> 8), byte(payloadLen), // Payload Length
		0x06, // Next Header = TCP
		0x40, // Hop Limit = 64
	}
	ipHeader = append(ipHeader, net.ParseIP(srcIP).To16()...)
	ipHeader = append(ipHeader, net.ParseIP(dstIP).To16()...)

	tcpHeader := []byte{
		byte(srcPort >> 8), byte(srcPort),
		byte(dstPort >> 8), byte(dstPort),
		byte(seq >> 24), byte(seq >> 16), byte(seq >> 8), byte(seq),
		byte(ack >> 24), byte(ack >> 16), byte(ack >> 8), byte(ack),
		0x50,       // Data Offset (5) + Reserved (0)
		0x18,       // Flags (ACK + PSH)
		0xff, 0xff, // Window Size
		0x00, 0x00, // Checksum (to be filled)
		0x00, 0x00, // Urgent Pointer
	}

	// 伪头部
	pseudoHeader := make([]byte, 40)
	copy(pseudoHeader[0:16], ipHeader[8:24])   // 源地址
	copy(pseudoHeader[16:32], ipHeader[24:40]) // 目的地址
	pseudoHeader[32] = 0                       // Upper-Layer Packet Length 高字节
	pseudoHeader[33] = 0                       // Upper-Layer Packet Length 低字节
	pseudoHeader[34] = byte((len(tcpHeader) + len(payload)) >> 8)
	pseudoHeader[35] = byte(len(tcpHeader) + len(payload))
	pseudoHeader[39] = 0x06 // Next Header (TCP)

	tcpChecksumData := append(pseudoHeader, append(tcpHeader, payload...)...)
	tcpHeader[16], tcpHeader[17] = checksum(tcpChecksumData)

	packet := append(ethHeader, append(ipHeader, append(tcpHeader, payload...)...)...)
	return packet
}
