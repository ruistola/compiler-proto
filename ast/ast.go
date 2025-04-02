package ast

import (
	"fmt"
	"jru-test/lexer"
	"strings"
)

type Type interface {
	_type()
}

type Expr interface {
	expr()
}

type Stmt interface {
	stmt()
}

type SymbolType struct {
	Value string
}

func (n SymbolType) _type() {}

func (n SymbolType) String() string {
	return fmt.Sprintf("(Type %s)", n.Value)
}

type ArrayType struct {
	UnderlyingType Type
}

func (n ArrayType) _type() {}

func (n ArrayType) String() string {
	return fmt.Sprintf("(Type %s[])", n.UnderlyingType)
}

type StringExpr struct {
	Value string
}

func (n StringExpr) expr() {}

func (n StringExpr) String() string {
	return fmt.Sprintf("(String literal \"%s\")", n.Value)
}

type SymbolExpr struct {
	Value string
}

func (n SymbolExpr) expr() {}

func (n SymbolExpr) String() string {
	return fmt.Sprintf("(Symbol \"%s\")", n.Value)
}

type NumberExpr struct {
	Value float64
}

func (n NumberExpr) expr() {}

func (n NumberExpr) String() string {
	return fmt.Sprintf("(Number %v)", n.Value)
}

type UnaryExpr struct {
	Operator lexer.Token
	Rhs      Expr
}

func (n UnaryExpr) expr() {}

func (n UnaryExpr) String() string {
	return fmt.Sprintf("(%s%s)", n.Operator.Value, n.Rhs)
}

type BinaryExpr struct {
	Lhs      Expr
	Operator lexer.Token
	Rhs      Expr
}

func (n BinaryExpr) expr() {}

func (n BinaryExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", n.Operator.Value, n.Lhs, n.Rhs)
}

type BlockStmt struct {
	Body []Stmt
}

func (n BlockStmt) stmt() {}

func (n BlockStmt) String() string {
	var sb strings.Builder
	sb.WriteString("(")
	for _, expr := range n.Body {
		sb.WriteString(fmt.Sprintf("(Stmt %s) ", expr))
	}
	sb.WriteString(")")
	return fmt.Sprintf(sb.String())
}

type ExpressionStmt struct {
	Expr Expr
}

func (n ExpressionStmt) stmt() {}

func (n ExpressionStmt) String() string {
	return fmt.Sprintf("(ExprStmt %s)", n.Expr)
}

type VarDeclStmt struct {
	Var     TypedIdent
	InitVal Expr
}

func (n VarDeclStmt) stmt() {}

func (n VarDeclStmt) String() string {
	var initVal string
	if n.InitVal != nil {
		initVal = fmt.Sprintf("%s", n.InitVal)
	} else {
		initVal = "none"
	}
	return fmt.Sprintf("(Var name %s type %s initial value %s)", n.Var.Name, n.Var.Type, initVal)
}

type TypedIdent struct {
	Name string
	Type Type
}

type FuncDeclStmt struct {
	Name       string
	Parameters []TypedIdent
	ReturnType Type
	Body       BlockStmt
}

func (n FuncDeclStmt) stmt() {}

func (n FuncDeclStmt) String() string {
	var returnType string
	if n.ReturnType != nil {
		returnType = fmt.Sprintf("%s", n.ReturnType)
	} else {
		returnType = "none"
	}
	return fmt.Sprintf("(Func name %s params (%s) returntype %s body { %s })", n.Name, n.Parameters, returnType, n.Body)
}
