package main

type Lexer struct {
	tokens []Token
	pos    int
}

func (l *Lexer) advance() Token {
	token := l.peek()
	l.pos += 1
	return token
}

func (l *Lexer) peek() (token Token) {
	if l.pos+1 < len(l.tokens) {
		return l.tokens[l.pos+1]
	}
	panic("Passed end of tokens without encountering EOF")
}

type BindingPower int

const (
	defaultBp BindingPower = iota
	comma
	assignment
	logical
	relational
	additive
	multiplicative
	unary
	call
	member
	primary
)

type Stmt interface{}

type Expr interface{}

type HeadHandlerFunc func(*Lexer) Expr
type TailHandlerFunc func(*Lexer, Expr, BindingPower) Expr

type HeadHandler struct {
	bp     BindingPower
	handle HeadHandlerFunc
}

type TailHandler struct {
	bp     BindingPower
	handle TailHandlerFunc
}

var (
	headHandlers map[TokenType]HeadHandler = map[TokenType]HeadHandler{
		NUMBER: {primary, parsePrimaryExpr},
	}
	tailHandlers map[TokenType]TailHandler = map[TokenType]TailHandler{}
)

func parse(lexer *Lexer, rbp int) Stmt {
	return nil
	// token := lexer.next()
	// left := handleHead(token)
	// for lbp[lexer.peek()] > rbp {
	// 	token = lexer.next()
	// 	left = handleRest(left)
	// }
	// return left
}

func parsePrimaryExpr(lexer *Lexer) Expr {
	return nil
}
