package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type BaseMessage struct {
	Method string
}

func EncodeMessage(msg interface{}) string {
	content, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}

func DecodeMessage(msg []byte) (string, []byte, error) {
	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return "", nil, errors.New("Did not find separator")
	}

	contentLenBytes := header[len("Content-Length: "):]
	contentLen, err := strconv.Atoi(string(contentLenBytes))
	if err != nil {
		return "", nil, err
	}

	var baseMsg BaseMessage
	if err := json.Unmarshal(content[:contentLen], &baseMsg); err != nil {
		return "", nil, err
	}

	return baseMsg.Method, content[:contentLen], nil
}

func Split(data []byte, _ bool) (advance int, token []byte, err error) {
	header, content, found := bytes.Cut(data, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return 0, nil, nil
	}

	contentLenBytes := header[len("Content-Length: "):]
	contentLen, err := strconv.Atoi(string(contentLenBytes))
	if err != nil {
		return 0, nil, err
	}

	if len(content) < contentLen {
		return 0, nil, nil
	}

	totalLen := len(header) + 4 + contentLen

	return totalLen, data[:totalLen], nil
}
