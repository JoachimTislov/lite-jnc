package parser

import (
	"bufio"
	"io"
	"os"
	"unicode"
)

type stateFn func(*lexer) stateFn

type lexer struct {
	reader    *bufio.Reader
	line      int
	column    int
	currRunes []rune
	state     stateFn
	tokens    chan *token
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
		r := l.next()
		if isWhitespace(r) {
			continue
		}
		l.currRunes = []rune{r}
		switch r {
		case TOKEN_QUOTE:
			return lexStringLiteral
		case '(':
			return lexParen
		case '{':
			return lexBrace
		case '[':
			// return lexBracket
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

func lexParen(l *lexer) stateFn {
	l.emit(PUNCTUATION)
	return lexText
}

func lexBrace(l *lexer) stateFn {
	l.emit(PUNCTUATION)
	return lexText
}

func lexStringLiteral(l *lexer) stateFn {
	for {
		r := l.next()
		l.currRunes = append(l.currRunes, r)
		if r == TOKEN_QUOTE {
			l.emit(LITERAL)
			return lexText
		}
	}
}

func (l *lexer) readToken() string {
	for {
		r := l.next()
		if unicode.IsLetter(r) {
			l.currRunes = append(l.currRunes, r)
		} else {
			l.backup()
			break
		}
	}
	return string(l.currRunes)
}

// readToken reads a sequence of letters and goes into lexIdentifier state
func lexToken(l *lexer) stateFn {
	switch l.readToken() {
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
		l.column = 0
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

func (l *lexer) backup() {
	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}
	l.column--
}

func (l *lexer) emit(kind tokenKind) {
	l.tokens <- &token{l.pos(), string(l.currRunes), kind}
	l.currRunes = nil
}

func (l *lexer) pos() pos {
	tokenLen := len(l.currRunes)
	start := l.column - tokenLen
	if tokenLen == 1 {
		start = l.column
	}
	return pos{l.line, start, l.column}
}
