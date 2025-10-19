package spec

import (
	"os"
)

type Runner interface {
	Run(out string) *os.File
}
