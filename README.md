# Lite Java - Transpilers And Native Compilers

- Native Java compiler written in Go.
- Translates Java source to x86-64 assembly (Intel syntax) and compiles to standalone ELF executables for Linux.
- Supports a subset of Java features.
- Transpiles to multiple languages

## Resources and considerations

- Using a Makefile to ease the build process and commands. I recommend viewing [Makefile tutorial](https://makefiletutorial.com/) if unfamiliar with Makefiles.
- Compiler uses [LLVM - pkg](https://pkg.go.dev/tinygo.org/x/go-llvm) to support multiple architectures.
    - [LLVM Language Reference Manual](https://llvm.org/docs/LangRef.html)
- [Lexical analysis in Go - Rob Pike](https://www.youtube.com/watch?v=HxaD_trXwRE)
    - [template lexer go source](https://go.dev/src/text/template/parse/lex.go)
- [aaronraff blog - how to write a lexer in Go](https://aaronraff.dev/blog/how-to-write-a-lexer-in-go)
- [Go template package - lex.go](https://go.dev/src/text/template/parse/lex.go)
- [Crafting Interpreters - Bob Nystrom](https://craftinginterpreters.com/)
- Java references
    - [Language Specification](https://docs.oracle.com/javase/specs/jls/se17/html/index.html)
    - [Cheatsheet](https://introcs.cs.princeton.edu/java/11cheatsheet/)
- [Tree-sitter](https://tree-sitter.github.io/tree-sitter/), [Java bindings](https://github.com/tree-sitter/java-tree-sitter)

## Project status

[Status overview](./Status.md)

## TODOs

- Handle missing ending punctionations convieniently
    - Either to the end of file or if the expected punctation is not found
    - A missing '"' or '`' can't be handled by any other means than reaching the end of the file. Logically the lexer can't know when the user intends to end the string. Same goes for alot of other unclosed punctuations:
        - `/*`
        - `(`, `{`, `[` - possible to avoid end of file before realizing that there is no closing `)`
- Standardize the AST
- Implement the LLVM pk
    - Add Intermediate representation?
