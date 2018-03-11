package cmd

import (
	"time"

	"github.com/briandowns/spinner"
)

type spinnerOptions struct {
	text  string
	color string
}

func startSpinner(opts spinnerOptions) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Start()
	if opts.color != "" {
		s.Color(opts.color)
	}
	if opts.text != "" {
		s.Suffix = " " + opts.text
	}
	return s
}
