package cmdUtils

import (
	"time"

	spinnerPkg "github.com/briandowns/spinner"
)

// SpinnerOptions is a struct that contains all details for the spinner
type SpinnerOptions struct {
	Text  string
	Color string
}

// StartSpinner create a new spinner for the terminal
func StartSpinner(opts SpinnerOptions) (spinner *spinnerPkg.Spinner) {
	spinner = spinnerPkg.New(spinnerPkg.CharSets[11], 100*time.Millisecond)
	spinner.Start()
	if opts.Color != "" {
		spinner.Color(opts.Color)
	}
	if opts.Text != "" {
		spinner.Suffix = " " + opts.Text
	}
	return
}
