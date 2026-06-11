package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 1024)
	},
}

var db = map[string]string{}

func MessageReader(conn net.Conn) ([]string, error) {
	reader := bufio.NewReader(conn)

	messageType, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	if messageType != '*' {
		return nil, errors.New("invalid message type")
	}

	arrLenLine, _ := reader.ReadString('\n')
	arrLen, _ := strconv.Atoi(strings.TrimSpace(arrLenLine))

	commands := make([]string, arrLen)

	for i := 0; i < arrLen; i++ {
		reader.ReadByte() //Pulando o byte $

		wordLenLine, _ := reader.ReadString('\n')
		wordLen, _ := strconv.Atoi(strings.TrimSpace(wordLenLine))

		wordBuffer := make([]byte, wordLen)
		io.ReadFull(reader, wordBuffer)

		commands[i] = string(wordBuffer)

		reader.ReadByte()
		reader.ReadByte()
	}

	return commands, nil
}

func handleMessage(conn net.Conn) {
	defer conn.Close()

	bufInterface := bufferPool.Get()
	buffer := bufInterface.([]byte)

	defer bufferPool.Put(bufInterface)

	for {
		_, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Printf("Error reading: %s\n", err)
			continue
		}

		words, err := MessageReader(conn)
		if err != nil {
			fmt.Printf("Error reading: %s\n", err)
		}

		fmt.Println(words)

		_, err = conn.Write([]byte{})
		if err != nil {
			fmt.Printf("Error writing: %s\n", err)
			continue
		}
	}
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

		go handleMessage(conn)
	}
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	handleConnection()
}
