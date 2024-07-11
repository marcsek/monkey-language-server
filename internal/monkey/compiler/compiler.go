package compiler

import (
	"fmt"
	"log"
	"sort"

	"github.com/marcsek/monkey-language-server/internal/lsp"
	"github.com/marcsek/monkey-language-server/internal/lsp/CompletionItemKind"
	"github.com/marcsek/monkey-language-server/internal/monkey/ast"
	"github.com/marcsek/monkey-language-server/internal/monkey/object"
	"github.com/marcsek/monkey-language-server/internal/monkey/token"
)

type relativeSymbolTable struct {
	symbolTable *SymbolTable
	outer       *relativeSymbolTable
	tableRange  token.Range
}

type Compiler struct {
	relativeSymbolTable *relativeSymbolTable
	symbolTableMap      map[string]relativeSymbolTable
	logger              *log.Logger

	scopeIndex int
}

func New(logger *log.Logger) *Compiler {
	symbolTable := NewSymbolTable()

	relativeST := &relativeSymbolTable{
		symbolTable: symbolTable,
		tableRange:  token.Range{},
	}
	symbolTableMap := map[string]relativeSymbolTable{}

	return &Compiler{
		relativeSymbolTable: relativeST,
		symbolTableMap:      symbolTableMap,
		scopeIndex:          0,

		logger: logger,
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}

	case *ast.LetStatement:
		c.relativeSymbolTable.symbolTable.Define(node.Name.Value)
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}

	case *ast.IntegerLiteral, *ast.StringLiteral, *ast.Boolean:
		return nil

	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}
		return nil

	case *ast.InfixExpression:
		if node.Operator == "<" {
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}

			err = c.Compile(node.Left)
			if err != nil {
				return err
			}

			return nil
		}

		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Right)
		if err != nil {
			return err
		}

	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}

		if node.Alternative != nil {
			err = c.Compile(node.Alternative)
			if err != nil {
				return err
			}
		}

	case *ast.BlockStatement:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.Identifier:
		_, ok := c.relativeSymbolTable.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable %s", node.Value)
		}

	case *ast.ArrayLiteral:
		for _, s := range node.Elements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.HashLiteral:
		keys := []ast.Expression{}
		for k := range node.Pairs {
			keys = append(keys, k)
		}

		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})

		for _, key := range keys {
			err := c.Compile(key)
			if err != nil {
				return err
			}

			err = c.Compile(node.Pairs[key])
			if err != nil {
				return err
			}
		}

	case *ast.IndexExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Index)
		if err != nil {
			return err
		}

	case *ast.FunctionLiteral:
		c.enterScope(node.Body.Range())

		defer c.leaveScope()

		if node.Name != "" {
			c.relativeSymbolTable.symbolTable.DefineFunctionName(node.Name)
		}

		for _, p := range node.Parameters {
			c.relativeSymbolTable.symbolTable.Define(p.Value)
		}

		err := c.Compile(node.Body)
		if err != nil {
			return err
		}

	case *ast.ReturnStatement:
		err := c.Compile(node.ReturnValue)
		if err != nil {
			return err
		}

	case *ast.CallExpression:
		err := c.Compile(node.Function)
		if err != nil {
			return err
		}

		for _, a := range node.Arguments {
			err := c.Compile(a)
			if err != nil {
				return err
			}
		}

	default:
		return fmt.Errorf("unknown operator %s", node.TokenLiteral())
	}

	return nil
}

func (c *Compiler) enterScope(tableRange token.Range) {
	newRelativeST := &relativeSymbolTable{
		symbolTable: NewEnclosedSymbolTable(c.relativeSymbolTable.symbolTable),
		tableRange:  tableRange,
		outer:       c.relativeSymbolTable,
	}
	c.relativeSymbolTable = newRelativeST
}

func (c *Compiler) leaveScope() {
	c.symbolTableMap[c.relativeSymbolTable.tableRange.String()] = *c.relativeSymbolTable
	c.relativeSymbolTable = c.relativeSymbolTable.outer
}

func (c *Compiler) Completion(position token.Position) []lsp.CompletionItem {
	items := []lsp.CompletionItem{}

	for _, name := range object.Constants {
		items = append(items, lsp.CompletionItem{Label: name, Kind: completion_item_kind.Constant})
	}

	for _, name := range object.Builtins {
		items = append(items, lsp.CompletionItem{Label: name, Kind: completion_item_kind.Function})
	}

	for _, name := range object.Keywords {
		items = append(items, lsp.CompletionItem{Label: name, Kind: completion_item_kind.Keyword})
	}

	st := c.findMostSpecificScope(position)
	for _, symbol := range st.ResolveAll() {
		kind := completion_item_kind.Variable
		if symbol.Scope == FunctionScope {
			kind = completion_item_kind.Function
		}

		items = append(
			items,
			lsp.CompletionItem{Label: symbol.Name, Kind: kind},
		)
	}

	return items
}

func (c *Compiler) findMostSpecificScope(position token.Position) *SymbolTable {
	mostSpecific := c.relativeSymbolTable.symbolTable
	depth := c.relativeSymbolTable.symbolTable.depth

	for _, symbol := range c.symbolTableMap {
		if symbol.tableRange.Start.Line < position.Line &&
			symbol.tableRange.End.Line > position.Line && depth < symbol.symbolTable.depth {
			mostSpecific = symbol.symbolTable
			depth = mostSpecific.depth
		}
	}

	return mostSpecific
}
