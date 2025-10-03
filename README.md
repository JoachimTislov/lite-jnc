# Lite Java - Native Compiler

- Native Java compiler written in Go. 
- Translates Java source to x86-64 assembly (Intel syntax) and compiles to standalone ELF executables for Linux. 
- Supports a subset of Java features.

## Resources and considerations

- Consider using: [LLVM - pkg](https://pkg.go.dev/tinygo.org/x/go-llvm) to support multiple architectures.
- Should write transpilers to prevent a challenging climb.
   - Java -> ( C or Go )
   - Which can then be compiled to machine code with a C or Go compiler.
- The difficult part is to write a well-designed lexer and parser. Afterwards, it can be challenging to translate to machine code, but there are multiple ways to do it, and I should do most of them instead of narrowing down to supporting every feature in Java.
