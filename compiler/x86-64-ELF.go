package compiler

import "log"

func (c *compiler) ELF() {
	for {
		_, err := c.Lexer.NextToken()
		if err != nil {
			log.Fatal(err)
			break
		}
		// symbol := c.Parser.Parse()
	}
}
