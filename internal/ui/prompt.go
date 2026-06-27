package ui

import "charm.land/huh/v2"

type askOptions struct {
	description string
	placeholder string
	validate    func(string) error
}

type AskOption func(*askOptions)

func WithDescription(description string) AskOption {
	return func(o *askOptions) { o.description = description }
}

func WithPlaceholder(placeholder string) AskOption {
	return func(o *askOptions) { o.placeholder = placeholder }
}

func WithValidation(validate func(string) error) AskOption {
	return func(o *askOptions) { o.validate = validate }
}

func Ask(question string, answer *string, opts ...AskOption) error {
	var o askOptions
	for _, opt := range opts {
		opt(&o)
	}

	input := huh.NewInput().
		Title(question).
		Value(answer)

	if o.description != "" {
		input = input.Description(o.description)
	}

	if o.placeholder != "" {
		input = input.Placeholder(o.placeholder)
	}

	if o.validate != nil {
		input = input.Validate(o.validate)
	}

	err := input.WithTheme(Theme()).Run()

	if err == nil {
		Say(question, *answer)
	}

	return err
}
