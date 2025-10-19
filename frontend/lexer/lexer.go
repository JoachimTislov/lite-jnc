package lexer

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/JoachimTislov/lite-jnc/spec"
)

type Lexer struct {
	Source source
	Lines  []Line
	CLine  int // current line
	CToken int // current token
}

type source struct {
	Path   string
	Length int
}

type Line struct {
	Indent        int
	TrailingSpace int
	Length        int
	Tokens        []spec.Token
}

func New(src string) (*Lexer, error) {
	bytes, err := os.ReadFile(src)
	if err != nil {
		return nil, err
	}
	trimmed := strings.TrimSpace(string(bytes))
	lines := strings.Split(trimmed, "\n")
	l := make([]Line, len(lines))
	for i, line := range lines {
		LenLine := len(line)
		symbols := strings.Fields(line)
		tokens := make([]spec.Token, len(symbols))
		for i, field := range symbols {
			startIndex := strings.Index(line, field)
			tokens[i] = spec.Token{
				Pos: spec.Pos{
					Line:  i + 1,
					Start: startIndex + 1,
					End:   startIndex + len(field),
				},
				Value: field,
			}
		}
		l[i] = Line{
			Indent:        LenLine - len(strings.TrimLeft(line, " ")),
			TrailingSpace: LenLine - len(strings.TrimRight(line, " ")),
			Length:        LenLine,
			Tokens:        tokens,
		}
	}
	return &Lexer{
		Source: source{Path: src, Length: len(lines)},
		Lines:  l,
	}, nil
}

func (l *Lexer) NextToken() (*Token, error) {
	if l.Source.Length < l.CLine {
		return nil, EOF
	}
	if len(l.Lines[l.CLine].Tokens) <= l.CToken {
		l.CToken = 0
		l.CLine++
	}
	token := l.Lines[l.CLine].Tokens[l.CToken]
	l.CToken++
	return &token, nil
}

func (l *Lexer) Json() (string, error) {
	data, err := json.MarshalIndent(*l, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
