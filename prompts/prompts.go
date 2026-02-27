package prompts

import (
	"bufio"
	"fmt"
	"go-cli-db/select_option"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func StringPrompt(label string, emptyAllowed bool) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stdout, label+" ")
		s, _ = r.ReadString('\n')
		s = strings.TrimSpace(s)
		if s != "" || s == "" && emptyAllowed {
			break
		}
	}
	return strings.TrimSpace(s)
}

func PasswordPrompt(label string) string {
	var s string
	for {
		fmt.Fprint(os.Stdout, label+" ")
		b, _ := term.ReadPassword(int(syscall.Stdin))
		s = string(b)
		if s != "" {
			break
		}
	}
	fmt.Println()
	return s
}

func YesOrNoPrompt(label string, def bool) bool {
	var s string
	r := bufio.NewReader(os.Stdin)
	choices := "Y/n"
	if !def {
		choices = "n/Y"
	}

	for {
		fmt.Fprintf(os.Stdout, "%s (%s) ", label, choices)
		s, _ = r.ReadString('\n')
		s = strings.TrimSpace(s)
		if s == "" {
			return def
		}
		s = strings.ToLower(s)
		if s == "y" || s == "yes" {
			return true
		}
		if s == "n" || s == "no" {
			return false
		}
	}
}

func SelectPrompt(label string, opts select_option.Options) (int, error) {

	selectOptions := select_option.NewSelectMenu(label, opts)
	selectedOption, err := selectOptions.DisplayMenu()
	return selectedOption, err
}
