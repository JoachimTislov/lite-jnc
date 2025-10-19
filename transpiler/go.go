package transpiler

import (
	"os"
)

type System struct{}

func (t *transpiler) Go(name string) error {
	return os.WriteFile(name, make([]byte, 0), os.ModePerm)
}
