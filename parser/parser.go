package parser

import (
	"fmt"
	"jru-test/ast"
	"jru-test/lexer"
	"strconv"
)

func Parse(src string) ast.Node {
	p := parser{
		tokens: lexer.Tokenize(src),
		pos:    0,
	}
	fmt.Printf("Received tokens: %v\n", p.tokens)
	return p.parseExpr(0)
}

type parser struct {
	tokens []lexer.Token
	pos    int
}

func (p *parser) advance() lexer.Token {
	if p.pos < len(p.tokens) {
		token := p.tokens[p.pos]
		p.pos++
		return token
	} else {
		panic("Passed end of tokens without encountering EOF")
	}
}

func (p *parser) peek() lexer.Token {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	panic("Passed end of tokens without encountering EOF")
}

func headPrecedence(tokenType lexer.TokenType) int {
	switch tokenType {
	case lexer.EOF:
		return 0
	case lexer.NUMBER, lexer.STRING, lexer.SYMBOL:
		return 1
	case lexer.PLUS, lexer.DASH:
		return 5
	default:
		panic(fmt.Sprintf("Cannot determine prefix binding power for '%s'", tokenType))
	}
}

func tailPrecedence(tokenType lexer.TokenType) (int, int) {
	switch tokenType {
	case lexer.EOF:
		return 0, 0
	case lexer.PLUS, lexer.DASH:
		return 1, 2
	case lexer.STAR, lexer.SLASH:
		return 3, 4
	default:
		panic(fmt.Sprintf("Cannot determine infix binding power for '%s'", tokenType))
	}
}

func (p *parser) parseExpr(min_bp int) ast.Node {
	token := p.advance()
	left := p.parseHeadExpr(token)

	for {
		nextToken := p.peek()
		if lbp, rbp := tailPrecedence(nextToken.Type); lbp <= min_bp {
			break
		} else {
			left = p.parseTailExpr(left, rbp)
		}
	}
	return left
}

func (p *parser) parseTailExpr(head ast.Node, rbp int) ast.Node {
	token := p.advance()
	switch token.Type {
	case lexer.PLUS, lexer.DASH, lexer.STAR, lexer.SLASH:
		tail := p.parseExpr(rbp)
		return ast.BinaryExprNode{
			Lhs:      head,
			Operator: token,
			Rhs:      tail,
		}
	default:
		panic(fmt.Sprintf("Failed to parse tail expression from token %v\n", token))
	}
}

func (p *parser) parseHeadExpr(token lexer.Token) ast.Node {
	switch token.Type {
	case lexer.NUMBER:
		value, err := strconv.ParseFloat(token.Value, 64)
		if err != nil {
			panic(fmt.Sprintf("Failed to parse number token: %v\n", token))
		}
		return ast.NumberNode{
			Value: value,
		}
	case lexer.STRING:
		return ast.StringNode{
			Value: token.Value,
		}
	case lexer.IDENTIFIER:
		return ast.SymbolNode{
			Value: token.Value,
		}
	case lexer.PLUS, lexer.DASH:
		rbp := headPrecedence(token.Type)
		rhs := p.parseExpr(rbp)
		return ast.UnaryExprNode{
			Operator: token,
			Rhs:      rhs,
		}
	default:
		panic(fmt.Sprintf("Failed to parse head expression from token %v\n", token))
	}
}
