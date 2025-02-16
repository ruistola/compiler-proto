package main

import (
	"fmt"
)

type Lexer struct {
	tokens []Token
	pos    int
}

func (l *Lexer) advance() Token {
	if l.pos < len(l.tokens) {
		token := l.tokens[l.pos]
		l.pos++
		return token
	} else {
		panic("Passed end of tokens without encountering EOF")
	}
}

func (l *Lexer) peek() (token Token) {
	if l.pos < len(l.tokens) {
		return l.tokens[l.pos]
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
	lexer := Lexer{
		tokens: tokenize(src),
	}
	fmt.Printf("Received tokens: %v\n", lexer.tokens)
	return parseExpr(&lexer, 0)
}

func prefixPrec(tokenType TokenType) int {
	switch tokenType {
	case EOF:
		return 0
	case NUMBER, STRING, SYMBOL:
		return 1
	default:
		panic(fmt.Sprintf("Cannot determine prefix binding power for '%s'", tokenType))
	}
}

func infixPrec(tokenType TokenType) (int, int) {
	switch tokenType {
	case EOF:
		return -1, -1
	case PLUS, DASH:
		return 1, 2
	case STAR, SLASH:
		return 3, 4
	default:
		panic(fmt.Sprintf("Cannot determine infix binding power for '%s'", tokenType))
	}
}

func parseExpr(lexer *Lexer, min_bp int) *AstNode {
	token := lexer.advance()
	left := parsePrefixExpr(token)

	for {
		operator := lexer.peek()
		if lbp, rbp := infixPrec(operator.Type); lbp < min_bp {
			break
		} else {
			token = lexer.advance()
			right := parseExpr(lexer, rbp)
			left = &AstNode{
				token: token,
				left:  left,
				right: right,
			}
		}
	}
	return left
}

func parsePrefixExpr(token Token) *AstNode {
	switch token.Type {
	case NUMBER, STRING, IDENTIFIER:
		node := &AstNode{
			token: token,
			left:  nil,
			right: nil,
		}
		return node
	default:
		panic(fmt.Sprintf("Failed to parse primary expression from token %v\n", token))
	}
}
