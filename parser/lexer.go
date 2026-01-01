package parser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

type lexStateFn func(*lexer) lexStateFn

type lexer struct {
	reader    *bufio.Reader
	line      int
	column    int
	runes     []rune
	state     lexStateFn
	prevToken string
	tokens    chan *token
	running   bool
	cleanup   func()
}

func newLexer(src string) (*lexer, error) {
	file, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	chanTokens := make(chan *token)
	return &lexer{
		reader: bufio.NewReader(file),
		line:   1,
		column: 1,
		state:  lexClass,
		tokens: chanTokens,
		cleanup: func() {
			file.Close()
			close(chanTokens)
		},
	}, nil
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r' || r == '\n'
}

func isNumber(r rune) bool {
	return '0' <= r && r <= '9'
}

func (l *lexer) enforceWhitespace(kind tokenKind) {
	if !isWhitespace(l.next()) {
		l.emit(
			ERROR,
			fmt.Sprintf("%s must be followed by a space and then '%s'", l.prevToken, kind),
		)
	}
}

// lexClass lexes class declarations and returns lexField state
func lexClass(l *lexer) lexStateFn {
	w := l.readWord()
	if w == "" {
		return nil
	}
	if l.isAccessModifier() {
		l.enforceWhitespace(CLASS)
		w = l.readWord()
	}
	if w != TOKEN_CLASS {
		l.emit(
			ERROR,
			"missing class declaration",
			"class declarations can start with an access modifier (public, private or protected) followed by the 'class' keyword",
		)
	} else {
		l.emit(CLASS)
	}
	l.enforceWhitespace(IDENTIFIER)
	className := l.readWord()
	l.emit(IDENTIFIER)
	if className == "" {
		l.emit(
			ERROR,
			"missing class name",
			"an identifier must follow the 'class' keyword",
		)
	}
	r := l.read()
	if r != TOKEN_OBRACE {
		l.emit(
			ERROR,
			fmt.Sprintf("expected '%c' found '%s'", TOKEN_OBRACE, l.readToken()),
			"multiple identifiers are not supported",
		)
		// Try to recover by reading until the opening brace
		l.readUntil(TOKEN_OBRACE)
		l.emit(
			ERROR,
			"missing opening brace for class body",
		)
	}
	l.emit(OBRACE)
	return lexField
}

// lexField lexes fields inside a class
// returns lexMethod
func lexField(l *lexer) lexStateFn {
	if l.read() == TOKEN_CBRACE {
		l.emit(CBRACE)
		return lexClass
	}
	l.readWord()
	if l.isAccessModifier() {
		l.enforceWhitespace(KEYWORD)
		l.readWord()
	}
	if l.isFieldModifier() {
		l.enforceWhitespace(KEYWORD)
		l.readType()
	}
	l.enforceWhitespace(KEYWORD)
	l.readToken()
	l.emit(IDENTIFIER)

	// Read access modifier, optional static or final modifier, read type, read identifier ...
	// '(' indicates method
	// ';' indicates declaration
	// '=' indicates declaration with initialization

	r := l.read()
	switch r {
	case TOKEN_OPAREN:
		l.emit(OPAREN)
		return lexMethod
	case TOKEN_SEMICOLON:
		l.emit(SEMICOLON)
		// TODO: handle field declaration
	case '=':
		l.emit(ASSIGN)
		// TODO: handle field initialization
	default:
		l.emit(ERROR, "expected '(', ';', or '=' after identifier")
	}
	return lexField
}

