package ast

import (
	"jru-test/lexer"
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
	TypeName string
}

func (n SymbolType) _type() {}

type ArrayType struct {
	UnderlyingType Type
}

func (n ArrayType) _type() {}

type StringExpr struct {
	Value string
}

func (n StringExpr) expr() {}

type SymbolExpr struct {
	Value string
}

func (n SymbolExpr) expr() {}

type NumberExpr struct {
	Value string
}

func (n NumberExpr) expr() {}

type UnaryExpr struct {
	Operator lexer.Token
	Rhs      Expr
}

func (n UnaryExpr) expr() {}

type BinaryExpr struct {
	Lhs      Expr
	Operator lexer.Token
	Rhs      Expr
}

func (n BinaryExpr) expr() {}

type BlockStmt struct {
	Body []Stmt
}

func (n BlockStmt) stmt() {}

type ExpressionStmt struct {
	Expr Expr
}

func (n ExpressionStmt) stmt() {}

type GroupExpr struct {
	Expr Expr
}

func (n GroupExpr) expr() {}

type VarDeclStmt struct {
	Var     TypedIdent
	InitVal Expr
}

func (n VarDeclStmt) stmt() {}

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

type FuncCallExpr struct {
	Func Expr
	Args []Expr
}

func (n FuncCallExpr) expr() {}

type StructDeclStmt struct {
	Name    string
	Members []TypedIdent
}

func (n StructDeclStmt) stmt() {}

type StructLiteralExpr struct {
	Struct  Expr
	Members []AssignExpr
}

func (n StructLiteralExpr) expr() {}

type IfStmt struct {
	Cond Expr
	Then Stmt
	Else Stmt
}

func (n IfStmt) stmt() {}

type ForStmt struct {
	Init Stmt
	Cond Expr
	Iter Expr
	Body []Stmt
}

func (n ForStmt) stmt() {}

type AssignExpr struct {
	Assigne       Expr
	AssignedValue Expr
}

func (n AssignExpr) expr() {}
