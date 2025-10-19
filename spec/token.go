package spec

type Token struct {
	Pos
	Value string
	Type  tokenType
}

type Pos struct {
	Line  int
	Start int
	End   int
}

type tokenType int

const (
	// Keywords
	CLASS tokenType = iota
	ABSTRACT
	INTERFACE
	VOID

	// Separators
	SEMI
	LPAREN
	RPAREN
	LBRACK
	RBRACK
	LBRACE
	RBRACE

	// Modifiers
	PUB
	PRIV
	PROTEC
	STATIC

	// Logical
	ASSIGN
	EQUAL
	ADD
	SUB
	MULT
	DIV
)

var tokens = map[tokenType]string{
	CLASS:     "class",
	ABSTRACT:  "abstract",
	INTERFACE: "interface",
	VOID:      "void",
	SEMI:      ";",
	LPAREN:    "(",
	RPAREN:    ")",
	LBRACE:    "{",
	RBRACE:    "}",
	LBRACK:    "[",
	RBRACK:    "]",
	PUB:       "public",
	PRIV:      "private",
	PROTEC:    "protected",
	STATIC:    "static",
	ASSIGN:    "=",
	EQUAL:     "==",
	ADD:       "+",
	SUB:       "-",
	MULT:      "*",
	DIV:       "/",
}

func (t tokenType) String() string {
	return tokens[t]
}
