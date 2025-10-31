package transpiler

import (
	"os"

	"github.com/JoachimTislov/lite-jnc/parser"
)

type transpiled int

const (
	Go transpiled = iota
	Python
	JavaScript
	C
	CSharp
)

var languages = map[string]transpiled{
	"Go":         Go,
	"Python":     Python,
	"JavaScript": JavaScript,
	"C":          C,
	"CSharp":     CSharp,
}

func Supports(lang string) bool {
	_, exists := languages[lang]
	return exists
}

type transpiler struct{ *parser.Parser }

func New(p *parser.Parser) *transpiler {
	return &transpiler{p}
}

func (t *transpiler) Run(out string) *os.File {
	// Implement transpilation logic here
	return nil
}