// lexMethod lexes parameters inside method parentheses and returns lexMethodBody
func lexMethod(l *lexer) lexStateFn {
	for {
		r := l.read()
		if r == TOKEN_COMMA {
			l.emit(COMMA)
		}
		if !l.readType() {
			l.emit(
				ERROR,
				"invalid parameter type",
				"expected a valid type in parameter list",
			)
		}
		l.enforceWhitespace(PARAMETER)
		l.readToken()
		l.emit(PARAMETER)
		r = l.read()
		if r == TOKEN_COMMA {
			continue
		}
		if r != TOKEN_CPAREN {
			l.emit(
				ERROR,
				"missing closing parenthesis in parameter list",
				"parameters must be separated by commas",
				"end the list with a closing parenthesis",
			)
		}
		// TODO: This is not robust, its fails if theres no whitespace between type and parameter name
		l.emit(CPAREN)
		l.enforceWhitespace(OBRACE)
		r = l.read()
		if r != TOKEN_OBRACE {
			l.emit(
				ERROR,
				"missing opening brace for method body",
			)
		}
		l.emit(OBRACE)
		return lexMethodBody
	}
}

func lexMethodBody(l *lexer) lexStateFn {
	for {
		switch r := l.read(); r {
		case TOKEN_CBRACE:
			l.emit(CBRACE)
			return lexField
		case TOKEN_OPAREN:
			l.emit(OPAREN)
			return lexMethodArguments
		case '.':
			l.emit(DOT)
		case TOKEN_SEMICOLON:
			l.emit(SEMICOLON)
		case '+':
			l.emit(PLUS)
		case '-':
			l.emit(MINUS)
		case '*':
			l.emit(MULTIPLY)
		case '/':
			l.emit(DIVIDE)
		case '%':
			l.emit(PERCENT)
		case '=':
			l.emit(ASSIGN)
		case '!':
			l.emit(NOT)
		case '<':
			l.emit(LT)
		case '>':
			l.emit(GT)
		default:
			if unicode.IsLetter(r) {
				l.readWhile(unicode.IsLetter)
				l.emit(IDENTIFIER)
			}
			if isNumber(r) {
				l.readWhile(isNumber)
				l.emit(LITERAL)
			}
		}
	}
}

func (l *lexer) isAccessModifier() bool {
	switch l.currToken() {
	case "public":
		l.emit(PUBLIC)
	case "private":
		l.emit(PRIVATE)
	case "protected":
		l.emit(PROTECTED)
	}
	return l.runesIsEmpty()
}

func (l *lexer) isFieldModifier() bool {
	switch l.currToken() {
	case "static":
		l.emit(STATIC)
	case "final":
		l.emit(FINAL)
	}
	return l.runesIsEmpty()
}

func (l *lexer) readStringLiteral() {
	l.nextUntil(TOKEN_QUOTE)
	l.read() // consume closing quote
}

func lexMethodArguments(l *lexer) lexStateFn {
	for {
		switch r := l.read(); r {
		case TOKEN_QUOTE:
			l.readStringLiteral()
		case TOKEN_COMMA:
			// TODO: improve if statement, to check runes buffer for unexpected commas
			// consider storing prevtoken as runes buffer
			if l.prevToken != "" && l.prevToken != "," {
				l.emit(COMMA)
				l.enforceWhitespace(ARGUMENT)
			} else {
				l.emit(ERROR, "unexpected comma in argument list")
			}
		case TOKEN_CPAREN:
			l.emit(CPAREN)
			return lexMethodBody
		default:
			l.readToken()
			// l.emit(ERROR, fmt.Sprintf("unexpected token '%s' in argument list", l.readToken()))
		}
		l.emit(ARGUMENT)
	}
}

// readWhile reads runes while the condition is true
// Only call this if the invalid rune should be available for the next read
func (l *lexer) readWhile(cond func(rune) bool) {
	for {
		r := l.next()
		if cond(r) {
			l.runes = append(l.runes, r)
		} else {
			l.backup()
			break
		}
	}
}

// nextUntil adds runes until the delimiter is found.
// reads everything except the delimiter
func (l *lexer) nextUntil(delim rune) {
	l.until(l.next, delim)
}

// readUntil adds runes until the delimiter is found.
// ignores spaces and quotes.
// reads everything except the delimiter
func (l *lexer) readUntil(delim rune) {
	l.until(l.read, delim)
}

func (l *lexer) until(fn func() rune, delim rune) {
	for r := fn(); r != delim; r = fn() {
		l.runes = append(l.runes, r)
	}
	l.backup()
}

