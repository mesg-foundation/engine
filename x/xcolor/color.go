package xcolor

import "github.com/fatih/color"

var colors = []color.Attribute{
	color.FgMagenta,
	color.FgYellow,
	color.FgBlue,
	color.FgCyan,
	color.FgGreen,
	color.FgRed,
}

var (
	colorsLen = len(colors)
	i         int
)

// NextColor returns the next color attribute.
func NextColor() color.Attribute {
	c := colors[i]
	i = (i + 1) % colorsLen
	return c
}
