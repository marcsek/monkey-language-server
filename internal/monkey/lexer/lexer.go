package lexer

import "github.com/marcsek/monkey-language-server/internal/monkey/token"

type Lexer struct {
	input        string
	position     int
	readPosition int
	line         int
	linePosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	// Line position increases with readPosition, so it needs to be decremented
	// (because it's pointing at peek)
	startPosition := token.Position{Character: l.linePosition - 1, Line: l.line}

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{
				Type:    token.EQ,
				Literal: literal,
				Range:   createSingleLineRange(startPosition.Character, startPosition.Line, 2),
			}
		} else {
			tok = newToken(token.ASSIGN, l.ch, startPosition)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch, startPosition)
	case '-':
		tok = newToken(token.MINUS, l.ch, startPosition)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{
				Type:    token.NOT_EQ,
				Literal: literal,
				Range:   createSingleLineRange(startPosition.Character, startPosition.Line, 2),
			}
		} else {
			tok = newToken(token.BANG, l.ch, startPosition)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch, startPosition)
	case '*':
		tok = newToken(token.ASTERISK, l.ch, startPosition)
	case '<':
		tok = newToken(token.LT, l.ch, startPosition)
	case '>':
		tok = newToken(token.GT, l.ch, startPosition)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, startPosition)
	case ',':
		tok = newToken(token.COMMA, l.ch, startPosition)
	case '(':
		tok = newToken(token.LPAREN, l.ch, startPosition)
	case ')':
		tok = newToken(token.RPAREN, l.ch, startPosition)
	case '{':
		tok = newToken(token.LBRACE, l.ch, startPosition)
	case '}':
		tok = newToken(token.RBRACE, l.ch, startPosition)
	case '[':
		tok = newToken(token.LBRACKET, l.ch, startPosition)
	case ']':
		tok = newToken(token.RBRACKET, l.ch, startPosition)
	case ':':
		tok = newToken(token.COLON, l.ch, startPosition)
	case '"':
		tok.Type = token.STRING
		literal, length := l.readString()
		tok.Literal = literal
		tok.Range = createSingleLineRange(startPosition.Character, startPosition.Line, length)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			literal, length := l.readIndetifier()
			tok.Literal = literal
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Range = createSingleLineRange(startPosition.Character, startPosition.Line, length)
			return tok
		} else if isDigit(l.ch) {
			literal, length := l.readNumber()
			tok.Literal = literal
			tok.Type = token.INT
			tok.Range = createSingleLineRange(startPosition.Character, startPosition.Line, length)
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch, startPosition)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readNumber() (string, int) {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position], l.position - position
}

func (l *Lexer) readIndetifier() (string, int) {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position], l.position - position
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
	l.linePosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.advanceLine()
		}
		l.readChar()
	}
}

func (l *Lexer) advanceLine() {
	l.line += 1
	l.linePosition = 0
}

func (l *Lexer) readString() (string, int) {
	position := l.position + 1

	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position], l.position - position + 2
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType token.TokenType, ch byte, start token.Position) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
		Range: token.Range{
			Start: start,
			End:   token.Position{Character: start.Character + 1, Line: start.Line},
		},
	}
}

func createSingleLineRange(start, line, length int) token.Range {
	return token.Range{
		Start: token.Position{
			Line:      line,
			Character: start,
		},
		End: token.Position{
			Line:      line,
			Character: start + length,
		},
	}
}
