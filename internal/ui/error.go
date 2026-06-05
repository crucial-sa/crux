package ui

import (
	"os"

	"charm.land/lipgloss/v2"
)

var style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Red).
	PaddingTop(1)

var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Red)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))

func Error(message string, err error) {
	lipgloss.Println(style.Render("ERROR:", message))
	lipgloss.Println(errorStyle.Render("\n\t", err.Error(), "\n"))
	lipgloss.Println(helpStyle.Render("Run with --verbose or CRUX_VERBOSE=1 for more information"))
}

func Panic(message string, err error) {
	Error(message, err)
	os.Exit(1)
}
