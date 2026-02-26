package main

import (
	"bufio"
	"fmt"
	"go-cli-db/select_option"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func main() {
	// name := StringPrompt("Name?")
	// fmt.Printf("Hello %s", name)
	// fmt.Println()

	// pass := PasswordPrompt("Pass?")
	// fmt.Printf("pass %s", pass)
	// fmt.Println()

	// doit := YesOrNoPrompt("Do it?", true)
	// if doit {
	// 	fmt.Println("Let's go")
	// } else {
	// 	fmt.Println("Sad")
	// }

	fruits := make(map[int]string)
	fruits[0] = "Banana"
	fruits[1] = "Maca"
	fruits[3] = "Pera"
	fruits[4] = "Laranja"
	fruits[5] = "Melao"
	SelectPrompt("Fruits", fruits)
}

func StringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

func PasswordPrompt(label string) string {
	var s string
	for {
		fmt.Fprint(os.Stderr, label+" ")
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
		fmt.Fprintf(os.Stderr, "%s (%s) ", label, choices)
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

func SelectPrompt(label string, opts map[int]string) int {

	selectOptions := select_option.NewSelectMenu(opts)

	// r := bufio.NewReader(os.Stdin)
	// var s string
	selectOptions.DisplayMenu()
	return 1
}
