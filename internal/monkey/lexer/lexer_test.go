package lexer

import (
	"testing"

	"github.com/marcsek/monkey-language-server/internal/monkey/token"
)

func TestNextToken(t *testing.T) {
	input := `+-*/=!
*/==!=

*   /

"asd"

"daco"
""
==

 !=

ahoj
let ten = 10;

{"foo": "bar"}
`
	type Test struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedRange   token.Range
	}

	tests := []Test{
		{token.PLUS, "+", createSingleLineRange(0, 0, 1)},
		{token.MINUS, "-", createSingleLineRange(1, 0, 1)},
		{token.ASTERISK, "*", createSingleLineRange(2, 0, 1)},
		{token.SLASH, "/", createSingleLineRange(3, 0, 1)},
		{token.ASSIGN, "=", createSingleLineRange(4, 0, 1)},
		{token.BANG, "!", createSingleLineRange(5, 0, 1)},
		{token.ASTERISK, "*", createSingleLineRange(0, 1, 1)},
		{token.SLASH, "/", createSingleLineRange(1, 1, 1)},
		{token.EQ, "==", createSingleLineRange(2, 1, 2)},
		{token.NOT_EQ, "!=", createSingleLineRange(4, 1, 2)},
		{token.ASTERISK, "*", createSingleLineRange(0, 3, 1)},
		{token.SLASH, "/", createSingleLineRange(4, 3, 1)},
		{token.STRING, "asd", createSingleLineRange(0, 5, 5)},
		{token.STRING, "daco", createSingleLineRange(0, 7, 6)},
		{token.STRING, "", createSingleLineRange(0, 8, 2)},
		{token.EQ, "==", createSingleLineRange(0, 9, 2)},
		{token.NOT_EQ, "!=", createSingleLineRange(1, 11, 2)},
		{token.IDENT, "ahoj", createSingleLineRange(0, 13, 4)},
		{token.LET, "let", createSingleLineRange(0, 14, 3)},
		{token.IDENT, "ten", createSingleLineRange(4, 14, 3)},
		{token.ASSIGN, "=", createSingleLineRange(8, 14, 1)},
		{token.INT, "10", createSingleLineRange(10, 14, 2)},
		{token.SEMICOLON, ";", createSingleLineRange(12, 14, 1)},
		{token.LBRACE, "{", createSingleLineRange(0, 16, 1)},
		{token.STRING, "foo", createSingleLineRange(1, 16, 5)},
		{token.COLON, ":", createSingleLineRange(6, 16, 1)},
		{token.STRING, "bar", createSingleLineRange(8, 16, 5)},
		{token.RBRACE, "}", createSingleLineRange(13, 16, 1)},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}

		if !compareRange(tok.Range, tt.expectedRange) {
			t.Fatalf("tests[%d] - range wrong. expected=%s, got=%s",
				i, tt.expectedRange, tok.Range)
		}
	}
}

func compareRange(r1, r2 token.Range) bool {
	return r1.String() == r2.String()
}

//func TestNextToken(t *testing.T) {
//	input := `let five = 5;
//let ten = 10;
//
//let add = fn(x, y) {
//  x + y;
//};
//
//let result = add(five, ten);
//!-/*5;
//5 < 10 > 5;
//
//if (5 < 10) {
//	return true;
//} else {
//	return false;
//}
//
//10 == 10;
//10 != 9;
//
//"foobar"
//"foo bar"
//[1, 2];
//{"foo": "bar"}
//`
//	type Test struct {
//		expectedType    token.TokenType
//		expectedLiteral string
//	}
//
//	tests := []Test{
//		{token.LET, "let"},
//		{token.IDENT, "five"},
//		{token.ASSIGN, "="},
//		{token.INT, "5"},
//		{token.SEMICOLON, ";"},
//		{token.LET, "let"},
//		{token.IDENT, "ten"},
//		{token.ASSIGN, "="},
//		{token.INT, "10"},
//		{token.SEMICOLON, ";"},
//		{token.LET, "let"},
//		{token.IDENT, "add"},
//		{token.ASSIGN, "="},
//		{token.FUNCTION, "fn"},
//		{token.LPAREN, "("},
//		{token.IDENT, "x"},
//		{token.COMMA, ","},
//		{token.IDENT, "y"},
//		{token.RPAREN, ")"},
//		{token.LBRACE, "{"},
//		{token.IDENT, "x"},
//		{token.PLUS, "+"},
//		{token.IDENT, "y"},
//		{token.SEMICOLON, ";"},
//		{token.RBRACE, "}"},
//		{token.SEMICOLON, ";"},
//		{token.LET, "let"},
//		{token.IDENT, "result"},
//		{token.ASSIGN, "="},
//		{token.IDENT, "add"},
//		{token.LPAREN, "("},
//		{token.IDENT, "five"},
//		{token.COMMA, ","},
//		{token.IDENT, "ten"},
//		{token.RPAREN, ")"},
//		{token.SEMICOLON, ";"},
//		{token.BANG, "!"},
//		{token.MINUS, "-"},
//		{token.SLASH, "/"},
//		{token.ASTERISK, "*"},
//		{token.INT, "5"},
//		{token.SEMICOLON, ";"},
//		{token.INT, "5"},
//		{token.LT, "<"},
//		{token.INT, "10"},
//		{token.GT, ">"},
//		{token.INT, "5"},
//		{token.SEMICOLON, ";"},
//		{token.IF, "if"},
//		{token.LPAREN, "("},
//		{token.INT, "5"},
//		{token.LT, "<"},
//		{token.INT, "10"},
//		{token.RPAREN, ")"},
//		{token.LBRACE, "{"},
//		{token.RETURN, "return"},
//		{token.TRUE, "true"},
//		{token.SEMICOLON, ";"},
//		{token.RBRACE, "}"},
//		{token.ELSE, "else"},
//		{token.LBRACE, "{"},
//		{token.RETURN, "return"},
//		{token.FALSE, "false"},
//		{token.SEMICOLON, ";"},
//		{token.RBRACE, "}"},
//		{token.INT, "10"},
//		{token.EQ, "=="},
//		{token.INT, "10"},
//		{token.SEMICOLON, ";"},
//		{token.INT, "10"},
//		{token.NOT_EQ, "!="},
//		{token.INT, "9"},
//		{token.SEMICOLON, ";"},
//		{token.STRING, "foobar"},
//		{token.STRING, "foo bar"},
//		{token.LBRACKET, "["},
//		{token.INT, "1"},
//		{token.COMMA, ","},
//		{token.INT, "2"},
//		{token.RBRACKET, "]"},
//		{token.SEMICOLON, ";"},
//		{token.LBRACE, "{"},
//		{token.STRING, "foo"},
//		{token.COLON, ":"},
//		{token.STRING, "bar"},
//		{token.RBRACE, "}"},
//		{token.EOF, ""},
//	}
//
//	l := New(input)
//
//	for i, tt := range tests {
//		tok := l.NextToken()
//
//		if tok.Type != tt.expectedType {
//			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
//				i, tt.expectedType, tok.Type)
//		}
//
//		if tok.Literal != tt.expectedLiteral {
//			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
//				i, tt.expectedLiteral, tok.Literal)
//		}
//	}
//}
