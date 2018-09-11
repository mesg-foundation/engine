package pretty

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

var (
	// spinnerCharset is the default animation.
	spinnerCharset = spinner.CharSets[11]

	// spinnerDuration is the default duration for spinning.
	spinnerDuration = 100 * time.Millisecond
)

// pg is the default Pretty printer.
var pg = New()

// Signs to print for success, warn and fail.
var (
	SuccessSign = pg.Success("✔")
	WarnSign    = pg.Warn("?")
	FailSign    = pg.Fail("⨯")
)

// Predefiend colors
var (
	FgBlack   = color.New(color.FgBlack)
	FgYellow  = color.New(color.FgYellow)
	FgBlue    = color.New(color.FgBlue)
	FgMagenta = color.New(color.FgMagenta)
	FgCyan    = color.New(color.FgCyan)
	FgWhite   = color.New(color.FgWhite)
	FgRed     = color.New(color.FgRed)
	FgGreen   = color.New(color.FgGreen)
)

// Pretty handles pretty printing for terminal and string.
type Pretty struct {
	noColor   bool
	noSpinner bool

	successColor *color.Color
	warnColor    *color.Color
	failColor    *color.Color
	boldColor    *color.Color

	*color.Color

	*spinner.Spinner
}

// New returns a new pretty object with all values set to default one.
func New() *Pretty {
	return &Pretty{
		noColor:   runtime.GOOS == "windows",
		noSpinner: runtime.GOOS == "windows",

		successColor: color.New(color.FgGreen),
		warnColor:    color.New(color.FgYellow),
		failColor:    color.New(color.FgRed),
		boldColor:    color.New(color.Bold),

		Color: color.New(),

		Spinner: spinner.New(spinnerCharset, spinnerDuration),
	}
}

// DisableColor disables the color output.
func (p *Pretty) DisableColor() {
	p.noColor = true
	p.successColor.DisableColor()
	p.warnColor.DisableColor()
	p.failColor.DisableColor()
	p.boldColor.DisableColor()
	p.Color.DisableColor()
}

// EnableColor enables the color output.
func (p *Pretty) EnableColor() {
	// windows terminal dosen't support colors
	if runtime.GOOS == "windows" {
		return
	}

	p.noColor = false
	p.successColor.EnableColor()
	p.warnColor.EnableColor()
	p.failColor.EnableColor()
	p.boldColor.EnableColor()
	p.Color.EnableColor()
}

// DisableSpinner disables the spinner.
func (p *Pretty) DisableSpinner() {
	p.noSpinner = true
}

// EnableSpinner enables the spinner.
func (p *Pretty) EnableSpinner() {
	// windows terminal dosen't support spinner
	if runtime.GOOS == "windows" {
		return
	}

	p.noSpinner = false
}

// Success formats using the default formats for its operands and
// returns the resulting string with success color foreground.
// Spaces are added between operands when neither is a string.
func (p *Pretty) Success(msg string) string {
	return p.successColor.Sprint(msg)
}

// Successf formats according to a format specifier and
// returns the resulting string with success color foreground.
func (p *Pretty) Successf(format string, a ...interface{}) string {
	return p.successColor.Sprintf(format, a...)
}

// Successln formats using the default formats for its operands and
// returns the resulting string with success color foreground.
// Spaces are always added between operands and a newline is appended.
func (p *Pretty) Successln(a ...interface{}) string {
	return p.successColor.Sprintln(a...)
}

// Warn formats using the default formats for its operands and
// returns the resulting string with warn color foreground.
// Spaces are added between operands when neither is a string.
func (p *Pretty) Warn(msg string) string {
	return p.warnColor.Sprint(msg)
}

// Warnf formats according to a format specifier and
// returns the resulting string with warn color foreground.
func (p *Pretty) Warnf(format string, a ...interface{}) string {
	return p.warnColor.Sprintf(format, a...)
}

// Warnln formats using the default formats for its operands and
// returns the resulting string with warn color foreground.
// Spaces are always added between operands and a newline is appended.
func (p *Pretty) Warnln(a ...interface{}) string {
	return p.warnColor.Sprintln(a...)
}

// Fail formats using the default formats for its operands and
// returns the resulting string with fail color foreground.
// Spaces are added between operands when neither is a string.
func (p *Pretty) Fail(msg string) string {
	return p.failColor.Sprint(msg)
}

// Failf formats according to a format specifier and
// returns the resulting string with fail color foreground.
func (p *Pretty) Failf(format string, a ...interface{}) string {
	return p.failColor.Sprintf(format, a...)
}

// Failln formats using the default formats for its operands and
// returns the resulting string with fail color foreground.
// Spaces are always added between operands and a newline is appended.
func (p *Pretty) Failln(a ...interface{}) string { return p.failColor.Sprintln(a...) }

// Bold formats using the default formats for its operands and
// returns the resulting bolded string.
// Spaces are added between operands when neither is a string.
func (p *Pretty) Bold(msg string) string {
	return p.boldColor.Sprint(msg)
}

// Boldf formats according to a format specifier and
// returns the resulting bolded string.
func (p *Pretty) Boldf(format string, a ...interface{}) string {
	return p.boldColor.Sprintf(format, a...)
}

