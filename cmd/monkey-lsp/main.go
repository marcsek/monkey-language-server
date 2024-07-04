package main

import (
	"bufio"
	"log"
	"os"

	"github.com/marcsek/monkey-language-server/internal/analysis"
	"github.com/marcsek/monkey-language-server/internal/messageHandler"
	"github.com/marcsek/monkey-language-server/internal/rpc"
)

const LOG_FILE_PATH = "/home/marek/Personal/code/compiler_go/monkey-language-server/log.txt"

func main() {
	logger := getLogger(LOG_FILE_PATH)
	logger.Println("Logger started!")

	reader := os.Stdin
	writer := os.Stdout

	scanner := bufio.NewScanner(reader)
	scanner.Split(rpc.SplicFunc)

	state := analysis.NewState()

	messageHandler := messageHandler.New(reader, writer, state, logger)

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Got and error: %s\n", err)
		}

		messageHandler.HandleMessage(method, contents)
	}
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("wrong log file")
	}

	return log.New(logfile, "[monkey-language-server]", log.Ldate|log.Ltime|log.Lshortfile)
}
