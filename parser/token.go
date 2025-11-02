package parser

import "fmt"

type token struct {
	*pos
	value string
	kind  tokenKind
}

func (t *token) String() string {
	if t.kind == ERROR {
		return fmt.Sprintf("ERROR(%v): \n\t%s", t.pos, t.value)
	}
	return fmt.Sprintf("kind: %-20s value: %-20s pos: %v", t.kind, t.value, t.pos)
}

type pos struct {
	line  int
	start int
	end   int
}

type tokenKind int

const (
	NOT_SUPPORTED tokenKind = iota
	IDENTIFIER

	PUBLIC
	PRIVATE
	PROTECTED
	STATIC
	FINAL

	// keywords
	CLASS

	LITERAL
	KEYWORD
	EOF
	ERROR

	// punctuations
	OPAREN
	CPAREN
	OBRACE
	CBRACE
	OBRACKET
	CBRACKET

	// types
	PRIMITIVE
	REFERENCE

	// others
	PARAMETER
	ARGUMENT
	DELIMITER
	OPERAND
	OPERATOR
)

const (
	// strings
	TOKEN_EOF   = "EOF"
	TOKEN_CLASS = "class"
	// runes
	TOKEN_QUOTE  = '"'
	TOKEN_COMMA  = ','
	TOKEN_CPAREN = ')'
	TOKEN_OPAREN = '('
	TOKEN_OBRACE = '{'
	TOKEN_CBRACE = '}'
	SPACE        = ' '
)

func (t tokenKind) String() string {
	switch t {
	case EOF:
		return TOKEN_EOF
	case ERROR:
		return "error"
	case IDENTIFIER:
		return "identifier"
	case CLASS:
		return "class"
	case ARGUMENT:
		return "argument"
	case NOT_SUPPORTED:
		return "not supported"
	case PUBLIC, PRIVATE, PROTECTED:
		return "access modifier"
	case STATIC, FINAL:
		return "persistence modifier"
	case PRIMITIVE:
		return "primitive type"
	case REFERENCE:
		return "reference"
	case KEYWORD:
		return "keyword"
	case LITERAL:
		return "literal"
	case OPAREN:
		return "oparen"
	case CPAREN:
		return "cparen"
	case OBRACE:
		return "obrace"
	case CBRACE:
		return "cbrace"
	case OBRACKET:
		return "obracket"
	case CBRACKET:
		return "cbracket"
	case PARAMETER:
		return "parameter"
	case DELIMITER:
		return "delimiter"
	case OPERAND:
		return "operand"
	case OPERATOR:
		return "operator"
	default:
		return "unknown"
	}
}