// readNumber reads a sequence of digits
// func (l *lexer) readNumber() {
// 	l.readWhile(unicode.IsDigit)
// }

// TODO: Write a more robust and strict switch statement
func (l *lexer) readType() bool {
	l.readWhile(func(r rune) bool {
		return unicode.IsLetter(r) || r == '[' || r == ']' || r == '<' || r == '>' || r == '.'
	})
	return l.isType()
}

// readToken reads an alphanumeric token (letters, digits, underscores)
func (l *lexer) readToken() string {
	l.readWhile(func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
	})
	return l.currToken()
}

func (l *lexer) readWord() string {
	l.readWhile(unicode.IsLetter)
	return l.currToken()
}

// currToken returns the current token as a string
func (l *lexer) currToken() string {
	return string(l.runes)
}

// contains returns true if the current token contains any of the provided runes
// func (l *lexer) contains(runes ...rune) bool {
// 	for _, r := range l.runes {
// 		if slices.Contains(runes, r) {
// 			return true
// 		}
// 	}
// 	return false
// }

// isTokenGenericType returns true if the current token ends with '>'
func (l *lexer) isGeneric() bool {
	return l.runes[len(l.runes)-1] == '>'
}

func (l *lexer) runesIsEmpty() bool {
	return len(l.runes) == 0
}

// isType emits PRIMITIVE or REFERENCE token if current token is a type
// Strings with [ or ] are handled as standard types by removing the brackets before checking
// TODO: Add robust handling for generics and arrays, and error reporting
func (l *lexer) isType() bool {
	if l.isGeneric() {
		l.emit(NOT_SUPPORTED)
		return l.runesIsEmpty()
	}
	// Handle array types by removing brackets
	switch strings.Split(l.currToken(), "[")[0] {
	case "void":
		l.emit(VOID)
	case "boolean", "Boolean":
		l.emit(BOOLEAN)
	case "int", "Integer":
		l.emit(INT)
	case "float", "Float":
		l.emit(FLOAT)
	case "double", "Double":
		l.emit(DOUBLE)
	case "char", "Character":
		l.emit(CHAR)
	case "String":
		l.emit(STRING)
	default:
		l.emit(NOT_SUPPORTED)
	}
	return l.runesIsEmpty()
}

// read returns the next non-whitespace rune.
// whitespace and quotes are not added to runes buffer
func (l *lexer) read() rune {
	for {
		r := l.next()
		if isWhitespace(r) {
			continue
		}
		if r != TOKEN_QUOTE {
			l.runes = append(l.runes, r)
		}
		return r
	}
}

// next returns the next rune from the reader.
// no filtering of any kind is done here
func (l *lexer) next() rune {
	r, _, e := l.reader.ReadRune()
	if e != nil {
		if e == io.EOF {
			l.emit(EOF)
		} else {
			panic(e)
		}
	}
	if r == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
	return r
}

func (l *lexer) nextToken() *token {
	if !l.running {
		l.running = true
		go func() {
			defer l.cleanup()
			// Start the state machine
			for l.state != nil {
				l.state = l.state(l)
			}
			l.emit(EOF)
		}()
	}
	return <-l.tokens
}

// func (l *lexer) peek() rune {
// 	r := l.next()
// 	l.backup()
// 	return r
// }

func (l *lexer) backup() {
	if err := l.reader.UnreadRune(); err != nil {
		l.emit(CRITICAL, "unable to backup lexer reader")
	}
	// Handles backing up when rune is a newline
	if l.column == 1 {
		l.line--
	} else {
		l.column--
	}
}

func (l *lexer) emit(kind tokenKind, msgs ...string) {
	t := &token{l.pos(), l.currToken(), kind}
	if len(msgs) > 0 {
		t.value += ": " + strings.Join(msgs, "\n\t- ")
	}
	l.prevToken = t.value
	l.runes = nil
	l.tokens <- t
}

func (l *lexer) pos() *pos {
	return &pos{
		line:  l.line,
		start: l.column - len(l.runes),
		end:   l.column - 1,
	}
}
