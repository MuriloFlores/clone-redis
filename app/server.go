package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 1024)
	},
}

func handleMessage(conn net.Conn, message string) {
	defer conn.Close()

	bufInterface := bufferPool.Get()
	buffer := bufInterface.([]byte)

	defer bufferPool.Put(bufInterface)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Printf("Error reading: %s\n", err)
			continue
		}

		_, err = conn.Write(cleanMessage(buffer[:n]))
		if err != nil {
			fmt.Printf("Error writing: %s\n", err)
			continue
		}
	}
}

func cleanMessage(message []byte) []byte {
	if strings.Contains(strings.TrimSpace(string(message)), "PING") {
		return []byte("PONG " + strings.Split(string(message), " ")[1])
	}

	return []byte("")
}

func handleConnection() net.Conn {

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {

		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleMessage(conn, "+PONG\r\n")
	}
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	handleConnection()

}
