package transpiler

import (
	"os"

	"github.com/JoachimTislov/lite-jnc/frontend"
)

type transpiler struct{ frontend *frontend.Frontend }
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

func New(frontend *frontend.Frontend) *transpiler {
	return &transpiler{frontend}
}

func (t *transpiler) Run(out string) *os.File {
	// Implement transpilation logic here
	return nil
}
