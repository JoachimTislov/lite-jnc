package main

import (
	"flag"
	"log"
	"path"

	"github.com/JoachimTislov/lite-jnc/compiler"
	"github.com/JoachimTislov/lite-jnc/env"
	"github.com/JoachimTislov/lite-jnc/parser"
	"github.com/JoachimTislov/lite-jnc/spec"
	"github.com/JoachimTislov/lite-jnc/transpiler"
)

func main() {
	language := flag.String("l", "ELF", "Select language to either compile or transpile")
	source := path.Join(env.Home(), "projects/lite-jnc/src/Main.javaa")
	path := flag.String("p", source, "Path to the source file")
	out := flag.String("o", "out", "name of output file")
	// compile := flag.Bool("c", false, "Compile the transpiled language. Nothing happens when compiling directly to machine code")
	flag.Parse()

	p, err := parser.New(*path, *language)
	if err != nil {
		log.Fatal(err)
	}

	var runner spec.Runner
	if transpiler.Supports(*language) {
		runner = transpiler.New(p)
	} else {
		runner = compiler.New(p)
	}
	runner.Run(*out)
}
