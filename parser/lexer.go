package parser

import (
	"bufio"
	"io"
	"os"
	"strings"
	"unicode"
)

type stateFn func(*lexer) stateFn

type lexer struct {
	reader  *bufio.Reader
	line    int
	column  int
	runes   []rune
	error   string
	state   stateFn
	tokens  chan *token
	cleanup func()
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
		state:  lexText,
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

// lex switches to the appropriate state function based on the rune
func lexText(l *lexer) stateFn {
	for {
		r := l.read()
		switch r {
		case TOKEN_QUOTE:
			l.readStringLiteral()
		case '(':
			return lexParen
		case '[', '{':
			l.emit(OPUNCTUATION)
		case ']', '}':
			l.emit(CPUNCTUATION)
		case ':', ',', '.', ';':
			l.emit(DELIMITER)
		case '+', '-', '*', '/', '%', '=', '!', '<', '>':
			l.emit(OPERATOR)
		default:
			if unicode.IsLetter(r) {
				return lexToken
			}
			l.emit(ILLEGAL)
		}
	}
}

func (l *lexer) readStringLiteral() {
	l.readEmit(TOKEN_QUOTE, LITERAL)
}

func (l *lexer) readEmit(r rune, kind tokenKind) {
	l.readUntil(r)
	l.emit(kind)
}

// lexParen lexes a parameter list inside parentheses
func lexParen(l *lexer) stateFn {
	l.emit(OPUNCTUATION)
	for {
		r := l.read()
		switch r {
		case TOKEN_QUOTE:
			l.readStringLiteral()
			continue
		case ',':
			l.emit(DELIMITER)
			continue
		case ')':
			l.emit(CPUNCTUATION)
			return lexText
		}

		l.readType()
		l.emit(TYPE)

		if isWhitespace(l.next()) {
			l.readToken()
			l.emit(PARAMETER)
		} else {
			l.emit(ERROR, "expected whitespace after type in parameter list", "parameters must be separated by commas")
			l.backup()
		}
	}
}

// readWhile reads runes while the condition is true
// Only call this if the invalid rune should be available for the next read
func (l *lexer) readWhile(cond func(rune) bool) string {
	for {
		r := l.next()
		if cond(r) {
			l.runes = append(l.runes, r)
		} else {
			l.backup()
			break
		}
	}
	return string(l.runes)
}

// readUntil adds runes the delimiter is found
// Similar to reader.readString() but does not include the delimiter
func (l *lexer) readUntil(delim rune) {
	for {
		r := l.next()
		if r == delim {
			break
		}
		l.runes = append(l.runes, r)
	}
}

// readNumber reads a sequence of digits
// func (l *lexer) readNumber() string {
// 	return l.readWhile(func(r rune) bool {
// 		return unicode.IsDigit(r)
// 	})
// }

func (l *lexer) readType() string {
	return l.readWhile(func(r rune) bool {
		return unicode.IsLetter(r) || r == '[' || r == ']'
	})
}

// readToken reads an alphanumeric token (letters, digits, underscores)
func (l *lexer) readToken() string {
	return l.readWhile(func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
	})
}

// readWord reads a sequence of letters
func (l *lexer) readWord() string {
	return l.readWhile(func(r rune) bool {
		return unicode.IsLetter(r)
	})
}

// func (l *lexer) eatWhitespace() {
// 	for {
// 		r := l.next()
// 		if !isWhitespace(r) {
// 			l.backup()
// 			break
// 		}
// 	}
// }

// readToken reads a sequence of letters and goes into lexIdentifier state
func lexToken(l *lexer) stateFn {
	switch l.readWord() {
	case "public", "private", "protected", "static":
		l.emit(MODIFIER)
	case "interface", "abstract", "class", "extends", "implements", "new", "return", "if", "else", "for", "while", "break", "continue":
		l.emit(KEYWORD)
	case "void", "int", "float", "double", "char", "boolean", "String":
		l.emit(TYPE)
	default:
		l.emit(IDENTIFIER)
	}
	return lexText
}

// read returns the next non-whitespace rune.
// Does not add whitespace or quotes to the runes buffer
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

func (l *lexer) run() {
	defer l.cleanup()
	for l.state != nil {
		l.state = l.state(l)
	}
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) backup() {
	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}
	// Handles backing up when rune is a newline
	if l.column == 1 {
		l.line--
	} else {
		l.column--
	}
}

func (l *lexer) emit(kind tokenKind, errors ...string) {
	var t *token
	if len(errors) > 0 {
		t = &token{l.pos(), strings.Join(errors, "; "), ERROR}
	} else {
		t = &token{l.pos(), string(l.runes), kind}
		l.runes = nil
	}
	l.tokens <- t
}

func (l *lexer) pos() pos {
	return pos{
		line:  l.line,
		start: l.column - len(l.runes),
		end:   l.column - 1,
	}
}
