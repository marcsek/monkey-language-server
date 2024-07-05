package token

import "fmt"

type TokenType string

type Token struct {
	Type    TokenType
	Range   Range
	Literal string
}

type Range struct {
	Start Position
	End   Position
}

type Position struct {
	Line      int
	Character int
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT" // add, foobar, x, y, ...
	INT    = "INT"   // 1343456
	STRING = "STRING"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

	EQ     = "=="
	NOT_EQ = "!="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func (r Range) String() string {
	return fmt.Sprintf(
		"(%d, %d) - (%d, %d)",
		r.Start.Line,
		r.Start.Character,
		r.End.Line,
		r.End.Character,
	)
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
