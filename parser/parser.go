package parser

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
)

// TOOD: Add fault tolerance, so that parser can continue parsing after an error
func New(path string, language string) (*Parser, error) {
	if lexer, err := newLexer(path); err != nil {
		return nil, err
	} else {
		execFile := &file{path: path}
		return &Parser{
			lexer:  lexer,
			Target: language,
			ast:    &AST{files: []*file{execFile}},
			curr: curr{
				file: execFile,
			},
			peekToken:  lexer.nextToken(),
			state:      parseClass,
			errorLimit: 5,
		}, nil
	}
}

// Parse parses the tokens and returns the AST or an error if parsing fails.
func (p *Parser) Parse() (*AST, []error) {
	if !p.running {
		p.running = true
		for p.state != nil && p.peekToken.kind != EOF && len(p.errors) < p.errorLimit {
			p.state = p.state(p)
		}
		p.running = false
	}
	return p.ast, p.errors
}

func (t *token) node() node {
	return node{t.value, t.pos}
}

func funcCaller(skip int) (string, string, int, bool) {
	pc, file, line, ok := runtime.Caller(skip)
	fn := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	fnName := fn[len(fn)-1]
	if fnName == "expect" || fnName == "expectNext" {
		return funcCaller(skip + 2)
	}
	return fnName, filepath.Base(file), line, ok
}

// nextToken advances the parser to the next token.
// TODO: Load data to a file field and add it to the ast
func (p *Parser) nextToken() {
	p.prevToken = p.curr.token
	p.curr.token = p.peekToken

	// Block advancement past EOF
	if p.peekToken.kind == EOF {
		funcName, file, line, ok := funcCaller(2)
		fmt.Println("nextToken was called after reaching end of file")
		if ok {
			fmt.Printf("caller: %s (%s:%d)\n\n", funcName, filepath.Base(file), line)
		}
		return
	}

	t := p.lexer.nextToken()
	if t == nil {
		panic("lexer returned nil token. This can happen if token channel is closed")
	}
	switch t.kind {
	case NOT_SUPPORTED:
		p.errorf("token: %s, is not supported", t.value)
	case CRITICAL:
		// TODO: Add fatal error handling
	case ERROR, WARNING:
		p.errorf(fmt.Sprintf("%s: %s", t.kind, t.value))
	case INFO:
		// TODO: Handle info logs
	default:
		p.peekToken = t
	}
}

func (p *Parser) expect(kind ...tokenKind) bool {
	if slices.Contains(kind, p.curr.token.kind) {
		return true
	}
	p.errorf("expected one of %v, got %s", kind, p.curr.token.kind)
	return false
}

func (p *Parser) expectNext(kind ...tokenKind) bool {
	p.nextToken()
	return p.expect(kind...)
}

func (k tokenKind) isModifier() bool {
	switch k {
	case PUBLIC, PRIVATE, PROTECTED, STATIC, FINAL:
		return true
	default:
		return false
	}
}

// parseModifiers parses visibility and other modifiers for classes, methods, and fields.
func (p *Parser) parseModifiers() (mods modifiers, isFinal bool) {
	mods = modifiers{
		visibility: PACKAGE,
	}
	for p.peekToken.kind.isModifier() {
		switch p.nextToken(); p.curr.token.kind {
		case PUBLIC, PRIVATE, PROTECTED:
			if mods.isStatic || isFinal {
				p.errorf("Visibility modifier must be declared before static and final")
			}
			if mods.visibility != PACKAGE {
				p.errorf("Multiple visibility modifiers declared")
			}
			mods.visibility = p.curr.token.kind
		case STATIC:
			mods.isStatic = true
		case FINAL:
			if p.peekToken.kind == STATIC {
				p.errorf("Final can't be declared before static")
			}
			isFinal = true
		}
	}
	return mods, isFinal
}

// TODO: parseClass only operates per token instance...
// debate on what is best practice. One state iteration should live longer than one token
// new structure could be
// eat mods with switch statement until class keyword
// then expect identifier for class name
// then expect opening brace
// then enter declaration parsing state
// This state should arguably not handle closing brace...
// could implement a buffer to hold tokens until closing brace is found
// but not doing anything while in that state seems wasteful
func parseClass(p *Parser) parseStateFn {
	mods, _ := p.parseModifiers()
	p.expectNext(CLASS)
	p.expectNext(IDENTIFIER)
	p.expectNext(OBRACE)
	p.curr.class = &class{
		modifiers: mods,
		node:      p.prevToken.node(),
	}
	return parseDeclaration
}

