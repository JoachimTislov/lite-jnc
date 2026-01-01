package parser

import "fmt"

func (t token) String() string {
	if t.kind == ERROR {
		return fmt.Sprintf("ERROR(%v): \n\t%s", t.pos, t.value)
	}
	return fmt.Sprintf("\n%s, %s, %s", t.kind, t.value, t.pos)
}

func (p pos) String() string {
	rang := fmt.Sprintf("%d-%d", p.start, p.end)
	if p.start >= p.end {
		rang = fmt.Sprintf("%d", p.end)
	}
	if p.end == 0 {
		return fmt.Sprintf("%d", p.line)
	}
	return fmt.Sprintf("%d:%s", p.line, rang)
}
