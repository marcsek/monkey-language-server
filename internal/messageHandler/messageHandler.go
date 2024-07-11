package messageHandler

import (
	"encoding/json"
	"io"
	"log"

	"github.com/marcsek/monkey-language-server/internal/analysis"
	"github.com/marcsek/monkey-language-server/internal/lsp"
	"github.com/marcsek/monkey-language-server/internal/rpc"
)

type MessageHandler struct {
	reader io.Reader
	writer io.Writer
	state  *analysis.State
	logger *log.Logger
}

func New(
	reader io.Reader,
	writer io.Writer,
	state *analysis.State,
	logger *log.Logger,
) *MessageHandler {
	return &MessageHandler{
		reader: reader,
		writer: writer,
		state:  state,
		logger: logger,
	}
}

func (mh *MessageHandler) HandleMessage(method string, contents []byte) {
	mh.logger.Printf("Received %s", method)

	switch method {
	case "initialize":
		request := parseMessage[lsp.InitializeRequest](contents, mh.logger, method)

		msg := lsp.NewInitializeResponse(request.ID)
		mh.sendMessage(msg)

	case "textDocument/didOpen":
		request := parseMessage[lsp.DidOpenTextDocumentNotification](contents, mh.logger, method)

		diagnostics := mh.state.OpenDocument(
			request.Params.TextDocument.URI,
			request.Params.TextDocument.Text,
		)
		mh.sendMessage(lsp.PublishDiagnosticsNotification{
			Notification: lsp.Notification{
				RPC:    "2.0",
				Method: "textDocument/publishDiagnostics",
			},
			Params: lsp.PublishDiagnosticsParams{
				URI:         request.Params.TextDocument.URI,
				Diagnostics: diagnostics,
			},
		})

	case "textDocument/didChange":
		request := parseMessage[lsp.TextDocumentDidChangeNotification](contents, mh.logger, method)

		for _, change := range request.Params.ContentChanges {
			diagnostics := mh.state.UpdateDocument(request.Params.TextDocument.URI, change.Text)
			mh.sendMessage(lsp.PublishDiagnosticsNotification{
				Notification: lsp.Notification{
					RPC:    "2.0",
					Method: "textDocument/publishDiagnostics",
				},
				Params: lsp.PublishDiagnosticsParams{
					URI:         request.Params.TextDocument.URI,
					Diagnostics: diagnostics,
				},
			})
		}

	case "textDocument/hover":
		request := parseMessage[lsp.HoverRequest](contents, mh.logger, method)

		response := mh.state.Hover(
			request.ID,
			request.Params.TextDocument.URI,
			request.Params.Position,
		)
		mh.sendMessage(response)

	case "textDocument/definition":
		request := parseMessage[lsp.DefinitionRequest](contents, mh.logger, method)

		response := mh.state.Definition(
			request.ID,
			request.Params.TextDocument.URI,
			request.Params.Position,
		)
		mh.sendMessage(response)

	case "textDocument/codeAction":
		request := parseMessage[lsp.CodeActionRequest](contents, mh.logger, method)

		response := mh.state.TextDocumentCodeAction(request.ID, request.Params.TextDocument.URI)
		mh.sendMessage(response)

	case "textDocument/completion":
		request := parseMessage[lsp.CompletionRequest](contents, mh.logger, method)
		mh.logger.Println(request.Params.Position.Line)

		response := mh.state.TextDocumentCompletion(
			request.ID,
			request.Params.Position,
			request.Params.TextDocument.URI,
		)
		mh.sendMessage(response)
	}
}

func (mh *MessageHandler) sendMessage(msg any) {
	reply := rpc.EncodeMessage(msg)
	mh.writer.Write([]byte(reply))
}

func parseMessage[T any](contents []byte, logger *log.Logger, method string) T {
	var request T
	if err := json.Unmarshal(contents, &request); err != nil {
		logger.Printf("Couldn't parse message (%s) in \"%s\"\n", err, method)
	}
	return request
}
