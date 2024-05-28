package spray

/*
Package represenes colorful text to debugging and logging data in console
First letter means color
R = red
G = green
M = magenta
Y = yellow
B = blue
*/

import (
	"github.com/fatih/color"
)

func Rspray(text string) string {
	_color := color.New(color.FgRed)
	return _color.Sprint(text)
}

func Gspray(text string) string {
	_color := color.New(color.FgGreen)
	return _color.Sprint(text)
}

func Mspray(text string) string {
	_color := color.New(color.FgMagenta)
	return _color.Sprint(text)
}

func Yspray(text string) string {
	_color := color.New(color.FgYellow)
	return _color.Sprint(text)
}

func Bspray(text string) string {
	_color := color.New(color.FgBlue)
	return _color.Sprint(text)
}
