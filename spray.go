package main

import(
"github.com/fatih/color"
)

//here functions to paint text
func rspray(text string) string {
	_color := color.New(color.FgRed)
	return _color.Sprint(text)
}

func mspray(text string) string {
	_color := color.New(color.FgMagenta)
	return _color.Sprint(text)
}

func yspray(text string) string {
	_color := color.New(color.FgYellow)
	return _color.Sprint(text)
}

func bspray(text string) string {
	_color := color.New(color.FgBlue)
	return _color.Sprint(text)
}

