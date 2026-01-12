package parser

type Node interface {
	Name() string
	Position() *pos
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

func (n node) Name() string   { return n.name }
func (n node) Position() *pos { return n.pos }

type decl struct {
	node
	// kind is the returnType for methods
	kind    tokenKind
	isFinal bool
	modifiers
}

type parameter struct {
	name node
	kind node
}

// reference represents an expression referencing a field, variable or method
// inverse definition, end token is most relevant. Example; System.out.print(...)
// becomes "print <- out <- System"
type reference struct {
	node
	parent *reference
}

type fn struct {
	*reference
	args []node
}

func (fn *fn) Evaluate() {}
func (fn *fn) Position() *pos {
	return fn.reference.pos
}

type field struct {
	*decl
	initValue string
}

type method struct {
	*decl
	parameters []*parameter
	body
}

type body struct {
	references  []*reference
	statements  []Statement
	expressions []Expression
}

type modifiers struct {
	visibility tokenKind
	isStatic   bool
}

type class struct {
	node
	modifiers
	isClosed bool
	fields   []*field
	methods  []*method
}

type pkg struct {
	name    string
	classes []class
}

type importDecl struct {
	pkgName string
	alias   string
}

type file struct {
	path string
	// pkg     pkg
	imports []importDecl
	classes []*class
}

type AST struct {
	packages []*pkg
	files    []*file
}

type parseStateFn func(*Parser) parseStateFn

// curr holds the current parsing context
// per state iteration (e.g., parsing a class, method, field, etc.)
// acts with alot of side effects. Not ideal. TODO: Should reconsider design
// keep in mind that its either these or more code per state function and thereby fewer state functions
type curr struct {
	*token
	class     *class
	method    *method
	reference *reference
	file      *file
	decl      *decl
}

// TODO: Consider adding a token buffer, possibly replace prev, peek with a slice of tokens
type Parser struct {
	Target string
	*lexer
	prevToken *token
	peekToken *token
	curr
	running    bool
	ast        *AST
	state      parseStateFn
	errors     []error
	errorLimit int
}
