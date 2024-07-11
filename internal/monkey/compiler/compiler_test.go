package compiler

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	"github.com/marcsek/monkey-language-server/internal/lsp"
	completion_item_kind "github.com/marcsek/monkey-language-server/internal/lsp/CompletionItemKind"
	"github.com/marcsek/monkey-language-server/internal/monkey/ast"
	"github.com/marcsek/monkey-language-server/internal/monkey/lexer"
	"github.com/marcsek/monkey-language-server/internal/monkey/parser"
	"github.com/marcsek/monkey-language-server/internal/monkey/token"
)

var (
	buff       bytes.Buffer
	MockLogger = log.New(&buff, "", log.LstdFlags)
)

type compilerTestCase struct {
	input string
}

type completionItemWrapper struct {
	lsp.CompletionItem
	shouldAppear bool
}

func TestCompletion(t *testing.T) {
	input := `let glob = 0
  let f = fn(){}

  let ff = fn(x){
    let smth = 0

    let dd = fn(b) {
      let smthh = 1

    }

  }
  `

	testCompletionHelper(
		t,
		input,
		token.Position{Line: 200, Character: 0},
		createCompletionItem("glob", "", "", completion_item_kind.Variable, true),
		createCompletionItem("len", "", "", completion_item_kind.Function, true),
		createCompletionItem("return", "", "", completion_item_kind.Keyword, true),
		createCompletionItem("f", "", "", completion_item_kind.Variable, true),
		createCompletionItem("smth", "", "", completion_item_kind.Variable, false),
		createCompletionItem("x", "", "", completion_item_kind.Variable, false),
	)

	testCompletionHelper(
		t,
		input,
		token.Position{Line: 10, Character: 0},
		createCompletionItem("glob", "", "", completion_item_kind.Variable, true),
		createCompletionItem("len", "", "", completion_item_kind.Function, true),
		createCompletionItem("return", "", "", completion_item_kind.Keyword, true),
		createCompletionItem("ff", "", "", completion_item_kind.Function, true),
		createCompletionItem("smth", "", "", completion_item_kind.Variable, true),
		createCompletionItem("x", "", "", completion_item_kind.Variable, true),
		createCompletionItem("dd", "", "", completion_item_kind.Variable, true),
		createCompletionItem("b", "", "", completion_item_kind.Variable, false),
		createCompletionItem("smthh", "", "", completion_item_kind.Variable, false),
	)

	testCompletionHelper(
		t,
		input,
		token.Position{Line: 8, Character: 0},
		createCompletionItem("glob", "", "", completion_item_kind.Variable, true),
		createCompletionItem("len", "", "", completion_item_kind.Function, true),
		createCompletionItem("return", "", "", completion_item_kind.Keyword, true),
		createCompletionItem("ff", "", "", completion_item_kind.Function, true),
		createCompletionItem("smth", "", "", completion_item_kind.Variable, true),
		createCompletionItem("x", "", "", completion_item_kind.Variable, true),
		createCompletionItem("smth", "", "", completion_item_kind.Variable, true),
		createCompletionItem("x", "", "", completion_item_kind.Variable, true),
		createCompletionItem("dd", "", "", completion_item_kind.Function, true),
		createCompletionItem("b", "", "", completion_item_kind.Variable, true),
		createCompletionItem("smthh", "", "", completion_item_kind.Variable, true),
	)
}

func createCompletionItem(
	label, detail, documentation string,
	kind int,
	shouldAppear bool,
) completionItemWrapper {
	return completionItemWrapper{
		CompletionItem: lsp.CompletionItem{
			Label:         label,
			Detail:        detail,
			Documentation: documentation,
			Kind:          kind,
		},
		shouldAppear: shouldAppear,
	}
}

func testCompletionHelper(
	t *testing.T,
	input string,
	position token.Position,
	expected ...completionItemWrapper,
) {
	t.Helper()

	compiler, err := runCompiler(compilerTestCase{input})
	if err != nil {
		t.Fatal(err)
	}

	result := compiler.Completion(position)

	for _, exp := range expected {
		found := false
		for _, res := range result {
			if exp.Label == res.Label {
				if !exp.shouldAppear {
					t.Fatalf("Label `%s` should not appear", exp.Label)
				}

				found = true
				if exp.Detail != res.Detail {
					t.Fatalf("Wrong `detail` field, want=%s; got=%s", exp.Detail, res.Detail)
				}

				if exp.Documentation != res.Documentation {
					t.Fatalf(
						"Wrong `documentation` field, want=%s; got=%s",
						exp.Documentation,
						res.Documentation,
					)
				}

				if exp.Kind != res.Kind {
					t.Fatalf("Wrong `kind` field, want=%d; got=%d", exp.Kind, res.Kind)
				}

			}
		}
		if !found && exp.shouldAppear {
			t.Fatalf("Couldn't find label with name %s", exp.Label)
		}
	}

}

func runCompiler(test compilerTestCase) (*Compiler, error) {
	program := parse(test.input)

	compiler := New(MockLogger)
	err := compiler.Compile(program)
	if err != nil {
		return nil, fmt.Errorf("compiler error: %s", err)
	}

	return compiler, nil
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}
