package prompts

import (
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
)

type Option struct {
	Id    int
	Label string
}

type Options []Option

func StringPrompt(label string, emptyAllowed bool) string {
	var result string
	field := huh.NewInput().Title(label).Value(&result)
	if !emptyAllowed {
		field = field.Validate(func(s string) error {
			if s == "" {
				return fmt.Errorf("campo obrigatório")
			}
			return nil
		})
	}
	err := huh.NewForm(huh.NewGroup(field)).Run()
	handleErr(err)
	return result
}

func PasswordPrompt(label string) string {
	var result string
	err := huh.NewForm(huh.NewGroup(
		huh.NewInput().
			Title(label).
			EchoMode(huh.EchoModePassword).
			Value(&result).
			Validate(func(s string) error {
				if s == "" {
					return fmt.Errorf("campo obrigatório")
				}
				return nil
			}),
	)).Run()
	handleErr(err)
	return result
}

func YesOrNoPrompt(label string, def bool) bool {
	result := def
	err := huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title(label).
			Value(&result),
	)).Run()
	handleErr(err)
	return result
}

func SelectPrompt(label string, opts Options) (int, error) {
	var selected int
	if len(opts) > 0 {
		selected = opts[0].Id
	}
	huhOpts := make([]huh.Option[int], len(opts))
	for i, o := range opts {
		huhOpts[i] = huh.NewOption(o.Label, o.Id)
	}
	err := huh.NewForm(huh.NewGroup(
		huh.NewSelect[int]().
			Title(label).
			Options(huhOpts...).
			Value(&selected),
	)).Run()
	if errors.Is(err, huh.ErrUserAborted) {
		return 0, err
	}
	return selected, err
}

func handleErr(err error) {
	if err == nil {
		return
	}
	if errors.Is(err, huh.ErrUserAborted) {
		os.Exit(0)
	}
	fmt.Fprintln(os.Stderr, "Error:", err)
}
