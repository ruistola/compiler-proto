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

func (p *parser) peek() lexer.Token {
	result := lexer.Token{}
	if p.pos < len(p.tokens) {
		result = p.tokens[p.pos]
	}
	return result
}

func (p *parser) consume(expected ...lexer.TokenType) lexer.Token {
	token := p.peek()
	if len(expected) > 0 && !slices.Contains(expected, token.Type) {
		panic(fmt.Sprintf("Expected %s, found %s\n", expected, token.Type))
	}
	p.pos++
	return token
}

func headPrecedence(tokenType lexer.TokenType) int {
	switch tokenType {
	case lexer.EOF, lexer.SEMI_COLON, lexer.OPEN_PAREN:
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
	case lexer.EOF, lexer.SEMI_COLON, lexer.CLOSE_PAREN:
		return 0, 0
	case lexer.ASSIGNMENT, lexer.PLUS_EQUALS, lexer.MINUS_EQUALS:
		return 1, 2
	case lexer.EQUALS, lexer.NOT_EQUALS:
		return 3, 4
	case lexer.LESS, lexer.LESS_EQUALS, lexer.GREATER, lexer.GREATER_EQUALS:
		return 5, 6
	case lexer.PLUS, lexer.DASH:
		return 7, 8
	case lexer.STAR, lexer.SLASH, lexer.PERCENT:
		return 9, 10
	case lexer.OPEN_PAREN:
		return 11, 0
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
		program.Body = append(program.Body, p.parseStmt())
	}
	return program
}

