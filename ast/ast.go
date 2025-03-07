package ast

import (
	"fmt"
	"jru-test/lexer"
)

type Node interface {
	node()
}

type StringNode struct {
	Value string
}

func (n StringNode) node() {}

func (n StringNode) String() string {
	return n.Value
}

type SymbolNode struct {
	Value string
}

func (n SymbolNode) node() {}

func (n SymbolNode) String() string {
	return n.Value
}

type NumberNode struct {
	Value float64
}

func (n NumberNode) node() {}

func (n NumberNode) String() string {
	return fmt.Sprintf("%f", n.Value)
}

type UnaryExprNode struct {
	Operator lexer.Token
	Rhs      Node
}

func (n UnaryExprNode) node() {}

func (n UnaryExprNode) String() string {
	return fmt.Sprintf("(%s %s)", n.Operator.Value, n.Rhs)
}

type BinaryExprNode struct {
	Lhs      Node
	Operator lexer.Token
	Rhs      Node
}

func (n BinaryExprNode) node() {}

func (n BinaryExprNode) String() string {
	return fmt.Sprintf("(%s %s %s)", n.Operator.Value, n.Lhs, n.Rhs)
}
