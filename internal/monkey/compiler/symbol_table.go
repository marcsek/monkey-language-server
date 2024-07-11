package compiler

import "github.com/marcsek/monkey-language-server/internal/monkey/token"

type SymbolScope string

const (
	GlobalScope   SymbolScope = "GLOBAL"
	LocalScope    SymbolScope = "LOCAL"
	BuiltinScope  SymbolScope = "BUILTIN"
	FreeScope     SymbolScope = "FREE"
	FunctionScope SymbolScope = "FUNCTION"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	Outer *SymbolTable

	store          map[string]Symbol
	numDefinitions int
	depth          int

	tableRange token.Range

	FreeSymbols []Symbol
}

func NewSymbolTable(tableRange token.Range) *SymbolTable {
	s := make(map[string]Symbol)
	free := []Symbol{}
	return &SymbolTable{store: s, FreeSymbols: free, depth: 0, tableRange: tableRange}
}

func NewEnclosedSymbolTable(outer *SymbolTable, tableRange token.Range) *SymbolTable {
	s := NewSymbolTable(tableRange)
	s.depth = outer.depth + 1
	s.Outer = outer
	return s
}

func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions}
	if s.Outer == nil {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}

	s.store[name] = symbol
	s.numDefinitions++
	return symbol
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := s.store[name]
	if !ok && s.Outer != nil {
		obj, ok = s.Outer.Resolve(name)
		if !ok {
			return obj, ok
		}

		if obj.Scope == GlobalScope || obj.Scope == BuiltinScope {
			return obj, ok
		}

		free := s.defineFree(obj)
		return free, true
	}
	return obj, ok
}

func (s *SymbolTable) ResolveAll() []Symbol {
	symbols := map[string]Symbol{}

	s.resolveHelper(&symbols)

	values := make([]Symbol, 0, len(symbols))
	for _, value := range symbols {
		values = append(values, value)
	}

	return values
}

func (s *SymbolTable) resolveHelper(symbols *map[string]Symbol) {
	for key, value := range s.store {
		_, ok := (*symbols)[key]
		if !ok {
			(*symbols)[key] = value
		}
	}
	if s.Outer != nil {
		s.Outer.resolveHelper(symbols)
	}
}

func (s *SymbolTable) DefineBuiltin(index int, name string) Symbol {
	symbol := Symbol{Name: name, Scope: BuiltinScope, Index: index}
	s.store[name] = symbol
	return symbol
}

func (s *SymbolTable) DefineFunctionName(name string) Symbol {
	symbol := Symbol{Name: name, Scope: FunctionScope, Index: 0}
	s.store[name] = symbol
	return symbol
}

func (s *SymbolTable) defineFree(original Symbol) Symbol {
	s.FreeSymbols = append(s.FreeSymbols, original)

	symbol := Symbol{Name: original.Name, Index: len(s.FreeSymbols) - 1}
	symbol.Scope = FreeScope

	s.store[original.Name] = symbol
	return symbol
}
