package cmdUtils

import "github.com/fatih/color"

// SuccessColor is a function to format using the success color
var SuccessColor = color.New(color.FgGreen, color.Bold).SprintFunc()

// WarningColor is a function to format using the warning color
var WarningColor = color.New(color.FgYellow, color.Bold).SprintFunc()
