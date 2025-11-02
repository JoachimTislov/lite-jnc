package compiler

import (
	"fmt"
	"log"
)

func (c *compiler) ELF() {
	ast, err := c.Parser.Parse()
	fmt.Printf("AST: %+v\n", ast)
	if err != nil {
		log.Print(err)
	}
}
