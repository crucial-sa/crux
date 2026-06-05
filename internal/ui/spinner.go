package ui

import (
	"context"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2/spinner"
)

func Spinner(title string, ctx context.Context, action func(context.Context) error) error {
	s := spinner.New().
		Title(title).
		Context(ctx).
		WithTheme(SpinnerTheme()).
		WithViewHook(func(v tea.View) tea.View {
			v.ProgressBar = tea.NewProgressBar(tea.ProgressBarIndeterminate, 1)
			return v
		})

	if action != nil {
		s = s.ActionWithErr(action)
	}

	return s.Run()
}
