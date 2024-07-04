package analysis

import (
	"fmt"
	"strings"

	"github.com/marcsek/monkey-language-server/internal/lsp"
)

type State struct {
	Documents map[string]string
}

func NewState() *State {
	return &State{Documents: map[string]string{}}
}

func (s *State) OpenDocument(uri, text string) {
	s.Documents[uri] = text
}

func (s *State) UpdateDocument(uri, text string) {
	s.Documents[uri] = text
}

func (s *State) Hover(id int, uri string, position lsp.Position) lsp.HoverResponse {

	document := s.Documents[uri]

	response := lsp.HoverResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.HoverResult{
			Contents: fmt.Sprintf("File %s, Characters: %d", uri, len(document)),
		},
	}

	return response
}

func (s *State) Definition(id int, uri string, position lsp.Position) lsp.DefinitionResponse {
	response := lsp.DefinitionResponse{
		Response: lsp.Response{RPC: "2.0", ID: &id},
		Result: lsp.Location{
			URI: uri,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      position.Line + 1,
					Character: 0,
				},
				End: lsp.Position{
					Line:      position.Line + 1,
					Character: 0,
				},
			},
		},
	}

	return response
}

func (s *State) TextDocumentCodeAction(id int, uri string) lsp.CodeActionResponse {
	text := s.Documents[uri]

	actions := []lsp.CodeAction{}
	for row, line := range strings.Split(text, "\n") {
		idx := strings.Index(line, "VS Code")
		if idx >= 0 {
			replaceChange := map[string][]lsp.TextEdit{}
			replaceChange[uri] = []lsp.TextEdit{
				{
					Range: lsp.Range{
						Start: lsp.Position{
							Line:      row,
							Character: idx,
						},
						End: lsp.Position{
							Line:      row,
							Character: idx + len("VS Code"),
						},
					},
					NewText: "Neovim",
				},
			}
			newCodeAction := lsp.CodeAction{
				Title: "Replace \"VS Code\" with superior text editor",
				Edit: &lsp.WorkspaceEdit{
					Changes: replaceChange,
				},
			}
			actions = append(actions, newCodeAction)
		}
	}

	return lsp.CodeActionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: actions,
	}
}
