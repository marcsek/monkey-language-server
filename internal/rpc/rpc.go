package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

type RequestMessage struct {
	Method string `json:"method"`
}

func EncodeMessage(msg any) string {
	content, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}

func DecodeMessage(msg []byte) (string, []byte, error) {
	header, content, found := bytes.Cut(msg, []byte("\r\n\r\n"))
	if !found {
		return "", nil, fmt.Errorf("Incorrect message format (did not find separator)")
	}

	contentLengthBytes := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return "", nil, fmt.Errorf("Incorrect message format (wrong content length)")
	}

	var requestMessage RequestMessage
	if err := json.Unmarshal(content[:contentLength], &requestMessage); err != nil {
		return "", nil, fmt.Errorf("Incorrect message foramt (can't decode content)")
	}

	return requestMessage.Method, content[:contentLength], nil
}

func SplicFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	header, content, found := bytes.Cut(data, []byte("\r\n\r\n"))
	if !found {
		return 0, nil, nil
	}

	contentLengthBytes := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return 0, nil, err
	}

	if len(content) < contentLength {
		return 0, nil, nil
	}

	totalLength := len(header) + 4 + contentLength
	return totalLength, data[:totalLength], nil
}
