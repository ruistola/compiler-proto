package parser

import (
	"fmt"
	"jru-test/ast"
	"jru-test/lexer"
	"strconv"
)

type parser struct {
	tokens []lexer.Token
	pos    int
}

func (p *parser) peek() lexer.Token {
	result := lexer.Token{}
	if p.pos < len(p.tokens) {
		result = p.tokens[p.pos]
	}
	return result
}

func (p *parser) advance() lexer.Token {
	nextToken := p.peek()
	p.pos++
	return nextToken
}

func (p *parser) expect(expected lexer.TokenType) lexer.Token {
	nextToken := p.advance()
	if nextToken.Type != expected {
		panic(fmt.Sprintf("Expected %s, found %s\n", expected, nextToken.Type))
	}
	return nextToken
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
		panic(fmt.Sprintf("Cannot determine binding power for '%s' as a head token", tokenType))
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
		panic(fmt.Sprintf("Cannot determine binding power for '%s' as a tail token", tokenType))
	}
}

func Parse(tokens []lexer.Token) ast.Node {
	p := parser{tokens, 0}
	return p.parseProgram()
}

func (p *parser) parseProgram() ast.Node {
	program := ast.BlockExprNode{}
	for p.peek().Type != lexer.EOF {
		program.Body = append(program.Body, p.parseExpr(0))
	}
	return program
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
	// case lexer.FUNC:
	//     rbp := headPrecedence(token.Type)
	//           name := p.expect(lexer.IDENTIFIER)
	//     _ = p.expect(lexer.OPEN_PAREN)
	// TODO: parse parameter list
	//     _ = p.expect(lexer.CLOSE_PAREN)
	// return ast.FuncDeclNode{}
	default:
		panic(fmt.Sprintf("Failed to parse head expression from token %v\n", token))
	}
}
