package lexer_test

import (
	"path"
	"testing"

	"github.com/JoachimTislov/lite-jnc/env"
	"github.com/JoachimTislov/lite-jnc/lexer"
)

func TestLexer(t *testing.T) {
	path := path.Join(env.Home(), "projects/lite-jnc/src/Main.javaa")
	lexer, err := lexer.New(path)
	if err != nil {
		t.Error(err)
	}
	first := lexer.NextToken()
	second := lexer.NextToken()
	if first == second {
		t.Error("First and second token is the same, state is not working")
	}
}
