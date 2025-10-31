# Lite Java - Transpilers And Native Compilers

- Native Java compiler written in Go.
- Translates Java source to x86-64 assembly (Intel syntax) and compiles to standalone ELF executables for Linux.
- Supports a subset of Java features.
- Transpiles to multiple languages

## Resources and considerations

- Consider using: [LLVM - pkg](https://pkg.go.dev/tinygo.org/x/go-llvm) to support multiple architectures.
- [Lexical analysis in Go - Rob Pike](https://www.youtube.com/watch?v=HxaD_trXwRE)
    - [template lexer go source](https://go.dev/src/text/template/parse/lex.go)
- [aaronraff blog - how to write a lexer in Go](https://aaronraff.dev/blog/how-to-write-a-lexer-in-go)
