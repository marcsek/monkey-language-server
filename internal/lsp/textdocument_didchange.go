package lsp

type TextDocumentDidChangeNotification struct {
	Notification
	Params DidChangeNotificationParams `json:"params"`
}

type DidChangeNotificationParams struct {
	TextDocument   VersionedTextDocumentIdentifier  `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

type TextDocumentContentChangeEvent struct {
	Text string `json:"text"`
}
