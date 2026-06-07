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

var db = map[string]string{}

func handleMessage(conn net.Conn, message string) {
	defer conn.Close()

	bufInterface := bufferPool.Get()
	buffer := bufInterface.([]byte)

	defer bufferPool.Put(bufInterface)

	for {
		byteLen, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Printf("Error reading: %s\n", err)
			continue
		}

		_, err = conn.Write(cleanMessage(buffer[:byteLen]))
		if err != nil {
			fmt.Printf("Error writing: %s\n", err)
			continue
		}
	}
}

func cleanMessage(message []byte) []byte {
	totalWords := strings.Split(string(message), " ")

	if strings.TrimSpace(strings.ToUpper(totalWords[0])) == "SET" {
		db[totalWords[1]] = totalWords[2]

		fmt.Println(db)

		return []byte(fmt.Sprintf("$+OK\r\n"))
	}

	if strings.TrimSpace(strings.ToUpper(totalWords[0])) == "GET" {
		for key, value := range db {
			if key == totalWords[1] {
				fmt.Println(db)

				return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value))
			}
		}
	}

	return []byte("$-1\r\n")
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
