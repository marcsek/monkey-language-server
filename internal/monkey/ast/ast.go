package ast

import (
	"fmt"
	"github.com/marcsek/monkey-language-server/internal/monkey/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
	Range() token.Range
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) Range() token.Range {
	if len(p.Statements) == 0 {
		return token.Range{}
	} else if len(p.Statements) == 1 {
		return p.Statements[0].Range()
	} else {
		return token.Range{
			Start: p.Statements[0].Range().Start,
			End:   p.Statements[len(p.Statements)-1].Range().End,
		}
	}
}

func (p *Program) String() string {
	var sb strings.Builder

	for _, s := range p.Statements {
		sb.WriteString(s.String())
	}

	return sb.String()
}

type LetStatement struct {
	Token      token.Token
	Name       *Identifier
	Value      Expression
	RangeValue token.Range
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) Range() token.Range   { return ls.RangeValue }

func (ls *LetStatement) String() string {
	var sb strings.Builder

	sb.WriteString(ls.TokenLiteral() + " ")
	sb.WriteString(ls.Name.String())
	sb.WriteString(" = ")

	if ls.Value != nil {
		sb.WriteString(ls.Value.String())
	}

	sb.WriteString(";")

	return sb.String()
}

type Identifier struct {
	Token      token.Token
	Value      string
	RangeValue token.Range
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) Range() token.Range   { return i.RangeValue }
func (i *Identifier) String() string       { return i.Value }

type ReturnStatement struct {
	Token       token.Token
	RangeValue  token.Range
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) Range() token.Range   { return rs.RangeValue }
func (rs *ReturnStatement) String() string {
	var sb strings.Builder

	sb.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		sb.WriteString(rs.ReturnValue.String())
	}

	sb.WriteString(";")

	return sb.String()
}

type ExpressionStatement struct {
	Token      token.Token
	RangeValue token.Range
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) Range() token.Range   { return es.RangeValue }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type IntegerLiteral struct {
	Token      token.Token
	RangeValue token.Range
	Value      int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) Range() token.Range   { return il.RangeValue }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type StringLiteral struct {
	Token      token.Token
	RangeValue token.Range
	Value      string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) Range() token.Range   { return sl.RangeValue }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type PrefixExpression struct {
	Token      token.Token
	RangeValue token.Range
	Operator   string
	Right      Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) Range() token.Range   { return pe.RangeValue }
func (pe *PrefixExpression) String() string {
	var sb strings.Builder

	sb.WriteString("(")
	sb.WriteString(pe.Operator)
	sb.WriteString(pe.Right.String())
	sb.WriteString(")")

	return sb.String()
}

type InfixExpression struct {
	Token      token.Token
	Left       Expression
	RangeValue token.Range
	Operator   string
	Right      Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) Range() token.Range   { return ie.RangeValue }
func (ie *InfixExpression) String() string {
	var sb strings.Builder

	sb.WriteString("(")
	sb.WriteString(ie.Left.String())
	sb.WriteString(" " + ie.Operator + " ")
	sb.WriteString(ie.Right.String())
	sb.WriteString(")")

	return sb.String()
}

type Boolean struct {
	Token      token.Token
	RangeValue token.Range
	Value      bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) Range() token.Range   { return b.RangeValue }
func (b *Boolean) String() string       { return b.Token.Literal }

type IfExpression struct {
	Token       token.Token
	RangeValue  token.Range
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) Range() token.Range   { return ie.RangeValue }
func (ie *IfExpression) String() string {
	var sb strings.Builder

	sb.WriteString("if")
	sb.WriteString(ie.Condition.String())
	sb.WriteString(" ")
	sb.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		sb.WriteString("else ")
		sb.WriteString(ie.Alternative.String())
	}

	return sb.String()
}

type BlockStatement struct {
	Token      token.Token
	RangeValue token.Range
	Statements []Statement
}

func (bs *BlockStatement) expressionNode()      {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) Range() token.Range   { return bs.RangeValue }
func (bs *BlockStatement) String() string {
	var sb strings.Builder

	for _, s := range bs.Statements {
		sb.WriteString(s.String())
	}

	return sb.String()
}

type FunctionLiteral struct {
	Token      token.Token
	RangeValue token.Range
	Parameters []*Identifier
	Body       *BlockStatement
	Name       string
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) Range() token.Range   { return fl.RangeValue }
func (fl *FunctionLiteral) String() string {
	var sb strings.Builder

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	sb.WriteString(fl.TokenLiteral())
	if fl.Name != "" {
		sb.WriteString(fmt.Sprintf("<%s>", fl.Name))
	}
	sb.WriteString("(")
	sb.WriteString(strings.Join(params, ", "))
	sb.WriteString(") ")
	sb.WriteString(fl.Body.String())

	return sb.String()
}

type CallExpression struct {
	Token      token.Token
	RangeValue token.Range
	Function   Expression
	Arguments  []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) Range() token.Range   { return ce.RangeValue }
func (ce *CallExpression) String() string {
	var sb strings.Builder

	args := []string{}
	for _, p := range ce.Arguments {
		args = append(args, p.String())
	}

	sb.WriteString(ce.Function.String())
	sb.WriteString("(")
	sb.WriteString(strings.Join(args, ", "))
	sb.WriteString(")")

	return sb.String()
}

type ArrayLiteral struct {
	Token      token.Token
	RangeValue token.Range
	Elements   []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) Range() token.Range   { return al.RangeValue }
func (al *ArrayLiteral) String() string {
	var sb strings.Builder

	elements := []string{}
	for _, p := range al.Elements {
		elements = append(elements, p.String())
	}

	sb.WriteString("[")
	sb.WriteString(strings.Join(elements, ", "))
	sb.WriteString("]")

	return sb.String()
}

type IndexExpression struct {
	Token      token.Token
	RangeValue token.Range
	Left       Expression
	Index      Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) Range() token.Range   { return ie.RangeValue }
func (ie *IndexExpression) String() string {
	var sb strings.Builder

	sb.WriteString("(")
	sb.WriteString(ie.Left.String())
	sb.WriteString("[")
	sb.WriteString(ie.Index.String())
	sb.WriteString("])")

	return sb.String()
}

type HashLiteral struct {
	Token      token.Token
	RangeValue token.Range
	Pairs      map[Expression]Expression
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) Range() token.Range   { return hl.RangeValue }
func (hl *HashLiteral) String() string {
	var sb strings.Builder

	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	sb.WriteString("{")
	sb.WriteString(strings.Join(pairs, ", "))
	sb.WriteString("}")

	return sb.String()
}
