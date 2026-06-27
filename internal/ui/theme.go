package ui

import (
	"image/color"

	"charm.land/huh/v2"
	"charm.land/huh/v2/spinner"
	"charm.land/lipgloss/v2"
)

var (
	colorFg1     color.Color = lipgloss.NoColor{}
	colorFg2     color.Color = lipgloss.BrightBlack
	colorFg3     color.Color = lipgloss.BrightBlack
	colorMuted   color.Color = lipgloss.BrightBlack
	colorBorder  color.Color = lipgloss.BrightBlack
	colorAccent  color.Color = lipgloss.Blue
	colorOnAccnt color.Color = lipgloss.BrightWhite
	colorRose    color.Color = lipgloss.Red
)

func Theme() huh.Theme {
	return huh.ThemeFunc(formStyles)
}

func formStyles(isDark bool) *huh.Styles {
	t := huh.ThemeBase(isDark)

	// Group / field container
	t.Focused.Base = t.Focused.Base.BorderForeground(colorAccent)
	t.Focused.Card = t.Focused.Base

	// Titles & descriptions.
	t.Focused.Title = t.Focused.Title.Foreground(colorAccent).Bold(true)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(colorAccent).Bold(true).MarginBottom(1)
	t.Focused.Description = t.Focused.Description.Foreground(colorFg3)
	t.Focused.Directory = t.Focused.Directory.Foreground(colorAccent)
	t.Focused.File = t.Focused.File.Foreground(colorFg1)

	// Errors.
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(colorRose)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(colorRose)

	// Select / multi-select.
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(colorAccent)
	t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(colorAccent)
	t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(colorAccent)
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(colorAccent)
	t.Focused.Option = t.Focused.Option.Foreground(colorFg1)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(colorAccent)
	t.Focused.SelectedPrefix = lipgloss.NewStyle().Foreground(colorAccent).SetString("✓ ")
	t.Focused.UnselectedPrefix = lipgloss.NewStyle().Foreground(colorMuted).SetString("• ")
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(colorFg2)

	// Buttons
	t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(colorOnAccnt).Background(colorAccent).Bold(true)
	t.Focused.BlurredButton = t.Focused.BlurredButton.Foreground(colorFg2).Background(lipgloss.NoColor{})
	t.Focused.Next = t.Focused.FocusedButton

	// Text input.
	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(colorAccent)
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(colorAccent)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(colorMuted)
	t.Focused.TextInput.Text = t.Focused.TextInput.Text.Foreground(colorFg1)

	// Blurred mirrors focused.
	t.Blurred = t.Focused
	t.Blurred.Base = t.Blurred.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.Card = t.Blurred.Base
	t.Blurred.NextIndicator = lipgloss.NewStyle()
	t.Blurred.PrevIndicator = lipgloss.NewStyle()

	// Help footer.
	t.Help.ShortKey = t.Help.ShortKey.Foreground(colorFg2)
	t.Help.ShortDesc = t.Help.ShortDesc.Foreground(colorFg3)
	t.Help.ShortSeparator = t.Help.ShortSeparator.Foreground(colorBorder)
	t.Help.FullKey = t.Help.FullKey.Foreground(colorFg2)
	t.Help.FullDesc = t.Help.FullDesc.Foreground(colorFg3)
	t.Help.FullSeparator = t.Help.FullSeparator.Foreground(colorBorder)
	t.Help.Ellipsis = t.Help.Ellipsis.Foreground(colorMuted)

	t.Group.Title = t.Focused.Title
	t.Group.Description = t.Focused.Description
	return t
}

func SpinnerTheme() spinner.Theme {
	return spinner.ThemeFunc(func(isDark bool) *spinner.Styles {
		return &spinner.Styles{
			Spinner: lipgloss.NewStyle().Foreground(colorAccent),
			Title:   lipgloss.NewStyle().Foreground(colorFg2),
		}
	})
}
