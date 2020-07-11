package verbose

import (
	"fmt"
	"os"
)

var V = false

func Printf(format string, a ...interface{}) {
	if V {
		fmt.Fprintf(os.Stderr, format, a...)
	}
}
