package lexer

type symbol string

type separator symbol

const (
	oParen    separator = "("
	cParen    separator = ")"
	obracket  separator = "{"
	cbracket  separator = "}"
	semiColon separator = ";"
)
