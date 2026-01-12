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
		exec := &file{path: path}
		return &Parser{
			lexer:  lexer,
			Target: language,
			ast:    &AST{files: []*file{exec}},
			curr: curr{
				file: exec,
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
func (p *Parser) nextToken() {
	p.prevToken = p.token
	p.token = p.peekToken

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

// expect recursively calls itself and nextToken until it hits one of the kinds provided at the initial call
func (p *Parser) expect(kind ...tokenKind) bool {
	if slices.Contains(kind, p.token.kind) {
		return true
	}
	if len(kind) > 1 {
		p.errorf("expected one of %v, got %s", kind, p.token.kind)
	} else {
		p.errorf("expected %v, got %s", kind[0], p.token.kind)
	}
	if p.token.kind == EOF {
		return false
	}
	return p.expectNext(kind...)
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
		switch p.nextToken(); p.token.kind {
		case PUBLIC, PRIVATE, PROTECTED:
			if mods.isStatic || isFinal {
				p.errorf("Visibility modifier must be declared before static and final")
			}
			if mods.visibility != PACKAGE {
				p.errorf("Multiple visibility modifiers declared")
			}
			mods.visibility = p.token.kind
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

func parseClass(p *Parser) parseStateFn {
	mods, _ := p.parseModifiers()
	p.expectNext(CLASS)
	p.expectNext(IDENTIFIER)
	p.expectNext(OBRACE)
	p.class = &class{
		modifiers: mods,
		node:      p.prevToken.node(),
	}
	return parseDeclaration
}

func parseDeclaration(p *Parser) parseStateFn {
	if p.token.kind == CBRACE {
		p.nextToken()
		p.addClass(p.class)
		return parseClass
	}
	mods, isFinal := p.parseModifiers()
	p.expectNext(INT, FLOAT, STRING, BOOLEAN, VOID)
	p.expectNext(IDENTIFIER)
	p.decl = &decl{
		modifiers: mods,
		isFinal:   isFinal,
		node:      p.token.node(),
		kind:      p.prevToken.kind,
	}
	switch p.nextToken(); p.token.kind {
	case OPAREN:
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
		p.errorf("unexpected token (declaration): %s", p.token.kind)
	}
	return parseDeclaration
}

func parseFunc(p *Parser) parseStateFn {
	fn := &fn{reference: p.reference}
	// p.reference = nil
	for p.token.kind != CPAREN {
		p.expectNext(ARGUMENT)
		fn.args = append(fn.args, p.token.node())
		p.expectNext(COMMA, CPAREN)
	}
	p.expectNext(SEMICOLON)
	p.method.expressions = append(p.method.expressions, fn)
	return parseMethodBody
}

func parseParams(p *Parser) parseStateFn {
	var params []*parameter
	for p.token.kind != CPAREN {
		p.expectNext(INT, FLOAT, STRING, BOOLEAN)
		p.expectNext(PARAMETER)
		params = append(params, &parameter{
			kind: p.prevToken.node(),
			name: p.token.node(),
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
	for p.token.kind != CBRACE {
		switch p.nextToken(); p.token.kind {
		case IDENTIFIER:
			return parseExpression
		case ASSIGN, LITERAL:
			return parseStatement
			// initialization
			// should handle literal tokens
		case SEMICOLON:
			p.addReference()
		case OPAREN:
			return parseFunc
		}
	}
	p.addMethod()
	return parseDeclaration
}

func parseExpression(p *Parser) parseStateFn {
	p.createReference()
	if p.peekToken.kind == DOT {
		p.nextToken() // consume DOT
		p.expectNext(IDENTIFIER)
		p.createReference()
	}
	return parseMethodBody
}

func parseStatement(p *Parser) parseStateFn {
	return nil
}

func parseField(p *Parser) parseStateFn {
	return nil
}

func (p *Parser) addReference() {
	p.method.references = append(p.method.references, p.reference)
	p.reference = nil
}

func (p *Parser) createReference() {
	if p.reference == nil {
		p.reference = &reference{
			node: p.token.node(),
		}
	} else {
		p.reference = &reference{
			node:   p.token.node(),
			parent: p.reference,
		}
	}
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
	format = fmt.Sprintf("%s: %s", p.token.pos, format)
	if ok {
		format = fmt.Sprintf("(%s) %s", caller, format)
	}
	p.errors = append(p.errors, fmt.Errorf(format, args...))
}

func (p *Parser) error(err string) {
	p.errors = append(p.errors, errors.New(err))
}
