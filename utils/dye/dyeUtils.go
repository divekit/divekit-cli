package dye

import "github.com/fatih/color"

var (
	Yellow  = color.New(color.FgYellow).SprintFunc()
	Grey    = color.New(color.FgHiBlack).SprintFunc()
	Red     = color.New(color.FgRed).SprintFunc()
	Green   = color.New(color.FgGreen).SprintFunc()
	Blue    = color.New(color.FgBlue).SprintFunc()
	Bold    = color.New(color.Bold).SprintFunc()
	White   = color.New(color.FgWhite).SprintFunc()
	Black   = color.New(color.FgBlack).SprintFunc()
	Cyan    = color.New(color.FgCyan).SprintFunc()
	Magenta = color.New(color.FgMagenta).SprintFunc()
	Orange  = color.New(color.FgHiRed).SprintFunc()     // closest to orange
	Purple  = color.New(color.FgHiMagenta).SprintFunc() // closest to purple
)
