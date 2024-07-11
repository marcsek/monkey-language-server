package analysis

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/marcsek/monkey-language-server/internal/lsp"
	"github.com/marcsek/monkey-language-server/internal/monkey/compiler"
	"github.com/marcsek/monkey-language-server/internal/monkey/lexer"
	"github.com/marcsek/monkey-language-server/internal/monkey/parser"
	"github.com/marcsek/monkey-language-server/internal/monkey/token"
)

type State struct {
	Documents map[string]string
	Compiler  *compiler.Compiler
	logger    *log.Logger
}

func NewState(logger *log.Logger) *State {
	return &State{Documents: map[string]string{}, logger: logger}
}

func (s *State) createCompiler(text string) *compiler.Compiler {
	//constants := []object.Object{}
	//globals := make([]object.Object, vm.GlobalsSize)

	start := time.Now()

	//for i, name := range object.Builtins {
	//	symbolTable.DefineBuiltin(i, name)
	//}

	l := lexer.New(text)
	p := parser.New(l)

	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		s.logger.Printf("Parser Errors: %v", p.Errors())
	}

	comp := compiler.New(s.logger)

	err := comp.Compile(program)
	if err != nil {
		s.logger.Printf("Compilation error: %s", err)
	}
	total := time.Since(start)

	s.logger.Printf("Compile time: %s", total)
	_ = program

	return comp
}

func (s *State) OpenDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text
	s.Compiler = s.createCompiler(text)

	return getDiagnosticsForFile(text)
}

func (s *State) UpdateDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text
	s.Compiler = s.createCompiler(text)

	return getDiagnosticsForFile(text)
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

func getDiagnosticsForFile(text string) []lsp.Diagnostic {
	diagnostics := []lsp.Diagnostic{}
	for row, line := range strings.Split(text, "\n") {
		idx := strings.Index(line, "VS Code")
		if idx >= 0 {
			diagnostics = append(diagnostics, lsp.Diagnostic{
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
				Severity: 1,
				Source:   "The god him self",
				Message:  "Fix this asap",
			})
		}
	}
	return diagnostics
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

func (s *State) TextDocumentCompletion(
	id int,
	position lsp.Position,
	uri string,
) lsp.CompletionResponse {
	items := s.Compiler.Completion(token.Position(position))

	return lsp.CompletionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: items,
	}
}
