package ui

import (
	"strings"

	"charm.land/lipgloss/v2"
)

func Summary(title string, fields [][2]string) {
	width := 0
	for _, f := range fields {
		if len(f[0]) > width {
			width = len(f[0])
		}
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Blue)
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.BrightBlack).Width(width + 2)

	var b strings.Builder
	if title != "" {
		b.WriteString(titleStyle.Render(title))
		b.WriteString("\n")
	}
	for _, f := range fields {
		b.WriteString("  ")
		b.WriteString(labelStyle.Render(f[0]))
		b.WriteString(f[1])
		b.WriteString("\n")
	}

	Say(b.String())
}
