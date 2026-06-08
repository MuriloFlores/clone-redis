package main

import (
	"bufio"
	"errors"
	"io"
	"net"
	"strconv"
	"strings"
)

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
