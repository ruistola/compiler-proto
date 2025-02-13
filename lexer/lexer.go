package lexer

import (
	"fmt"
	"regexp"
)

type TokenType int

const (
	EOF        TokenType = iota
	WHITESPACE           // pseudotype that will not participate in AST but is needed to unify the tokenization code
	COMMENT              // pseudotype that will not participate in AST but is kept as metadata
	SYMBOL               // pseudotype that will be refined into a symbol specific TokenType or IDENTIFIER
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

type tokenPattern struct {
	tokenType TokenType
	pattern   *regexp.Regexp
}

var tokenPatterns []tokenPattern = []tokenPattern{
	{WHITESPACE, regexp.MustCompile(`^\s+`)},
	{COMMENT, regexp.MustCompile(`^\/\/.*`)},
	{SYMBOL, regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*`)},
	{STRING, regexp.MustCompile(`^"[^"]*"`)},
	{NUMBER, regexp.MustCompile(`^[0-9]+(\.[0-9]+)?`)},
	{OPEN_BRACKET, regexp.MustCompile(`^\[`)},
	{CLOSE_BRACKET, regexp.MustCompile(`^\]`)},
	{OPEN_CURLY, regexp.MustCompile(`^\{`)},
	{CLOSE_CURLY, regexp.MustCompile(`^\}`)},
	{OPEN_PAREN, regexp.MustCompile(`^\(`)},
	{CLOSE_PAREN, regexp.MustCompile(`^\)`)},
	{EQUALS, regexp.MustCompile(`^==`)},
	{ASSIGNMENT, regexp.MustCompile(`^=`)},
	{NOT_EQUALS, regexp.MustCompile(`^!=`)},
	{NOT, regexp.MustCompile(`^!`)},
	{LESS, regexp.MustCompile(`^<`)},
	{LESS_EQUALS, regexp.MustCompile(`^<=`)},
	{GREATER, regexp.MustCompile(`^>`)},
	{GREATER_EQUALS, regexp.MustCompile(`^>=`)},
	{OR, regexp.MustCompile(`^\|\|`)},
	{AND, regexp.MustCompile(`^&&`)},
	{DOT, regexp.MustCompile(`^\.`)},
	{DOT_DOT, regexp.MustCompile(`^\.\.`)},
	{SEMI_COLON, regexp.MustCompile(`^;`)},
	{COLON, regexp.MustCompile(`^:`)},
	{QUESTION, regexp.MustCompile(`^\?`)},
	{COMMA, regexp.MustCompile(`^,`)},
	{PLUS_PLUS, regexp.MustCompile(`^\+\+`)},
	{MINUS_MINUS, regexp.MustCompile(`^--`)},
	{PLUS_EQUALS, regexp.MustCompile(`^\+=`)},
	{MINUS_EQUALS, regexp.MustCompile(`^-=`)},
	{NULLISH_ASSIGNMENT, regexp.MustCompile(`^\?\?=`)},
	{PLUS, regexp.MustCompile(`^\+`)},
	{DASH, regexp.MustCompile(`^-`)},
	{SLASH, regexp.MustCompile(`^/`)},
	{STAR, regexp.MustCompile(`^\*`)},
	{PERCENT, regexp.MustCompile(`^%`)},
}

var reservedKeywords map[string]TokenType = map[string]TokenType{
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
	case WHITESPACE:
		return "whitespace"
	case COMMENT:
		return "comment"
	case SYMBOL:
		return "symbol"
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

func tryMatchPattern(src string, re *regexp.Regexp, tokenType TokenType) (int, Token) {
	loc := re.FindStringIndex(src)
	if loc == nil {
		return 0, Token{}
	}

	if loc[0] > 0 {
		panic(fmt.Sprintf("Internal error: regex matched at non-zero index %d!", loc[0]))
	}

	match := src[:loc[1]]
	length := loc[1]

	if tokenType != SYMBOL {
		return length, Token{
			Type:  tokenType,
			Value: match,
		}
	} else if keywordTokenType, found := reservedKeywords[match]; found {
		return length, Token{
			Type:  keywordTokenType,
			Value: match,
		}
	} else {
		return length, Token{
			Type:  IDENTIFIER,
			Value: match,
		}
	}
}

func Tokenize(src string) []Token {
	pos := 0
	tokens := make([]Token, 0)

	for pos < len(src) {
		remainingSrc := src[pos:]
		for _, tp := range tokenPatterns {
			if length, newToken := tryMatchPattern(remainingSrc, tp.pattern, tp.tokenType); length != 0 {
				if newToken.Type != WHITESPACE {
					tokens = append(tokens, newToken)
				}
				pos += length
				break
			}
		}
	}

	tokens = append(tokens, Token{
		Type:  EOF,
		Value: "EOF",
	})
	return tokens
}