// Boldln formats using the default formats for its operands and
// returns the resulting bolded string.
// Spaces are always added between operands and a newline is appended.
func (p *Pretty) Boldln(a ...interface{}) string {
	return p.boldColor.Sprintln(a...)
}

// Colorize formats using the default formats for its operands and
// returns the resulting string with given color.
// Spaces are added between operands when neither is a string.
func (p *Pretty) Colorize(c *color.Color, msg string) string {
	if p.noColor {
		return msg
	}
	return c.Sprint(msg)
}

// ColorizeJSON colors keys and values of stringified JSON. On errors the original string is returned.
// If color is nil then key/value won't be colorize.
func (p *Pretty) ColorizeJSON(keyColor *color.Color, valueColor *color.Color, data []byte) []byte {
	if p.noColor {
		return data
	}

	var (
		in  map[string]interface{}
		out map[string]interface{}
	)

	if json.Unmarshal(data, &in) != nil {
		return data
	}

	if keyColor == nil {
		keyColor = color.New()
	}
	if valueColor == nil {
		valueColor = color.New()
	}

	for k, v := range in {
		out[keyColor.Sprint(k)] = valueColor.Sprint(v)
	}

	b, err := json.Marshal(out)
	if err != nil {
		return data
	}

	return b
}

// Progress prints spinner with the given message while calling fn function.
func (p *Pretty) Progress(message string, fn func()) {
	if p.noSpinner {
		fmt.Println(message)
		fn()
		return
	}

	p.Spinner.Suffix = " " + message
	p.Spinner.Start()
	fn()
	p.Spinner.Stop()
	p.Spinner.Suffix = ""

	if !strings.HasSuffix(message, "\n") {
		fmt.Println()
	}
}

// FgColors returns a slice with predefiend foreground color.
func (p *Pretty) FgColors() []*color.Color {
	return []*color.Color{
		FgBlack,
		FgYellow,
		FgBlue,
		FgMagenta,
		FgCyan,
		FgWhite,
		FgRed,
		FgGreen,
	}
}

// Default returns the default Pretty printer.
func Default() *Pretty { return pg }

// DisableColor disables the color output.
func DisableColor() { pg.DisableColor() }

// EnableColor enables the color output.
func EnableColor() { pg.EnableColor() }

// DisableSpinner disables the spinner.
func DisableSpinner() { pg.DisableSpinner() }

// EnableSpinner enables the spinner.
func EnableSpinner() { pg.EnableSpinner() }

// Success formats using the default formats for its operands and
// returns the resulting string with success color foreground.
// Spaces are added between operands when neither is a string.
func Success(msg string) string { return pg.Success(msg) }

// Successf formats according to a format specifier and
// returns the resulting string with success color foreground.
func Successf(format string, a ...interface{}) string { return pg.Successf(format, a...) }

// Successln formats using the default formats for its operands and
// returns the resulting string with success color foreground.
// Spaces are always added between operands and a newline is appended.
func Successln(a ...interface{}) string { return pg.Successln(a...) }

// Warn formats using the default formats for its operands and
// returns the resulting string with warn color foreground.
// Spaces are added between operands when neither is a string.
func Warn(msg string) string { return pg.Warn(msg) }

// Warnf formats according to a format specifier and
// returns the resulting string with warn color foreground.
func Warnf(format string, a ...interface{}) string { return pg.Warnf(format, a...) }

// Warnln formats using the default formats for its operands and
// returns the resulting string with warn color foreground.
// Spaces are always added between operands and a newline is appended.
func Warnln(a ...interface{}) string { return pg.Warnln(a...) }

// Fail formats using the default formats for its operands and
// returns the resulting string with fail color foreground.
// Spaces are added between operands when neither is a string.
func Fail(msg string) string { return pg.Fail(msg) }

// Failf formats according to a format specifier and
// returns the resulting string with fail color foreground.
func Failf(format string, a ...interface{}) string { return pg.Failf(format, a...) }

// Failln formats using the default formats for its operands and
// returns the resulting string with fail color foreground.
// Spaces are always added between operands and a newline is appended.
func Failln(a ...interface{}) string { return pg.Failln(a...) }

// Bold formats using the default formats for its operands and
// returns the resulting bolded string.
// Spaces are added between operands when neither is a string.
func Bold(msg string) string { return pg.Bold(msg) }

// Boldf formats according to a format specifier and
// returns the resulting bolded string.
func Boldf(format string, a ...interface{}) string { return pg.Boldf(format, a...) }

// Boldln formats using the default formats for its operands and
// returns the resulting bolded string.
// Spaces are always added between operands and a newline is appended.
func Boldln(a ...interface{}) string { return pg.Boldln(a...) }

// Colorize formats using the default formats for its operands and
// returns the resulting string with given color.
// Spaces are added between operands when neither is a string.
func Colorize(c *color.Color, msg string) string { return pg.Colorize(c, msg) }

// Progress prints spinner with the given message while calling fn function.
func Progress(message string, fn func()) { pg.Progress(message, fn) }

// ColorizeJSON colors keys and values of stringified JSON. On errors the original string is returned.
// If color is nil then key/value won't be colorize.
func ColorizeJSON(keyColor *color.Color, valueColor *color.Color, data []byte) []byte {
	return pg.ColorizeJSON(keyColor, valueColor, data)
}

// FgColors returns a slice with predefiend foreground color.
func FgColors() []*color.Color {
	return pg.FgColors()
}
