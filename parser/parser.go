package parser

import (
	"fmt"
)

type Parser struct {
	*lexer
	Target string
}

func New(path string, language string) (*Parser, error) {
	lexer, err := newLexer(path)
	if err != nil {
		return nil, err
	}
	return &Parser{
		lexer:  lexer,
		Target: language,
	}, nil
}

func (p *Parser) Parse() {
	go p.lexer.run()
	for t := range p.lexer.tokens {
		if t.kind == EOF {
			break
		}
		fmt.Println(t)
	}
	// for {
	// 	t := p.lexer.nextToken()
	// 	if t.kind == EOF {
	// 		close(p.lexer.tokens)
	// 		break
	// 	}
	// 	fmt.Println(t)
	// }
}
