package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"time"
)

func main() {
	listen, _ := net.Listen("tcp", ":9520")
	defer listen.Close()

	for {
		conn, _ := listen.Accept()
		request(conn)
		response(conn)
		conn.Close()
	}
}

// print conn info
func request(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		ln := scanner.Text()
		fmt.Println(ln)
		if ln == "" {
			break
		}
	}
}

func response(conn net.Conn) {
	msg := time.Now().Format("2006-01-02 15:04:05") + " OK"

	conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
	conn.Write([]byte("Content-Length: " + strconv.Itoa(len(msg)) + "\r\n"))
	conn.Write([]byte("Content-Type:text/html:charset=UTF-8\r\n"))
	conn.Write([]byte("\r\n"))
	conn.Write([]byte(msg))
}
