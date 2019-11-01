package verbose

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var V = false

func init() {
	V = viper.GetBool("verbose")
}
func Print(format string, a ...interface{}) {
	if V {
		fmt.Fprintf(os.Stderr, format, a...)
	}
}
