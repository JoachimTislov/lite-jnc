package parser

import "fmt"

func (a *AST) String() string {
	prettyPrintAST(a)
	return ""
}

func (a modifiers) String() string {
	return fmt.Sprintf("Visibility: %s, isStatic: %t", a.visibility, a.isStatic)
}

func (m *method) String() string {
	fmt.Printf("\tMethod: %s", m.name)
	for _, p := range m.parameters {
		fmt.Printf("\n\t  kind: %s\n\t  name: %s\n", p.kind.name, p.name.name)
	}
	return ""
}

func prettyPrintAST(a *AST) {
	for _, f := range a.files {
		fmt.Printf("File: %s\n", f.path)
		for _, c := range f.classes {
			fmt.Printf(" Class: %s (Visibility: %s, Static: %t)\n", c.name, c.visibility, c.isStatic)
			for _, fld := range c.fields {
				fmt.Printf("  Field: %s %s (%s, InitVal: %s)\n", fld.kind, fld.name, fld.modifiers, fld.initValue)
			}
			for _, m := range c.methods {
				fmt.Printf("  Method: return type: %s name: %s (%s)\n", m.kind, m.name, m.modifiers)
				fmt.Printf("    Parameters:\n")
				for _, p := range m.parameters {
					fmt.Printf("\t- kind: %s name: %s\n", p.kind.name, p.name.name)
				}
				fmt.Printf("   Body:\n")
				for _, ref := range m.body.references {
					fmt.Printf("    Reference: %s\n", ref.name)
				}
			}
		}
	}
}
