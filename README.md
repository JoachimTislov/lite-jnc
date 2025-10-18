# Lite Java - Transpilers And Native Compilers

- Native Java compiler written in Go.
- Translates Java source to x86-64 assembly (Intel syntax) and compiles to standalone ELF executables for Linux.
- Supports a subset of Java features.
- Transpiles to multiple languages

## Resources and considerations

- Consider using: [LLVM - pkg](https://pkg.go.dev/tinygo.org/x/go-llvm) to support multiple architectures.
