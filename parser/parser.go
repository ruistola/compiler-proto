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

func Parse(tokens []lexer.Token) ast.BlockStmt {
	p := parser{tokens, 0}
	return p.parseProgram()
}

func (p *parser) parseProgram() ast.BlockStmt {
	program := ast.BlockStmt{}
	for p.peek().Type != lexer.EOF {
		program.Body = append(program.Body, p.parseExpressionStmt())
	}
	return program
}

func (p *parser) parseStmt() ast.Stmt {
	switch p.peek().Type {
	case lexer.FUNC:
		return p.parseFuncDeclStmt()
	default:
		return p.parseExpressionStmt()
	}
}

func (p *parser) parseFuncDeclStmt() ast.FuncDeclStmt {
	p.expect(lexer.FUNC)
	name := p.expect(lexer.IDENTIFIER).Value

	// Parse function parameter list
	p.expect(lexer.OPEN_PAREN)
	params := make([]ast.FuncParm, 0)
	p.expect(lexer.CLOSE_PAREN)

	// is followed by a return type, if any
	var returnType ast.Type
	if p.peek().Type == lexer.COLON {
		p.advance()
		// TODO: parse type
	}

	// Parse function body
	funcBody := p.parseBlockStmt()

	return ast.FuncDeclStmt{
		Name:       name,
		Parameters: params,
		ReturnType: returnType,
		Body:       funcBody.Body,
	}
}

func (p *parser) parseExpressionStmt() ast.Stmt {
	expr := p.parseExpr(0)
	p.expect(lexer.SEMI_COLON)
	return ast.ExpressionStmt{
		Expr: expr,
	}
}

func (p *parser) parseBlockStmt() ast.BlockStmt {
	p.expect(lexer.OPEN_CURLY)
	body := []ast.Stmt{}
	for nextToken := p.peek(); nextToken.Type != lexer.EOF && nextToken.Type != lexer.CLOSE_CURLY; nextToken = p.peek() {
		body = append(body, p.parseStmt())
	}
	p.expect(lexer.CLOSE_CURLY)
	return ast.BlockStmt{
		Body: body,
	}
}

func (p *parser) parseExpr(min_bp int) ast.Expr {
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

func (p *parser) parseTailExpr(head ast.Expr, rbp int) ast.Expr {
	token := p.advance()
	switch token.Type {
	case lexer.PLUS, lexer.DASH, lexer.STAR, lexer.SLASH:
		tail := p.parseExpr(rbp)
		return ast.BinaryExpr{
			Lhs:      head,
			Operator: token,
			Rhs:      tail,
		}
	default:
		panic(fmt.Sprintf("Failed to parse tail expression from token %v\n", token))
	}
}

func (p *parser) parseHeadExpr(token lexer.Token) ast.Expr {
	switch token.Type {
	case lexer.NUMBER:
		value, err := strconv.ParseFloat(token.Value, 64)
		if err != nil {
			panic(fmt.Sprintf("Failed to parse number token: %v\n", token))
		}
		return ast.NumberExpr{
			Value: value,
		}
	case lexer.STRING:
		return ast.StringExpr{
			Value: token.Value,
		}
	case lexer.IDENTIFIER:
		return ast.SymbolExpr{
			Value: token.Value,
		}
	case lexer.PLUS, lexer.DASH:
		rbp := headPrecedence(token.Type)
		rhs := p.parseExpr(rbp)
		return ast.UnaryExpr{
			Operator: token,
			Rhs:      rhs,
		}
	default:
		panic(fmt.Sprintf("Failed to parse head expression from token %v\n", token))
	}
}