func parseDeclaration(p *Parser) parseStateFn {
	if p.curr.token.kind == CBRACE {
		p.nextToken()
		p.addClass(p.curr.class)
		return parseClass
	}
	mods, isFinal := p.parseModifiers()
	p.expectNext(INT, FLOAT, STRING, BOOLEAN, VOID)
	p.expectNext(IDENTIFIER)
	p.decl = &decl{
		modifiers: mods,
		isFinal:   isFinal,
		node:      p.curr.token.node(),
		kind:      p.prevToken.kind,
	}
	switch p.nextToken(); p.curr.token.kind {
	case OPAREN:
		p.nextToken()
		return parseParams
	case ASSIGN:
		return parseField
	case SEMICOLON:
		p.class.fields = append(p.class.fields, &field{decl: p.decl})
	case CBRACE:
		// if p.peekToken.kind == EOF {
		// 	p.errorf("class: %s is missing a closing bracket", p.getClass().name)
		// 	return nil
		// }
		return parseClass
	default:
		p.errorf("unexpected token (declaration): %s", p.curr.token.kind)
	}
	return parseDeclaration
}

func parseArguments(p *Parser) parseStateFn {
	m := &method{decl: &decl{node: p.prevToken.node()}}
	for p.curr.token.kind != CPAREN {
		p.expectNext(ARGUMENT)
		m.parameters = append(m.parameters, &decl{
			kind: p.curr.token.kind,
			node: p.curr.token.node(),
		})
		p.expectNext(COMMA, CPAREN)
	}
	p.object.methods = append(p.object.methods, m)
	p.expectNext(SEMICOLON)
	return parseMethodBody
}

func parseParams(p *Parser) parseStateFn {
	var params []*decl
	for p.curr.token.kind != CPAREN {
		p.expect(INT, FLOAT, STRING, BOOLEAN)
		p.expectNext(PARAMETER)
		params = append(params, &decl{
			kind: p.prevToken.kind,
			node: p.curr.token.node(),
		})
		p.expectNext(COMMA, CPAREN)
	}
	p.expectNext(OBRACE)
	p.method = &method{
		decl:       p.decl,
		parameters: params,
	}
	return parseMethodBody
}

// parseMethodBody parses lexer tokens until it reaches the end of the method
func parseMethodBody(p *Parser) parseStateFn {
	for p.curr.token.kind != CBRACE {
		switch p.nextToken(); p.curr.token.kind {
		case IDENTIFIER:
			return parseExpression
		case ASSIGN, LITERAL:
			// initialization
			// should handle literal tokens
		case OPAREN:
			return parseArguments
		}
	}
	p.addMethod()
	return parseDeclaration
}

func parseExpression(p *Parser) parseStateFn {
	// Todo, only works for the first object, next will overrite it selft and its never allowing nesting past 2 iterations.
	if p.object == nil {
		p.object = &object{
			node: p.curr.token.node(),
		}
	}
	if p.peekToken.kind == DOT {
		p.nextToken() // consume DOT
		p.expectNext(IDENTIFIER)
		obj := &object{
			node: p.curr.token.node(),
		}
		p.object.fields = append(p.object.fields, obj)
		p.method.body.objects = append(p.method.body.objects, p.object)
		p.object = obj
	} else if p.peekToken.kind != OPAREN {
		p.object = nil // important to allow new standalone objects to be created
		// not doing this results in an infinite nested object, even though do not relate to each other in the syntax.
	}
	return parseMethodBody
}

func parseField(p *Parser) parseStateFn {
	return nil
}

func (p *Parser) addClass(c *class) {
	p.file.classes = append(p.file.classes, c)
}

func (p *Parser) addMethod() {
	p.class.methods = append(p.class.methods, p.method)
}

func (p *Parser) addField(initValue string) {
	p.class.fields = append(p.class.fields, &field{decl: p.decl, initValue: initValue})
}

func (p *Parser) errorf(format string, args ...any) {
	caller, _, _, ok := funcCaller(2)
	format = fmt.Sprintf("%s: %s", p.curr.token.pos, format)
	if ok {
		format = fmt.Sprintf("(%s) %s", caller, format)
	}
	p.errors = append(p.errors, fmt.Errorf(format, args...))
}

func (p *Parser) error(err string) {
	p.errors = append(p.errors, errors.New(err))
}
