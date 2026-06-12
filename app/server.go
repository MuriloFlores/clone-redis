package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

var memDb = map[string]string{}
var dbMu sync.RWMutex

func MessageReader(conn net.Conn) ([]string, error) {
	reader := bufio.NewReader(conn)

	messageType, err := reader.ReadByte()
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}

		return nil, fmt.Errorf("erro de leitura da mensagem %v", err)

	}

	if messageType != '*' {
		return nil, fmt.Errorf("erro no tipo da mensagem da mensagem %v", err)
	}

	arrLenLine, _ := reader.ReadString('\n')
	arrLen, _ := strconv.Atoi(strings.TrimSpace(arrLenLine))

	commands := make([]string, arrLen)

	for i := 0; i < arrLen; i++ {
		reader.ReadByte()

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

func handleCommand(message []string) string {
	command := strings.ToUpper(message[0])

	if command == "SET" {
		dbMu.Lock()
		memDb[message[1]] = message[2]
		dbMu.Unlock()

		fmt.Println(memDb)

		return "+OK\r\n"
	}

	if command == "GET" {
		dbMu.RLock()
		value, ok := memDb[message[1]]
		dbMu.RUnlock()

		if ok {
			return fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
		}

		return "$-1\r\n"
	}

	return fmt.Sprintf("-ERR unknown command '%s'\r\n", message[0])
}

func handleMessage(conn net.Conn) {
	defer conn.Close()

	for {
		words, err := MessageReader(conn)

		if err != nil {
			fmt.Printf("Error reading: %s\n", err)
		}

		if len(words) == 0 {
			return
		}

		returnMessage := handleCommand(words)

		_, err = conn.Write([]byte(returnMessage))
		if err != nil {
			fmt.Printf("Error writing: %s\n", err)
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
