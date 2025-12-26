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
	PACKAGE
	IMPORT
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
	VOID
	INT
	FLOAT
	STRING
	BOOLEAN
	DOUBLE
	CHAR

	// others
	PARAMETER
	ARGUMENT
	DELIMITER
	SEMICOLON
	COMMA
	PERIOD
	QUESTION
	COLON

	OPERAND
	OPERATOR

	EQUALS
	ASSIGN
	PLUS
	MINUS
	MULTIPLY
	DIVIDE
	PERCENT
	NOT
	LT
	GT
)

const (
	// strings
	TOKEN_EOF   = "EOF"
	TOKEN_CLASS = "class"
	// runes
	TOKEN_QUOTE     = '"'
	TOKEN_COMMA     = ','
	TOKEN_SEMICOLON = ';'
	TOKEN_CPAREN    = ')'
	TOKEN_OPAREN    = '('
	TOKEN_OBRACE    = '{'
	TOKEN_CBRACE    = '}'
	SPACE           = ' '
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
	case SEMICOLON:
		return "semicolon"
	case COMMA:
		return "comma"
	case PERIOD:
		return "period"
	case QUESTION:
		return "question"
	case COLON:
		return "colon"
	case PACKAGE:
		return "package"
	case VOID:
		return "void"
	case IMPORT:
		return "import"
	case PUBLIC:
		return "public"
	case PRIVATE:
		return "private"
	case PROTECTED:
		return "protected"
	case STATIC:
		return "static"
	case FINAL:
		return "final"
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
