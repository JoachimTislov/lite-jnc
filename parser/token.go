package parser

type token struct {
	*pos
	value string
	kind  tokenKind
}

type pos struct {
	line  int
	start int
	end   int
}

type tokenKind int

const (
	NOT_SUPPORTED tokenKind = iota
	EOF

	// modifiers
	PUBLIC
	PRIVATE
	PROTECTED
	STATIC
	FINAL

	// keywords
	KEYWORD
	IDENTIFIER
	LITERAL
	PACKAGE
	IMPORT
	CLASS

	// errors
	CRITICAL
	ERROR
	WARNING
	INFO

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

	// method elements
	PARAMETER
	ARGUMENT

	// delimiter
	SEMICOLON
	COMMA
	DOT

	// operators
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
	TOKEN_QUOTE     = '"'
	TOKEN_COMMA     = ','
	TOKEN_SEMICOLON = ';'
	TOKEN_CPAREN    = ')'
	TOKEN_OPAREN    = '('
	TOKEN_OBRACE    = '{'
	TOKEN_CBRACE    = '}'
)

func (t tokenKind) String() string {
	switch t {
	case EOF:
		return "EOF"
	case ERROR:
		return "ERROR"
	case WARNING:
		return "WARNING"
	case CRITICAL:
		return "CRITICAL"
	case INFO:
		return "INFO"
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
	case DOT:
		return "period"
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
	case STRING:
		return "string"
	case INT:
		return "int"
	case FLOAT:
		return "float"
	case DOUBLE:
		return "double"
	case CHAR:
		return "char"
	case BOOLEAN:
		return "boolean"
	case EQUALS:
		return "equals"
	case ASSIGN:
		return "assign"
	case PLUS:
		return "plus"
	case MINUS:
		return "minus"
	case MULTIPLY:
		return "multiply"
	case DIVIDE:
		return "divide"
	case PERCENT:
		return "percent"
	case NOT:
		return "not"
	case LT:
		return "less than"
	case GT:
		return "greater than"
	default:
		return "unknown"
	}
}
