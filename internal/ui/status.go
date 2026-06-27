package ui

import (
	"strings"

	"charm.land/lipgloss/v2"
)

func StatusStyle(status string) string {
	style := lipgloss.NewStyle()

	switch strings.ToUpper(status) {
	case "RUNNING":
		style = style.Foreground(lipgloss.Green)
	case "PENDING", "BUILDING":
		style = style.Foreground(lipgloss.Yellow)
	case "FAILED":
		style = style.Foreground(lipgloss.Red)
	case "STOPPED":
		style = style.Foreground(lipgloss.BrightBlack)
	}

	return style.Render(status)
}
