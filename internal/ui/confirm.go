package ui

import "charm.land/huh/v2"

func Confirm(message string) bool {
	var confirm bool

	huh.NewConfirm().
		Title(message).
		Affirmative("Yes").
		Negative("No").
		Value(&confirm).WithTheme(Theme()).Run()

	return confirm
}
