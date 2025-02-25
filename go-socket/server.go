package main

import (
	"fmt"
	"log"
	"net"
)

const BUF_SIZE int = 1024

func main() {
	listener, err := net.Listen("tcp4", "127.0.0.1:9022")
	if err != nil {
		log.Fatal("Fail to listen: ", err.Error())
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Fail to accept connection:", err.Error())
			continue
		}

		// 监听到一个客户连接，处理
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, BUF_SIZE)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Fail to read data:", err.Error())
		return
	}

	fmt.Println("Data:", string(buf[:n]))

	// HTTP 响应，响应可有可无
	// 而且因为 client.go 没有接收响应的能力，所以收不到响应
	// 使用 curl 127.0.0.1:9022 能够收到响应
	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 13\r\n" +
		"Connection: close\r\n" +
		"\r\n" +
		"Hello, world!"

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Fail to send response:", err.Error())
	}
}
