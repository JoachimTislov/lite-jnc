package main

import (
	"flag"
	"path"

	"github.com/JoachimTislov/lite-jnc/compile"
	"github.com/JoachimTislov/lite-jnc/env"
)

func main() {
	language := flag.String("l", "ELF", "Select language to either compile or transpile")
	source := path.Join(env.Home(), "projects/lite-jnc/src/Main.javaa")
	path := flag.String("p", source, "Path to the source file")
	// compile := flag.Bool("c", false, "Compile the transpiled language. Nothing happens when compiling directly to machine code")
	flag.Parse()

	switch *language {
	case "ELF":
		compile.CompileELF(*path)
	}
}
