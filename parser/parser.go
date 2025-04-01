package parser

import (
	"fmt"
	"jru-test/ast"
	"jru-test/lexer"
	"slices"
	"strconv"
)

type parser struct {
	tokens []lexer.Token
	pos    int
}

func (p *parser) next() lexer.Token {
	result := lexer.Token{}
	if p.pos < len(p.tokens) {
		result = p.tokens[p.pos]
	}
	return result
}

func (p *parser) consume(expected ...lexer.TokenType) lexer.Token {
	token := p.next()
	if len(expected) > 0 && !slices.Contains(expected, token.Type) {
		panic(fmt.Sprintf("Expected %s, found %s\n", expected, token.Type))
	}
	p.pos++
	return token
}

func headPrecedence(tokenType lexer.TokenType) int {
	switch tokenType {
	case lexer.EOF, lexer.SEMI_COLON:
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
	case lexer.EOF, lexer.SEMI_COLON:
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
	for p.next().Type != lexer.EOF {
		program.Body = append(program.Body, p.parseStmt())
	}
	return program
}

func (p *parser) parseStmt() ast.Stmt {
	switch p.next().Type {
	case lexer.FUNC:
		return p.parseFuncDeclStmt()
	default:
		return p.parseExpressionStmt()
	}
}

func (p *parser) parseArrayType(innerType ast.Type) ast.Type {
	p.consume(lexer.OPEN_BRACKET)
	p.consume(lexer.CLOSE_BRACKET)
	arrayType := ast.ArrayType{
		UnderlyingType: innerType,
	}
	if p.next().Type == lexer.OPEN_BRACKET {
		return p.parseArrayType(arrayType)
	}
	return arrayType
}

func (p *parser) parseType() ast.Type {
	name := p.consume(lexer.IDENTIFIER).Value
	symbolType := ast.SymbolType{
		Value: name,
	}
	if p.next().Type == lexer.OPEN_BRACKET {
		return p.parseArrayType(symbolType)
	}
	return symbolType
}

func (p *parser) parseFuncParm() ast.FuncParm {
	name := p.consume(lexer.IDENTIFIER).Value
	p.consume(lexer.COLON)
	parmType := p.parseType()
	return ast.FuncParm{
		Name: name,
		Type: parmType,
	}
}

func (p *parser) parseFuncDeclStmt() ast.FuncDeclStmt {
	p.consume(lexer.FUNC)
	name := p.consume(lexer.IDENTIFIER).Value

	// Parse function parameter list
	params := make([]ast.FuncParm, 0)
	p.consume(lexer.OPEN_PAREN)

	// While not done with the parameter list...
	for p.next().Type != lexer.CLOSE_PAREN {
		// Parse one parameter (name, colon, type)
		params = append(params, p.parseFuncParm())
		// If followed by a comma, consume the comma and ensure that another parameter follows
		if p.next().Type == lexer.COMMA {
			p.consume(lexer.COMMA)
			if p.next().Type != lexer.IDENTIFIER {
				panic(fmt.Sprintf("Expected identifier after comma in function parameter list, found %s", p.next().Type))
			}
		}
	}
	p.consume(lexer.CLOSE_PAREN)

	// is followed by a return type, if any
	var returnType ast.Type
	if p.next().Type == lexer.COLON {
		p.consume(lexer.COLON)
		returnType = p.parseType()
	}

	// Parse function body
	funcBody := p.parseBlockStmt()

	return ast.FuncDeclStmt{
		Name:       name,
		Parameters: params,
		ReturnType: returnType,
		Body:       funcBody,
	}
}

func (p *parser) parseExpressionStmt() ast.Stmt {
	expr := p.parseExpr(0)
	p.consume(lexer.SEMI_COLON)
	return ast.ExpressionStmt{
		Expr: expr,
	}
}

func (p *parser) parseBlockStmt() ast.BlockStmt {
	p.consume(lexer.OPEN_CURLY)
	body := []ast.Stmt{}
	for nextToken := p.next(); nextToken.Type != lexer.EOF && nextToken.Type != lexer.CLOSE_CURLY; nextToken = p.next() {
		body = append(body, p.parseStmt())
	}
	p.consume(lexer.CLOSE_CURLY)
	return ast.BlockStmt{
		Body: body,
	}
}

func (p *parser) parseExpr(min_bp int) ast.Expr {
	token := p.consume()
	left := p.parseHeadExpr(token)

	for {
		nextToken := p.next()
		if lbp, rbp := tailPrecedence(nextToken.Type); lbp <= min_bp {
			break
		} else {
			left = p.parseTailExpr(left, rbp)
		}
	}
	return left
}

func (p *parser) parseTailExpr(head ast.Expr, rbp int) ast.Expr {
	token := p.consume()
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
