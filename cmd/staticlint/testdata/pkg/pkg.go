package pkg

import (
	"os"
)

func noErrInPkgFunc() {
	os.Exit(1)
}
