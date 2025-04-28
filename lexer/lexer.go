package lexer

import (
	"fmt"
	"regexp"
)

type TokenType int

const (
	EOF        TokenType = iota
	WHITESPACE           // pseudotype that will not participate in AST but is needed to unify the tokenization code
	WORD                 // pseudotype that will be refined into a keyword or IDENTIFIER
	COMMENT              // pseudotype that will not participate in AST but may in the future be kept as metadata
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

	// Equivalence
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
	SEMI_COLON
	COLON
	COMMA

	// Shorthand
	PLUS_EQUALS
	MINUS_EQUALS

	//Maths
	PLUS
	DASH
	SLASH
	STAR
	PERCENT

	// Reserved Keywords
	LET
	STRUCT
	TRUE
	FALSE
	FUNC
	IF
	ELSE
	FOR
	RETURN

	// Misc
	NUM_TOKENS
)

type tokenPattern struct {
	tokenType TokenType
	pattern   *regexp.Regexp
}

var tokenPatterns []tokenPattern = []tokenPattern{
	{WHITESPACE, regexp.MustCompile(`^\s+`)},
	{WORD, regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*`)},
	{COMMENT, regexp.MustCompile(`^\/\/.*`)},
	{NUMBER, regexp.MustCompile(`^[0-9]+(\.[0-9]+)?`)},
	{STRING, regexp.MustCompile(`^"[^"]*"`)},
	{OPEN_BRACKET, regexp.MustCompile(`^\[`)},
	{CLOSE_BRACKET, regexp.MustCompile(`^\]`)},
	{OPEN_CURLY, regexp.MustCompile(`^\{`)},
	{CLOSE_CURLY, regexp.MustCompile(`^\}`)},
	{OPEN_PAREN, regexp.MustCompile(`^\(`)},
	{CLOSE_PAREN, regexp.MustCompile(`^\)`)},
	{EQUALS, regexp.MustCompile(`^==`)},
	{NOT_EQUALS, regexp.MustCompile(`^!=`)},
	{ASSIGNMENT, regexp.MustCompile(`^=`)},
	{NOT, regexp.MustCompile(`^!`)},
	{LESS_EQUALS, regexp.MustCompile(`^<=`)},
	{LESS, regexp.MustCompile(`^<`)},
	{GREATER_EQUALS, regexp.MustCompile(`^>=`)},
	{GREATER, regexp.MustCompile(`^>`)},
	{OR, regexp.MustCompile(`^\|\|`)},
	{AND, regexp.MustCompile(`^&&`)},
	{DOT, regexp.MustCompile(`^\.`)},
	{SEMI_COLON, regexp.MustCompile(`^;`)},
	{COLON, regexp.MustCompile(`^:`)},
	{COMMA, regexp.MustCompile(`^,`)},
	{PLUS_EQUALS, regexp.MustCompile(`^\+=`)},
	{MINUS_EQUALS, regexp.MustCompile(`^-=`)},
	{PLUS, regexp.MustCompile(`^\+`)},
	{DASH, regexp.MustCompile(`^-`)},
	{SLASH, regexp.MustCompile(`^/`)},
	{STAR, regexp.MustCompile(`^\*`)},
	{PERCENT, regexp.MustCompile(`^%`)},
}

var reservedKeywords map[string]TokenType = map[string]TokenType{
	"let":    LET,
	"struct": STRUCT,
	"true":   TRUE,
	"false":  FALSE,
	"func":   FUNC,
	"if":     IF,
	"else":   ELSE,
	"for":    FOR,
	"return": RETURN,
}

func (tokenType TokenType) String() string {
	switch tokenType {
	case EOF:
		return "eof"
	case WHITESPACE:
		return "whitespace"
	case WORD:
		return "word"
	case COMMENT:
		return "comment"
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
	case SEMI_COLON:
		return "semi_colon"
	case COLON:
		return "colon"
	case COMMA:
		return "comma"
	case PLUS_EQUALS:
		return "plus_equals"
	case MINUS_EQUALS:
		return "minus_equals"
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
	case FUNC:
		return "func"
	case IF:
		return "if"
	case ELSE:
		return "else"
	case FOR:
		return "for"
	case STRUCT:
		return "struct"
	case RETURN:
		return "return"
	default:
		return fmt.Sprintf("unknown(%d)", tokenType)
	}
}

type Token struct {
	Type  TokenType
	Value string
}

func tryMatchPattern(src string, re *regexp.Regexp, tokenType TokenType) (int, Token) {
	matchRange := re.FindStringIndex(src)
	if matchRange == nil {
		return 0, Token{}
	}

	if matchRange[0] > 0 {
		panic(fmt.Sprintf("Internal error: regex matched at non-zero index %d!", matchRange[0]))
	}

	match := src[:matchRange[1]]
	length := matchRange[1]

	if tokenType != WORD {
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
				if newToken.Type != WHITESPACE && newToken.Type != COMMENT {
					tokens = append(tokens, newToken)
				}
				pos += length
				break
			}
		}
	}

	return tokens
}
