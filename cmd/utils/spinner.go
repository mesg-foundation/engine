package utils

import (
	"time"

	spinnerPkg "github.com/mesg-foundation/spinner"
)

var (
	// SpinnerCharset is the default animation.
	SpinnerCharset = spinnerPkg.CharSets[11]

	// SpinnerDuration is the default duration for spinning.
	SpinnerDuration = 100 * time.Millisecond
)

// SpinnerOptions contains all details for the spinner
type SpinnerOptions struct {
	Text  string
	Color string
}

// StartSpinner creates new spinner for terminal.
func StartSpinner(opts SpinnerOptions) (spinner *spinnerPkg.Spinner) {
	spinner = spinnerPkg.New(SpinnerCharset, SpinnerDuration)
	spinner.Start()
	if opts.Color != "" {
		spinner.Color(opts.Color)
	}
	if opts.Text != "" {
		spinner.Suffix = " " + opts.Text
	}
	return
}

// ShowSpinnerForFunc shows a spinner during the execution of the function.
func ShowSpinnerForFunc(opts SpinnerOptions, function func()) {
	s := StartSpinner(opts)
	defer s.Stop()
	function()
	return
}
