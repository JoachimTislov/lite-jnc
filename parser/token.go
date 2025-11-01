package parser

import "fmt"

type token struct {
	pos
	value string
	kind  tokenKind
}

func (t token) String() string {
	rang := fmt.Sprintf("%d-%d", t.pos.start, t.pos.end)
	if t.pos.start == t.pos.end {
		rang = fmt.Sprintf("%d", t.pos.end)
	}
	return fmt.Sprintf("kind: %-12s; value: %-15s %d:%s", t.kind, t.value, t.pos.line, rang)
}

type pos struct {
	line  int
	start int
	end   int
}

type tokenKind int

const (
	ILLEGAL tokenKind = iota
	IDENTIFIER
	MODIFIER
	LITERAL
	TYPE
	KEYWORD
	EOF
	ERROR
	OPUNCTUATION
	CPUNCTUATION
	PARAMETER
	DELIMITER
	OPERAND
	OPERATOR
)

const (
	TOKEN_EOF        = "EOF"
	TOKEN_IDENTIFIER = "identifier"
	TOKEN_QUOTE      = '"'
)

var tokens = map[tokenKind]string{
	EOF:          TOKEN_EOF,
	IDENTIFIER:   TOKEN_IDENTIFIER,
	ILLEGAL:      "invalid",
	MODIFIER:     "modifier",
	TYPE:         "type",
	KEYWORD:      "keyword",
	LITERAL:      "literal",
	OPUNCTUATION: "cpunctuation",
	CPUNCTUATION: "opunctuation",
	PARAMETER:    "parameter",
	DELIMITER:    "delimiter",
	OPERAND:      "operand",
	OPERATOR:     "operator",
}

func (t tokenKind) String() string {
	if str, ok := tokens[t]; ok {
		return str
	}
	return TOKEN_IDENTIFIER
}
