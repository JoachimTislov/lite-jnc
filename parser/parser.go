package parser

import (
	"errors"
	"fmt"
)

type Node interface {
	Name() string
	Position() pos
}

type Statement interface {
	Node
	Execute()
}

type Expression interface {
	Node
	Evaluate()
}

type node struct {
	name string
	pos  *pos
}

func (n *node) Name() string   { return n.name }
func (n *node) Position() *pos { return n.pos }

type decl struct {
	*node
	value string
	kind  tokenKind
}

// type variableDecl decl
type parameter decl

type modifiers []*tokenKind

type field struct {
	decl
	modifiers
	initVal string
}

type method struct {
	decl
	modifiers
	body
	parameters []*parameter
	returnType tokenKind
}

type body struct {
	statements  []*Statement
	Expressions []*Expression
}

type class struct {
	*node
	visibility tokenKind
	fields     []*field
	methods    []*method
}

type pkg struct {
	Name    string
	Classes []*class
}

type importDecl struct {
	pkgName string
	alias   string
}

type file struct {
	path string
	// pkg     pkg
	imports []*importDecl
	classes []*class
}

type AST struct {
	packages []*pkg
	files    []*file
}

func (a *AST) String() string {
	for _, f := range a.files {
		for _, c := range f.classes {
			fmt.Print(c)
			for _, m := range c.methods {
				fmt.Print("\n", m)
				for _, p := range m.parameters {
					fmt.Print("\n", p)
				}
			}
		}
	}
	return ""
}

func (m *method) String() string {
	return fmt.Sprintf("method: %s, return type: %s, pos: %v, params: %d", m.name, m.returnType, m.pos, len(m.parameters))
}

func (c *class) String() string {
	return fmt.Sprintf("class: %s, pos: %v, methods: %d, fields: %d", c.name, c.pos, len(c.methods), len(c.fields))
}

func (p *parameter) String() string {
	return fmt.Sprintf("param: %s, type: %s, pos: %v, value: %s", p.name, p.kind, p.pos, p.value)
}

func (f *field) String() string {
	return fmt.Sprintf("field: %s, type: %s, pos: %v, modifiers: %d, initVal: %s", f.name, f.kind, f.pos, len(f.modifiers), f.initVal)
}

type Parser struct {
	Target string
	*lexer
	ast      *AST
	currFile *file
	prev     *token
	curr     *token
	peek     *token
	errors   []error
}

func New(path string, language string) (*Parser, error) {
	lexer, err := newLexer(path)
	if err != nil {
		return nil, err
	}
	f := &file{path: path}
	p := &Parser{
		lexer:    lexer,
		Target:   language,
		currFile: f,
		ast:      &AST{files: []*file{f}},
	}

	// prime current and peek tokens
	p.nextToken()
	p.nextToken()

	return p, nil
}

// nextToken advances the parser to the next token.
func (p *Parser) nextToken() {
	p.prev = p.curr
	p.curr = p.peek
	p.peek = p.lexer.nextToken()
}

func (p *Parser) expect(kind tokenKind) bool {
	if p.nextToken(); p.curr.kind == kind {
		return true
	}
	p.errorf("expected %s, got %s", kind, p.curr.kind)
	return false
}

// Parse parses the tokens and returns the AST or an error if parsing fails.
func (p *Parser) Parse() (*AST, error) {
	for p.curr.kind != EOF {
		p.nextToken()
		switch p.curr.kind {
		case PACKAGE:
		// TODO: p.parsePackage()
		case IMPORT:
		// TODO: p.parseImport()
		case CLASS:
			class, err := p.parseClass()
			if err != nil {
				p.errorf("failed to parse class at %v: %v", class.pos, err)
			}
			p.currFile.classes = append(p.currFile.classes, class)
		case ERROR:
			p.errorf("lexical error(s) at %v: \n\t%s", p.curr.pos, p.curr.value)
		}
	}

	if len(p.errors) > 0 {
		return p.ast, fmt.Errorf("parsing errors:\n %v", errors.Join(p.errors...))
	}

	return p.ast, nil
}

func (t tokenKind) isAccessModifier() bool {
	return t == PUBLIC || t == PRIVATE || t == PROTECTED
}

// parseClass, enters at 'class' token, expects class name next and modifer to be previous
func (p *Parser) parseClass() (*class, error) {
	class := &class{
		node: &node{
			name: p.peek.value,
			pos:  p.curr.pos,
		},
	}
	if p.prev != nil {
		if !p.prev.kind.isAccessModifier() {
			return class, fmt.Errorf("expected modifier, got %s", p.prev.kind)
		}
		class.visibility = p.prev.kind
	}
	p.expect(IDENTIFIER)
	p.expect(OBRACE)

	closeBraceCount := 1
	var mods modifiers
	var d decl
	for closeBraceCount > 0 {
		p.nextToken()
		switch p.curr.kind {
		case PUBLIC, PRIVATE, PROTECTED, STATIC, FINAL:
			mods = append(mods, &p.curr.kind)
		case INT, FLOAT, STRING, BOOLEAN, VOID:
			d = decl{
				node: &node{
					name: p.peek.value,
					pos:  p.peek.pos,
				},
				kind: p.curr.kind,
			}
		case OPAREN:
			method := &method{
				decl:       d,
				modifiers:  mods,
				parameters: p.parseParameters(),
				returnType: d.kind,
			}
			// reset mods and decl
			// not good practice TODO: remove
			mods = nil
			d = decl{}
			// TODO: parse method body
			// skip to closing brace for now
			for p.curr.kind != CBRACE {
				p.nextToken()
			}
			class.methods = append(class.methods, method)
		case ASSIGN:
			f := &field{
				decl:      d,
				modifiers: mods,
				initVal:   p.peek.value,
			}
			class.fields = append(class.fields, f)
			p.nextToken() // move to init value
		case SEMICOLON:
			class.fields = append(class.fields, &field{
				decl:      d,
				modifiers: mods,
			})
		case CBRACE:
			closeBraceCount--
		case OBRACE:
			closeBraceCount++
		}
		if closeBraceCount == 0 {
			break
		}
	}
	return class, nil
}

func (p *Parser) parseParameters() []*parameter {
	parameters := []*parameter{}
	for p.peek.kind != CPAREN {
		p.nextToken()
		if p.curr.kind == PARAMETER {
			param := &parameter{
				node: &node{
					name: p.curr.value,
					pos:  p.curr.pos,
				},
				value: p.prev.value,
				kind:  p.prev.kind,
			}
			parameters = append(parameters, param)
		}
	}
	return parameters
}

func (p *Parser) parseMethod() (*method, error) {
	return &method{}, nil
}

func (p *Parser) errorf(format string, args ...any) {
	p.addError(fmt.Errorf(format, args...))
}

// func (p *Parser) error(err string) {
// 	p.addError(errors.New(err))
// }

func (p *Parser) addError(err error) {
	p.errors = append(p.errors, err)
}
