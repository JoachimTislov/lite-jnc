# Status overview

- [x] Basic state machine implementation
    - [x] lexer (strict, fault tolerant)
    - [x] parser (strict)
- [ ] Basic code generation setup

## Lexical analysis: Designed to be strict and robust

    - Errors
        - [x] Basic generic error handling
        - [ ] Critical
        - [ ] Warnings
        - [ ] Info
    - [x] Keywords
    - [x] Identifiers
    - [x] Modifiers
    - [x] Punctuations
    - [x] Class
        - [x] Method
            - [x] Parameters
            - [x] Return type
            - [x] Call
                - [x] Arguments
        - [x] Field
        - [ ] Variable
        - [ ] Constructor
        - [ ] Assingment
    - Literals
        - [x] String
        - [ ] Numeric
        - [ ] Boolean
        - [ ] Character
        - [ ] Null
    - Types
        - [x] Primitive types
        - [ ] Array types
        - [ ] Generic types
    - Operators
        - [ ] Arithmetic
        - [ ] Logical
        - [ ] Comparison

## Parsing to AST

    - [x] Classes
        - [x] Fields
        - [x] Methods
        - [ ] Variables
        - [ ] If-else
        - [ ] Switch
        - [ ] Loops
    - [ ] Interfaces
    - [ ] Inheritance
    - [ ] Enums
    - [ ] Packages

## Code generation

    - [ ] Native compilation
        - [ ] x86-64 Linux ELF
    - [ ] Intermediate representation
    - [ ] LLVM
    - Transpile
        - [ ] GO
        - [ ] JavaScript/TypeScript


