package parser

import (
	"errors"
	"fmt"
)

type Node interface {
	Name()
	Position() pos
	Token() token
}

type Declaration interface {
	Node
	Declare()
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
	name  string
	token token
}

func (n *node) Name() string   { return n.name }
func (n *node) Position() *pos { return n.token.pos }
func (n *node) Token() token   { return n.token }

type nodeKind int

// const (
// 	PRIMITIVE nodeKind = iota
// 	REFERENCE
// 	STRUCTURE
// 	INTERFACE
// 	ARRAY
// 	VARIABLE
// )

type MethodDecl struct {
	node
	IsStatic     bool
	IsFinal      bool
	Parameters   []ParameterDecl
	Declarations []Declaration
	Statements   []Statement
	Expressions  []Expression
	ReturnType   nodeKind
}

type decl struct {
	node
	Type nodeKind
}

type VariableDecl decl
type ParameterDecl decl
type FieldDecl struct {
	decl
	IsStatic bool
	IsFinal  bool
}

type accessibility int

// const (
// 	PUBLIC accessibility = iota
// 	PRIVATE
// 	PROTECTED
// )

type classDecl struct {
	node
	Visibility accessibility
	Parameters []ParameterDecl
	Fields     []FieldDecl
	Methods    []MethodDecl
	ReturnType nodeKind
}

type pkg struct {
	Name    string
	Classes []classDecl
}

type importDecl struct {
	pkgName string
	alias   string
}

type file struct {
	path string
	// pkg     pkg
	imports []importDecl
	classes []classDecl
}

type AST struct {
	packages []pkg
	files    []file
}

type Parser struct {
	Target string
	*lexer
	ast    *AST
	prev   *token
	curr   *token
	peek   *token
	errors []error
}

func New(path string, language string) (*Parser, error) {
	lexer, err := newLexer(path)
	if err != nil {
		return nil, err
	}
	p := &Parser{
		lexer:  lexer,
		Target: language,
		ast:    &AST{},
	}

	p.nextToken()
	p.nextToken()

	return p, nil
}

func (p *Parser) nextToken() {
	p.prev = p.curr
	p.curr = p.peek
	p.peek = p.lexer.nextToken()
}

func (p *Parser) Parse() (*AST, error) {
	for p.curr.kind != EOF {

		if p.curr.kind == ERROR {
			p.errorf("lexical error(s) at %v: \n\t%s", p.curr.pos, p.curr.value)
		}

		fmt.Println(p.curr)
		p.nextToken()
	}

	if len(p.errors) > 0 {
		return p.ast, fmt.Errorf("parsing errors:\n %v", errors.Join(p.errors...))
	}

	return p.ast, nil
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
