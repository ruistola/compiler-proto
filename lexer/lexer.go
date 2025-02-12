package lexer

import (
	"fmt"
	"regexp"
)

type TokenType int

const (
	EOF TokenType = iota
	NULL
	TRUE
	FALSE
	NUMBER
	STRING
	IDENTIFIER

	// Grouping & Braces
	OPEN_BRACKET
	CLOSE_BRACKET
	OPEN_CURLY
	CLOSE_CURLY
	OPEN_PAREN
	CLOSE_PAREN

	// Equivilance
	ASSIGNMENT
	EQUALS
	NOT_EQUALS
	NOT

	// Conditional
	LESS
	LESS_EQUALS
	GREATER
	GREATER_EQUALS

	// Logical
	OR
	AND

	// Symbols
	DOT
	DOT_DOT
	SEMI_COLON
	COLON
	QUESTION
	COMMA

	// Shorthand
	PLUS_PLUS
	MINUS_MINUS
	PLUS_EQUALS
	MINUS_EQUALS
	NULLISH_ASSIGNMENT // ??=

	//Maths
	PLUS
	DASH
	SLASH
	STAR
	PERCENT

	// Reserved Keywords
	LET
	CONST
	CLASS
	NEW
	IMPORT
	FROM
	FN
	IF
	ELSE
	FOREACH
	WHILE
	FOR
	EXPORT
	TYPEOF
	IN

	// Misc
	NUM_TOKENS
)

var keywordTokens map[string]TokenType = map[string]TokenType{
	"true":    TRUE,
	"false":   FALSE,
	"null":    NULL,
	"let":     LET,
	"const":   CONST,
	"class":   CLASS,
	"new":     NEW,
	"import":  IMPORT,
	"from":    FROM,
	"fn":      FN,
	"if":      IF,
	"else":    ELSE,
	"foreach": FOREACH,
	"while":   WHILE,
	"for":     FOR,
	"export":  EXPORT,
	"typeof":  TYPEOF,
	"in":      IN,
}

func (tokenType TokenType) String() string {
	switch tokenType {
	case EOF:
		return "eof"
	case NULL:
		return "null"
	case NUMBER:
		return "number"
	case STRING:
		return "string"
	case TRUE:
		return "true"
	case FALSE:
		return "false"
	case IDENTIFIER:
		return "identifier"
	case OPEN_BRACKET:
		return "open_bracket"
	case CLOSE_BRACKET:
		return "close_bracket"
	case OPEN_CURLY:
		return "open_curly"
	case CLOSE_CURLY:
		return "close_curly"
	case OPEN_PAREN:
		return "open_paren"
	case CLOSE_PAREN:
		return "close_paren"
	case ASSIGNMENT:
		return "assignment"
	case EQUALS:
		return "equals"
	case NOT_EQUALS:
		return "not_equals"
	case NOT:
		return "not"
	case LESS:
		return "less"
	case LESS_EQUALS:
		return "less_equals"
	case GREATER:
		return "greater"
	case GREATER_EQUALS:
		return "greater_equals"
	case OR:
		return "or"
	case AND:
		return "and"
	case DOT:
		return "dot"
	case DOT_DOT:
		return "dot_dot"
	case SEMI_COLON:
		return "semi_colon"
	case COLON:
		return "colon"
	case QUESTION:
		return "question"
	case COMMA:
		return "comma"
	case PLUS_PLUS:
		return "plus_plus"
	case MINUS_MINUS:
		return "minus_minus"
	case PLUS_EQUALS:
		return "plus_equals"
	case MINUS_EQUALS:
		return "minus_equals"
	case NULLISH_ASSIGNMENT:
		return "nullish_assignment"
	case PLUS:
		return "plus"
	case DASH:
		return "dash"
	case SLASH:
		return "slash"
	case STAR:
		return "star"
	case PERCENT:
		return "percent"
	case LET:
		return "let"
	case CONST:
		return "const"
	case CLASS:
		return "class"
	case NEW:
		return "new"
	case IMPORT:
		return "import"
	case FROM:
		return "from"
	case FN:
		return "fn"
	case IF:
		return "if"
	case ELSE:
		return "else"
	case FOREACH:
		return "foreach"
	case FOR:
		return "for"
	case WHILE:
		return "while"
	case EXPORT:
		return "export"
	case IN:
		return "in"
	default:
		return fmt.Sprintf("unknown(%d)", tokenType)
	}
}

type Token struct {
	Type  TokenType
	Value string
}

func (token Token) Debug() {
	switch token.Type {
	case IDENTIFIER, NUMBER, STRING:
		fmt.Printf("%s:%s", token.Type, token.Value)
	default:
		fmt.Printf("%s", token.Type)
	}
}

var (
	reWhitespace    = regexp.MustCompile(`^\s+`)
	reStringLiteral = regexp.MustCompile(`^"[^"]*"`)
	reNumber        = regexp.MustCompile(`^[0-9]+(\.[0-9]+)?`)
)

type patternHandler struct {
	pattern *regexp.Regexp
	handler func(lexer *lexer_t, pattern *regexp.Regexp) Token
}

type lexer_t struct {
	pos             int
	src             string
	tokens          []Token
	patternHandlers []patternHandler
}

func Tokenize(src string) []Token {
	var pos int
	tokens := make([]Token, 0)

	for pos < len(src) {
		remainingSrc := src[pos:]
		if loc := reWhitespace.FindStringIndex(remainingSrc); loc != nil {
			if loc[0] > 0 {
				panic(fmt.Sprintf("Internal error: anchored regex matched at index %d!", loc[0]))
			}
			// Consume whitespace
			pos += loc[1]
		} else if loc := reStringLiteral.FindStringIndex(remainingSrc); loc != nil {
			if loc[0] > 0 {
				panic(fmt.Sprintf("Internal error: anchored regex matched at index %d!", loc[0]))
			}
			// The token is a string literal
			tokens = append(tokens, Token{
				Type:  STRING,
				Value: remainingSrc[:loc[1]],
			})
			pos += loc[1]
		} else if loc := reNumber.FindStringIndex(remainingSrc); loc != nil {
			if loc[0] > 0 {
				panic(fmt.Sprintf("Internal error: anchored regex matched at index %d!", loc[0]))
			}
			// The token is a number
			tokens = append(tokens, Token{
				Type:  NUMBER,
				Value: remainingSrc[:loc[1]],
			})
			pos += loc[1]
		} // else if ...
	}

	tokens = append(tokens, Token{
		Type:  EOF,
		Value: "EOF",
	})
	return tokens
}
