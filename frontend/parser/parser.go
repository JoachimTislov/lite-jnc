package parser

import "github.com/JoachimTislov/lite-jnc/spec"

type Parser struct{}

func New() (*Parser, error) {
	return &Parser{}, nil
}

func (p *Parser) Parse(token spec.Token) string {
	switch token.Type {
	case spec.ADD:

	}
	return ""
}
