package cmdUtils

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
)

// HandleError display the error and stop the process if error exist
func HandleError(err error) {
	if err != nil {
		fmt.Println(aurora.Red(err))
		os.Exit(0)
	}
}
