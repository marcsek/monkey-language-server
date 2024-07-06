package parser

import (
	"fmt"
	"github.com/marcsek/monkey-language-server/internal/monkey/token"
	"strings"
)

var traceLevel int = 0

const traceIdentPlaceholder string = "\t"

func identLevel() string {
	return strings.Repeat(traceIdentPlaceholder, traceLevel-1)
}

func tracePrint(fs string) {
	fmt.Printf("%s%s\n", identLevel(), fs)
}

func incIdent() { traceLevel = traceLevel + 1 }
func decIdent() { traceLevel = traceLevel - 1 }

func trace(msg string, position token.Position) string {
	incIdent()
	tracePrint(fmt.Sprintf("BEGIN (%d, %d) %s", position.Line, position.Character, msg))
	return msg
}

func untrace(msg string, position token.Position) {
	tracePrint(fmt.Sprintf("END (%d, %d) %s", position.Line, position.Character, msg))
	decIdent()
}
