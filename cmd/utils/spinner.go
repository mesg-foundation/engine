package cmdUtils

import (
	"time"

	"github.com/briandowns/spinner"
)

// SpinnerOptions is a struct that contains all details for the spinner
type SpinnerOptions struct {
	Text  string
	Color string
}

// StartSpinner create a new spinner for the terminal
func StartSpinner(opts SpinnerOptions) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Start()
	if opts.Color != "" {
		s.Color(opts.Color)
	}
	if opts.Text != "" {
		s.Suffix = " " + opts.Text
	}
	return s
}
