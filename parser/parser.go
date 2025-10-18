package parser

type symbol string

const (
	class symbol = "class"
	void  symbol = "void"
)

type separator symbol

const (
	oParen    separator = "("
	cParen    separator = ")"
	obracket  separator = "{"
	cbracket  separator = "}"
	semiColon separator = ";"
)

type modifier symbol

const (
	pub    modifier = "public"
	priv   modifier = "private"
	proc   modifier = "protected"
	static modifier = "static"
)
