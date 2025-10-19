package frontend

import (
	"github.com/JoachimTislov/lite-jnc/frontend/lexer"
	"github.com/JoachimTislov/lite-jnc/frontend/parser"
)

type Frontend struct {
	*lexer.Lexer
	*parser.Parser
	Target string
}

func NewFrontend(path *string, language *string) (*Frontend, error) {
	lex, err := lexer.New(*path)
	if err != nil {
		return nil, err
	}
	parse, err := parser.New()
	if err != nil {
		return nil, err
	}
	return &Frontend{
		Lexer:  lex,
		Parser: parse,
		Target: *language,
	}, nil
}
