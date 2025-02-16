package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp4", "127.0.0.1:9022")
	if err != nil {
		fmt.Println("Failed to connect:", err.Error())
		return
	}
	defer conn.Close()

	message := "Hello, server!"
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Failed to send data:", err.Error())
		return
	}

	fmt.Println("Message sent successfully!")
}
