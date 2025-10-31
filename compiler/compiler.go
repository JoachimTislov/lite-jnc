package compiler

import (
	"log"
	"os"

	"github.com/JoachimTislov/lite-jnc/parser"
)

type Compiled int

const (
	X86_64ELF Compiled = iota
)

var languages = map[string]Compiled{
	"ELF": X86_64ELF,
}

type compiler struct {
	*parser.Parser
}

func New(p *parser.Parser) *compiler {
	return &compiler{p}
}

func (c *compiler) Run(out string) *os.File {
	switch languages[c.Target] {
	case X86_64ELF:
		c.ELF()
	default:
		log.Fatalf("Unsupported language: %s", c.Target)
		// Handle unsupported language
	}
	return nil
}
