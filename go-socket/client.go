package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9522")
	if err != nil {
		fmt.Println("Failed to connect:", err.Error())
		return
	}
	defer conn.Close()

	// 可以发送非 http 结构的 message
	request := "GET / HTTP/1.1\r\n" +
		"Host: 127.0.0.1:9523\r\n" +
		"Content-Length: 0\r\n" +
		"\r\n"

	_, err = conn.Write([]byte(request))
	if err != nil {
		fmt.Println("Failed to send data:", err.Error())
		return
	}

	fmt.Println("Message sent successfully!")
}
