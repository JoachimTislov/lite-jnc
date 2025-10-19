package lexer

type LexError int

const (
	EOF LexError = iota
)

var errorMessages = map[LexError]string{
	EOF: "End of source",
}

func (e LexError) Error() string {
	msg := "Lexer error"
	if errorMsg, ok := errorMessages[e]; ok {
		msg = errorMsg
	}
	return msg
}

func (e LexError) String() string {
	return e.Error()
}