func (p *parser) parseStmt() ast.Stmt {
	switch p.peek().Type {
	case lexer.OPEN_CURLY:
		return p.parseBlockStmt()
	case lexer.FUNC:
		return p.parseFuncDeclStmt()
	case lexer.IF:
		return p.parseIfStmt()
	case lexer.FOR:
		return p.parseForStmt()
	case lexer.LET:
		return p.parseVarDeclStmt()
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
	if p.peek().Type == lexer.OPEN_BRACKET {
		return p.parseArrayType(arrayType)
	}
	return arrayType
}

func (p *parser) parseType() ast.Type {
	name := p.consume(lexer.IDENTIFIER).Value
	symbolType := ast.SymbolType{
		TypeName: name,
	}
	if p.peek().Type == lexer.OPEN_BRACKET {
		return p.parseArrayType(symbolType)
	}
	return symbolType
}

func (p *parser) parseVarDeclStmt() ast.VarDeclStmt {
	p.consume(lexer.LET)
	varName := p.consume(lexer.IDENTIFIER).Value
	p.consume(lexer.COLON)
	varType := p.parseType()
	var initVal ast.Expr
	if p.peek().Type != lexer.SEMI_COLON {
		p.consume(lexer.ASSIGNMENT)
		initVal = p.parseExpr(0)
	}
	p.consume(lexer.SEMI_COLON)
	return ast.VarDeclStmt{
		Var: ast.TypedIdent{
			Name: varName,
			Type: varType,
		},
		InitVal: initVal,
	}
}

func (p *parser) parseFuncParm() ast.TypedIdent {
	name := p.consume(lexer.IDENTIFIER).Value
	p.consume(lexer.COLON)
	parmType := p.parseType()
	return ast.TypedIdent{
		Name: name,
		Type: parmType,
	}
}

func (p *parser) parseFuncDeclStmt() ast.FuncDeclStmt {
	p.consume(lexer.FUNC)
	name := p.consume(lexer.IDENTIFIER).Value

	// Parse function parameter list
	params := make([]ast.TypedIdent, 0)
	p.consume(lexer.OPEN_PAREN)

	// While not done with the parameter list...
	for p.peek().Type != lexer.CLOSE_PAREN {
		// Parse one parameter (name, colon, type)
		params = append(params, p.parseFuncParm())
		// If followed by a comma, consume the comma and ensure that another parameter follows
		if p.peek().Type == lexer.COMMA {
			p.consume(lexer.COMMA)
			if p.peek().Type != lexer.IDENTIFIER {
				panic(fmt.Sprintf("Expected identifier after comma in function parameter list, found %s", p.peek().Type))
			}
		}
	}
	p.consume(lexer.CLOSE_PAREN)

	// is followed by a return type, if any
	var returnType ast.Type
	if p.peek().Type == lexer.COLON {
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

func (p *parser) parseIfStmt() ast.Stmt {
	p.consume(lexer.IF)

	// Parse the condition enclosed by parens
	p.consume(lexer.OPEN_PAREN)
	cond := p.parseExpr(0)
	p.consume(lexer.CLOSE_PAREN)

	// Parse the consequent
	thenStmt := p.parseStmt()

	// Parse the alternate, if any
	var elseStmt ast.Stmt
	if p.peek().Type == lexer.ELSE {
		p.consume(lexer.ELSE)
		elseStmt = p.parseStmt()
	}

	return ast.IfStmt{
		Cond: cond,
		Then: thenStmt,
		Else: elseStmt,
	}
}

func (p *parser) parseForStmt() ast.Stmt {
	p.consume(lexer.FOR)
	p.consume(lexer.OPEN_PAREN)
	initStmt := p.parseStmt()
	condExpr := p.parseExpressionStmt().(ast.ExpressionStmt).Expr
	iterStmt := p.parseExpr(0)
	p.consume(lexer.CLOSE_PAREN)
	body := p.parseBlockStmt()
	return ast.ForStmt{
		Init: initStmt,
		Cond: condExpr,
		Iter: iterStmt,
		Body: body.Body,
	}
}

func (p *parser) parseFuncCallExpr(left ast.Expr) ast.FuncCallExpr {
	args := []ast.Expr{}
	for p.peek().Type != lexer.CLOSE_PAREN {
		// Parse the argument expression
		args = append(args, p.parseExpr(0))

		// If followed by a comma, consume the comma
		if p.peek().Type == lexer.COMMA {
			p.consume(lexer.COMMA)
		}
	}
	p.consume(lexer.CLOSE_PAREN)
	return ast.FuncCallExpr{
		Func: left,
		Args: args,
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
	for nextToken := p.peek(); nextToken.Type != lexer.EOF && nextToken.Type != lexer.CLOSE_CURLY; nextToken = p.peek() {
		body = append(body, p.parseStmt())
	}
	p.consume(lexer.CLOSE_CURLY)
	return ast.BlockStmt{
		Body: body,
	}
}

func (p *parser) parseExpr(min_bp int) ast.Expr {
	token := p.consume()
	parsedExpr := p.parseHeadExpr(token)

	for {
		token = p.peek()
		if lbp, rbp := tailPrecedence(token.Type); lbp <= min_bp {
			break
		} else {
			parsedExpr = p.parseTailExpr(parsedExpr, rbp)
		}
	}
	return parsedExpr
}

func (p *parser) parseTailExpr(head ast.Expr, rbp int) ast.Expr {
	token := p.consume()
	switch token.Type {
	case lexer.ASSIGNMENT:
		rhs := p.parseExpr(rbp)
		return ast.AssignExpr{
			Assigne:       head,
			AssignedValue: rhs,
		}
	case lexer.PLUS, lexer.DASH, lexer.STAR, lexer.SLASH, lexer.PERCENT,
		lexer.LESS, lexer.LESS_EQUALS, lexer.GREATER, lexer.GREATER_EQUALS,
		lexer.PLUS_EQUALS, lexer.MINUS_EQUALS:
		tail := p.parseExpr(rbp)
		return ast.BinaryExpr{
			Lhs:      head,
			Operator: token,
			Rhs:      tail,
		}
	case lexer.OPEN_PAREN:
		return p.parseFuncCallExpr(head)
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
	case lexer.OPEN_PAREN:
		rbp := headPrecedence(token.Type)
		rhs := p.parseExpr(rbp)
		p.consume(lexer.CLOSE_PAREN)
		return ast.GroupExpr{
			Expr: rhs,
		}
	default:
		panic(fmt.Sprintf("Failed to parse head expression from token %v\n", token))
	}
}
