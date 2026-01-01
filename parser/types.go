package parser

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

func (n node) Name() string   { return n.name }
func (n node) Position() *pos { return n.pos }

type decl struct {
	node
	// kind is the returnType for methods
	kind    tokenKind
	isFinal bool
	modifiers
}

// object represents an instance of a class
type object struct {
	node
	fields  []*object
	methods []*method
}

type field struct {
	*decl
	initValue string
}

type method struct {
	*decl
	parameters []*decl
	body
}

type body struct {
	objects     []*object
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
	packages []pkg
	files    []*file
}

type parseStateFn func(*Parser) parseStateFn

// curr holds the current parsing context
// per state iteration (e.g., parsing a class, method, field, etc.)
type curr struct {
	*token
	class  *class
	method *method
	object *object
	file   *file
	decl   *decl
}

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
