package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/marcsek/monkey-language-server/internal/analysis"
	"github.com/marcsek/monkey-language-server/internal/lsp"
	"github.com/marcsek/monkey-language-server/internal/rpc"
)

const LOG_FILE_PATH = "/home/marek/Personal/code/compiler_go/monkey-language-server/log.txt"

func main() {
	logger := getLogger(LOG_FILE_PATH)
	logger.Println("Logger started!")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.SplicFunc)

	state := analysis.NewState()
	writer := os.Stdout

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Got and error: %s\n", err)
		}

		handleMessage(logger, writer, state, method, contents)
	}
}

func handleMessage(
	logger *log.Logger,
	writer io.Writer,
	state *analysis.State,
	method string,
	contents []byte,
) {
	logger.Printf("Received %s", method)
	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Couldn't parse message (%s) in \"%s\"\n", err, method)
		}
		logger.Printf(
			"Connected to: %s %s",
			request.Params.ClientInfo.Name,
			request.Params.ClientInfo.Version,
		)

		msg := lsp.NewInitializeResponse(request.ID)
		writeResponse(writer, msg)

		logger.Println("Sent the reply")

	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Couldn't parse message (%s) in \"%s\"\n", err, method)
		}
		logger.Printf(
			"Opened: %s %s",
			request.Params.TextDocument.LanguageID,
			request.Params.TextDocument.URI,
		)

		state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)

	case "textDocument/didChange":
		var request lsp.TextDocumentDidChangeNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Couldn't parse message (%s) in \"%s\"\n", err, method)
		}
		logger.Printf(
			"Updated: %d %s",
			request.Params.TextDocument.Version,
			request.Params.TextDocument.URI,
		)

		for _, change := range request.Params.ContentChanges {
			state.UpdateDocument(request.Params.TextDocument.URI, change.Text)
		}

	case "textDocument/hover":
		var request lsp.HoverRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Couldn't parse message (%s) in \"%s\"\n", err, method)
		}
		logger.Printf(
			"Hovered : %d %d %s",
			request.Params.Position.Line,
			request.Params.Position.Character,
			request.Params.TextDocument.URI,
		)

		response := state.Hover(
			request.ID,
			request.Params.TextDocument.URI,
			request.Params.Position,
		)
		writeResponse(writer, response)

	case "textDocument/definition":
		var request lsp.DefinitionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Couldn't parse message (%s) in \"%s\"\n", err, method)
		}
		logger.Printf(
			"Definition: %d %d %s",
			request.Params.Position.Line,
			request.Params.Position.Character,
			request.Params.TextDocument.URI,
		)

		response := state.Definition(
			request.ID,
			request.Params.TextDocument.URI,
			request.Params.Position,
		)
		writeResponse(writer, response)

	case "textDocument/codeAction":
		var request lsp.CodeActionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Couldn't parse message (%s) in \"%s\"\n", err, method)
		}
		logger.Printf("Code action: %s", request.Params.TextDocument.URI)

		response := state.TextDocumentCodeAction(request.ID, request.Params.TextDocument.URI)
		writeResponse(writer, response)
	}
}

func writeResponse(writer io.Writer, msg any) {
	reply := rpc.EncodeMessage(msg)
	writer.Write([]byte(reply))
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("wrong log file")
	}

	return log.New(logfile, "[monkey-language-server]", log.Ldate|log.Ltime|log.Lshortfile)
}
