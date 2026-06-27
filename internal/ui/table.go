package ui

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
)

var (
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Blue).PaddingRight(3)
	cellStyle   = lipgloss.NewStyle().PaddingRight(3)
)

func RenderTable(headers []string, rows [][]string) string {
	return table.New().
		Border(lipgloss.HiddenBorder()).
		BorderTop(false).BorderBottom(false).
		BorderLeft(false).BorderRight(false).
		BorderColumn(false).
		BorderRow(false).
		BorderHeader(false). // removes the blank line under the header

		Headers(headers...).
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle
			}
			return cellStyle
		}).
		Render()
}
