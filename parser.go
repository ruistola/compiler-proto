package main

import (
	"fmt"
)

type Parser struct {
	tokens []Token
	pos    int
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

func (p *Parser) peek() (token Token) {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	panic("Passed end of tokens without encountering EOF")
}

type AstNode struct {
	token Token
	left  *AstNode
	right *AstNode
}

func (node *AstNode) String() string {
	switch {
	case node.left == nil && node.right == nil:
		return fmt.Sprintf("%s", node.token.Value)
	case node.left == nil:
		return fmt.Sprintf("(%s %s)", node.token.Value, node.right)
	case node.right == nil:
		return fmt.Sprintf("(%s %s)", node.token.Value, node.left)
	default:
		return fmt.Sprintf("(%s (%s %s))", node.token.Value, node.left, node.right)
	}
}

func parse(src string) *AstNode {
	parser := Parser{
		tokens: tokenize(src),
	}
	fmt.Printf("Received tokens: %v\n", parser.tokens)
	return parseExpr(&parser, 0)
}

func prefixPrec(tokenType TokenType) int {
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

func infixPrec(tokenType TokenType) (int, int) {
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

func parseExpr(p *Parser, min_bp int) *AstNode {
	token := p.advance()
	left := parsePrefixExpr(p, token)

	for {
		nextToken := p.peek()
		if lbp, rbp := infixPrec(nextToken.Type); lbp <= min_bp {
			break
		} else {
			token = p.advance()
			right := parseExpr(p, rbp)
			left = &AstNode{
				token: token,
				left:  left,
				right: right,
			}
		}
	}
	return left
}

func parsePrefixExpr(p *Parser, token Token) *AstNode {
	switch token.Type {
	case NUMBER, STRING, IDENTIFIER:
		node := &AstNode{
			token: token,
			left:  nil,
			right: nil,
		}
		return node
	case PLUS, DASH:
		bp := prefixPrec(token.Type)
		right := parseExpr(p, bp)
		node := &AstNode{
			token: token,
			left:  nil,
			right: right,
		}
		return node
	default:
		panic(fmt.Sprintf("Failed to parse primary expression from token %v\n", token))
	}
}
