package main

import (
	"fmt"
	"strconv"
)

type Parser struct {
	tokens []Token
	pos    int
}

func (p *Parser) currentToken() Token {
	return p.tokens[p.pos]
}

func (p *Parser) advance() Token {
	if p.pos < len(p.tokens) {
		token := p.tokens[p.pos]
		p.pos++
		return token
	} else {
		panic("Passed end of tokens without encountering EOF")
	}
}

func (p *Parser) peek() Token {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	panic("Passed end of tokens without encountering EOF")
}

type AstNode interface {
	node()
}

type StringNode struct {
	value string
}

func (n StringNode) node() {}

func (n StringNode) String() string {
	return n.value
}

type SymbolNode struct {
	value string
}

func (n SymbolNode) node() {}

func (n SymbolNode) String() string {
	return n.value
}

type NumberNode struct {
	value float64
}

func (n NumberNode) node() {}

func (n NumberNode) String() string {
	return fmt.Sprintf("%f", n.value)
}

type PrefixExprNode struct {
	operator Token
	rhs      AstNode
}

func (n PrefixExprNode) node() {}

func (n PrefixExprNode) String() string {
	return fmt.Sprintf("(%s %s)", n.operator.Value, n.rhs)
}

type InfixExprNode struct {
	lhs      AstNode
	operator Token
	rhs      AstNode
}

func (n InfixExprNode) node() {}

func (n InfixExprNode) String() string {
	return fmt.Sprintf("(%s %s %s)", n.operator.Value, n.lhs, n.rhs)
}

func parse(src string) AstNode {
	parser := Parser{
		tokens: tokenize(src),
		pos:    0,
	}
	fmt.Printf("Received tokens: %v\n", parser.tokens)
	return parseExpr(&parser, 0)
}

func headPrecedence(tokenType TokenType) int {
	switch tokenType {
	case EOF:
		return 0
	case NUMBER, STRING, SYMBOL:
		return 1
	case PLUS, DASH:
		return 5
	default:
		panic(fmt.Sprintf("Cannot determine prefix binding power for '%s'", tokenType))
	}
}

func tailPrecedence(tokenType TokenType) (int, int) {
	switch tokenType {
	case EOF:
		return 0, 0
	case PLUS, DASH:
		return 1, 2
	case STAR, SLASH:
		return 3, 4
	default:
		panic(fmt.Sprintf("Cannot determine infix binding power for '%s'", tokenType))
	}
}

func parseExpr(p *Parser, min_bp int) AstNode {
	token := p.advance()
	left := parseHeadExpr(p, token)

	for {
		nextToken := p.peek()
		if lbp, rbp := tailPrecedence(nextToken.Type); lbp <= min_bp {
			break
		} else {
			left = parseTailExpr(p, left, rbp)
		}
	}
	return left
}

func parseTailExpr(p *Parser, head AstNode, rbp int) AstNode {
	token := p.advance()
	switch token.Type {
	case PLUS, DASH, STAR, SLASH:
		tail := parseExpr(p, rbp)
		return InfixExprNode{
			lhs:      head,
			operator: token,
			rhs:      tail,
		}
	default:
		panic(fmt.Sprintf("Failed to parse tail expression from token %v\n", token))
	}
}

func parseHeadExpr(p *Parser, token Token) AstNode {
	switch token.Type {
	case NUMBER:
		value, err := strconv.ParseFloat(token.Value, 64)
		if err != nil {
			panic(fmt.Sprintf("Failed to parse number token: %v\n", token))
		}
		return NumberNode{
			value,
		}
	case STRING:
		return StringNode{
			value: token.Value,
		}
	case IDENTIFIER:
		return SymbolNode{
			value: token.Value,
		}
	case PLUS, DASH:
		rbp := headPrecedence(token.Type)
		rhs := parseExpr(p, rbp)
		return PrefixExprNode{
			operator: token,
			rhs:      rhs,
		}
	default:
		panic(fmt.Sprintf("Failed to parse head expression from token %v\n", token))
	}
}
