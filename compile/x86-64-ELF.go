package compile

import (
	"github.com/JoachimTislov/lite-jnc/lexer"
	"log"
)

func CompileELF(src string) {
	lexer, err := lexer.New(src)
	if err != nil {
		log.Fatal(err)
	}
	print(lexer.NextToken())
	print(lexer.NextToken())
}
