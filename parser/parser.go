package parser

import (
	"fmt"
	"jru-test/ast"
	"jru-test/lexer"
	"slices"
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
	case lexer.NUMBER, lexer.STRING, lexer.WORD:
		return 1
	case lexer.PLUS, lexer.DASH:
		return 9
	default:
		panic(fmt.Sprintf("Cannot determine binding power for '%s' as a head token", tokenType))
	}
}

func tailPrecedence(tokenType lexer.TokenType) (int, int) {
	switch tokenType {
	case lexer.EOF, lexer.SEMI_COLON, lexer.CLOSE_PAREN, lexer.COMMA, lexer.CLOSE_CURLY, lexer.CLOSE_BRACKET:
		return 0, 0
	case lexer.ASSIGNMENT, lexer.PLUS_EQUALS, lexer.MINUS_EQUALS:
		return 1, 2
	case lexer.EQUALS, lexer.NOT_EQUALS:
		return 3, 4
	case lexer.LESS, lexer.LESS_EQUALS, lexer.GREATER, lexer.GREATER_EQUALS:
		return 6, 5
	case lexer.PLUS, lexer.DASH:
		return 8, 7
	case lexer.STAR, lexer.SLASH, lexer.PERCENT:
		return 10, 9
	case lexer.OPEN_CURLY:
		return 11, 0
	case lexer.OPEN_PAREN, lexer.OPEN_BRACKET:
		return 12, 0
	case lexer.DOT:
		return 14, 13
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
	case lexer.STRUCT:
		return p.parseStructDeclStmt()
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

func (p *parser) parseExpr(min_bp int) ast.Expr {
	token := p.consume()
	parsedExpr := p.parseHeadExpr(token)

	for {
		nextToken := p.peek()
		if lbp, rbp := tailPrecedence(nextToken.Type); lbp <= min_bp {
			break
		} else {
			parsedExpr = p.parseTailExpr(parsedExpr, rbp)
		}
	}

	return parsedExpr
}

func (p *parser) parseHeadExpr(token lexer.Token) ast.Expr {
	switch token.Type {
	case lexer.NUMBER:
		return ast.NumberLiteralExpr{
			Number: token.Value,
		}
	case lexer.STRING:
		return ast.StringLiteralExpr{
			String: token.Value,
		}
	case lexer.IDENTIFIER:
		return ast.IdentExpr{
			Name: token.Value,
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

func (p *parser) parseTailExpr(head ast.Expr, rbp int) ast.Expr {
	token := p.consume()
	switch token.Type {
	case lexer.ASSIGNMENT:
		rhs := p.parseExpr(rbp)
		return ast.AssignExpr{
			Assigne:       head,
			AssignedValue: rhs,
		}
	case lexer.PLUS,
		lexer.DASH,
		lexer.STAR,
		lexer.SLASH,
		lexer.PERCENT,
		lexer.LESS,
		lexer.LESS_EQUALS,
		lexer.GREATER,
		lexer.GREATER_EQUALS,
		lexer.PLUS_EQUALS,
		lexer.MINUS_EQUALS:

		tail := p.parseExpr(rbp)
		return ast.BinaryExpr{
			Lhs:      head,
			Operator: token,
			Rhs:      tail,
		}
	case lexer.OPEN_PAREN:
		return p.parseFuncCallExpr(head)
	case lexer.OPEN_CURLY:
		return p.parseStructLiteralExpr(head)
	case lexer.OPEN_BRACKET:
		return p.parseArrayIndexExpr(head)
	case lexer.DOT:
		return p.parseStructMemberExpr(head)
	default:
		panic(fmt.Sprintf("Failed to parse tail expression from token %v\n", token))
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
	namedType := ast.NamedType{
		TypeName: name,
	}
	if p.peek().Type == lexer.OPEN_BRACKET {
		return p.parseArrayType(namedType)
	}
	return namedType
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

func (p *parser) parseFuncDeclStmt() ast.FuncDeclStmt {
	p.consume(lexer.FUNC)
	name := p.consume(lexer.IDENTIFIER).Value

	params := make([]ast.TypedIdent, 0)
	p.consume(lexer.OPEN_PAREN)

	for p.peek().Type != lexer.CLOSE_PAREN {
		paramName := p.consume(lexer.IDENTIFIER).Value
		p.consume(lexer.COLON)
		paramType := p.parseType()
		newParam := ast.TypedIdent{
			Name: paramName,
			Type: paramType,
		}
		params = append(params, newParam)

		if p.peek().Type == lexer.COMMA {
			p.consume(lexer.COMMA)
			if p.peek().Type != lexer.IDENTIFIER {
				panic(fmt.Sprintf("Expected identifier after comma in function parameter list, found %s", p.peek().Type))
			}
		}
	}
	p.consume(lexer.CLOSE_PAREN)

	var returnType ast.Type
	if p.peek().Type == lexer.COLON {
		p.consume(lexer.COLON)
		returnType = p.parseType()
	}

	funcBody := p.parseBlockStmt()

	return ast.FuncDeclStmt{
		Name:       name,
		Parameters: params,
		ReturnType: returnType,
		Body:       funcBody,
	}
}

func (p *parser) parseStructDeclStmt() ast.StructDeclStmt {
	p.consume(lexer.STRUCT)
	name := p.consume(lexer.IDENTIFIER).Value
	p.consume(lexer.OPEN_CURLY)

	members := make([]ast.TypedIdent, 0)
	for p.peek().Type != lexer.CLOSE_CURLY {
		memberName := p.consume(lexer.IDENTIFIER).Value
		p.consume(lexer.COLON)
		memberType := p.parseType()

		newMember := ast.TypedIdent{
			Name: memberName,
			Type: memberType,
		}
		members = append(members, newMember)
	}
	p.consume(lexer.CLOSE_CURLY)

	return ast.StructDeclStmt{
		Name:    name,
		Members: members,
	}
}

func (p *parser) parseIfStmt() ast.Stmt {
	p.consume(lexer.IF)

	p.consume(lexer.OPEN_PAREN)
	cond := p.parseExpr(0)
	p.consume(lexer.CLOSE_PAREN)

	thenStmt := p.parseStmt()

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
		args = append(args, p.parseExpr(0))
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

func (p *parser) parseStructLiteralExpr(left ast.Expr) ast.StructLiteralExpr {
	members := []ast.AssignExpr{}

	for p.peek().Type != lexer.CLOSE_CURLY {
		memberName := p.consume(lexer.IDENTIFIER).Value
		p.consume(lexer.COLON)
		value := p.parseExpr(0)
		members = append(members, ast.AssignExpr{
			Assigne: ast.IdentExpr{
				Name: memberName,
			},
			AssignedValue: value,
		})
		p.consume(lexer.COMMA)
	}

	p.consume(lexer.CLOSE_CURLY)

	return ast.StructLiteralExpr{
		Struct:  left,
		Members: members,
	}
}

func (p *parser) parseStructMemberExpr(left ast.Expr) ast.StructMemberExpr {
	return ast.StructMemberExpr{
		Struct: left,
		Member: ast.IdentExpr{
			Name: p.consume(lexer.IDENTIFIER).Value,
		},
	}
}

func (p *parser) parseArrayIndexExpr(left ast.Expr) ast.ArrayIndexExpr {
	indexExpr := p.parseExpr(0)
	p.consume(lexer.CLOSE_BRACKET)

	return ast.ArrayIndexExpr{
		Array: left,
		Index: indexExpr,
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
