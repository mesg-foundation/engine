package utils

import (
	"time"

	"github.com/briandowns/spinner"
)

var (
	// spinnerCharset is the default animation.
	spinnerCharset = spinner.CharSets[11]

	// spinnerDuration is the default duration for spinning.
	spinnerDuration = 100 * time.Millisecond
)

// SpinnerOptions contains all details for the spinner
type SpinnerOptions struct {
	Text  string
	Color string
}

// StartSpinner creates new spinner for terminal.
func StartSpinner(opts SpinnerOptions) *spinner.Spinner {
	spinner := spinner.New(spinnerCharset, spinnerDuration)
	spinner.Start()
	if opts.Color != "" {
		spinner.Color(opts.Color)
	}
	if opts.Text != "" {
		spinner.Suffix = " " + opts.Text
	}
	return spinner
}

// ShowSpinnerForFunc shows a spinner during the execution of the function.
func ShowSpinnerForFunc(opts SpinnerOptions, function func()) {
	s := StartSpinner(opts)
	defer s.Stop()
	function()
}
