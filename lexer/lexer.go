package lexer

import (
	"encoding/json"
	"os"
	"strings"
)

type source struct {
	Path   string
	Length int
}

type lexer struct {
	Source source
	Lines  []Line
	CLine  int // current line
	CToken int // current token
}

type Line struct {
	Indent   int
	Trailing int
	Length   int
	Tokens   []string
}

func New(src string) (*lexer, error) {
	bytes, err := os.ReadFile(src)
	if err != nil {
		return nil, err
	}
	trimmed := strings.TrimSpace(string(bytes))
	lines := strings.Split(trimmed, "\n")
	l := make([]Line, len(lines))
	for i, line := range lines {
		LenLine := len(line)
		l[i] = Line{
			Indent:   LenLine - len(strings.TrimLeft(line, " ")),
			Trailing: LenLine - len(strings.TrimRight(line, " ")),
			Length:   LenLine,
			Tokens:   strings.Fields(line),
		}
	}
	return &lexer{
		Source: source{Path: src, Length: len(lines)},
		Lines:  l,
	}, nil
}

func (l *lexer) NextToken() string {
	if l.Source.Length < l.CLine {
		return ""
	}
	if l.Lines[l.CLine].Length < l.CToken {
		l.CToken = 0
		l.CLine++
	}
	token := l.Lines[l.CLine].Tokens[l.CToken]
	l.CToken++
	return token
}

func (l *lexer) Json() (string, error) {
	data, err := json.MarshalIndent(*l, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
