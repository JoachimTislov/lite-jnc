package compiler

import (
	"fmt"
)

func (c *compiler) ELF() {
	ast, errors := c.Parser.Parse()
	fmt.Println(ast)
	for _, err := range errors {
		fmt.Println(err)
	}
}
