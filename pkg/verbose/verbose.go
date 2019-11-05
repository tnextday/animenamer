package verbose

import (
	"fmt"
	"os"
)

var V = false

func Print(format string, a ...interface{}) {
	if V {
		fmt.Fprintf(os.Stderr, format, a...)
	}
}
