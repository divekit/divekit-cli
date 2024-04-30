package dye

import "github.com/fatih/color"

var (
	Red     = color.New(color.FgRed).SprintFunc()
	Green   = color.New(color.FgGreen).SprintFunc()
	Yellow  = color.New(color.FgYellow).SprintFunc()
	Blue    = color.New(color.FgBlue).SprintFunc()
	Magenta = color.New(color.FgMagenta).SprintFunc()
	Cyan    = color.New(color.FgCyan).SprintFunc()
	White   = color.New(color.FgWhite).SprintFunc()
	Black   = color.New(color.FgBlack).SprintFunc()
	Grey    = color.New(color.FgHiBlack).SprintFunc()
	Bold    = color.New(color.Bold).SprintFunc()
	Under   = color.New(color.Underline).SprintFunc()
	Italic  = color.New(color.Italic).SprintFunc()
)

var Very = struct {
	Red     func(a ...interface{}) string
	Green   func(a ...interface{}) string
	Yellow  func(a ...interface{}) string
	Blue    func(a ...interface{}) string
	Magenta func(a ...interface{}) string
	Cyan    func(a ...interface{}) string
	White   func(a ...interface{}) string
	Black   func(a ...interface{}) string
	Grey    func(a ...interface{}) string
}{
	Red:     color.New(color.FgHiRed).SprintFunc(),
	Green:   color.New(color.FgHiGreen).SprintFunc(),
	Yellow:  color.New(color.FgHiYellow).SprintFunc(),
	Blue:    color.New(color.FgHiBlue).SprintFunc(),
	Magenta: color.New(color.FgHiMagenta).SprintFunc(),
	Cyan:    color.New(color.FgHiCyan).SprintFunc(),
	White:   color.New(color.FgHiWhite).SprintFunc(),
	Black:   color.New(color.FgHiBlack).SprintFunc(),
}
