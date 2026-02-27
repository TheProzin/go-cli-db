package select_option

import (
	"fmt"
	"go-cli-db/globals"
	"os"
	"strings"

	"golang.org/x/term"
)

type SelectOptions struct {
	Label             string
	Options           Options
	PointerPosition   int
	lastRenderedLines int
}

type Options []struct {
	Id    int
	Label string
}

func NewSelectMenu(label string, options Options) *SelectOptions {
	return &SelectOptions{
		Options: options,
		Label:   label,
	}
}

func (selectOptions *SelectOptions) DisplayMenu() (int, error) {
	selectOptions.RenderMenu()
	defer selectOptions.ResetScreen()

	for {
		keyCode, err := getInput()
		if err != nil {
			return 0, fmt.Errorf("error processing key input on select option")
		}

		switch keyCode {
		case globals.KeyUp:
			selectOptions.PointerPosition = (selectOptions.PointerPosition - 1 + len(selectOptions.Options)) % len(selectOptions.Options)
			selectOptions.RenderMenu()
		case globals.KeyDown:
			selectOptions.PointerPosition = (selectOptions.PointerPosition + 1) % len(selectOptions.Options)
			selectOptions.RenderMenu()
		case globals.KeyEnter:
			return selectOptions.Options[selectOptions.PointerPosition].Id, nil
		case globals.KeyCtrlC:
			return 0, fmt.Errorf("exit program")
		}
	}
}

func (selectOptions *SelectOptions) ResetScreen() {
	if selectOptions.lastRenderedLines > 0 {
		fmt.Printf("\033[%dA", selectOptions.lastRenderedLines)
		for i := 0; i < selectOptions.lastRenderedLines; i++ {
			fmt.Print("\r\033[K\n")
		}
		fmt.Printf("\033[%dA", selectOptions.lastRenderedLines)
	}
	selectOptions.lastRenderedLines = 0
	fmt.Print(globals.ShowCursor)
}

func (selectOptions *SelectOptions) RenderMenu() {
	if selectOptions.lastRenderedLines > 0 {
		fmt.Printf("\033[%dA", selectOptions.lastRenderedLines)
	}

	lines := 0

	fmt.Printf("\r\033[K%s\n", selectOptions.Label)
	lines++

	for optionIndex, option := range selectOptions.Options {
		optionSelected := " "
		if optionIndex == selectOptions.PointerPosition {
			optionSelected = "X"
		}
		option.Label = strings.TrimLeft(option.Label, " \t")
		fmt.Printf("\r\033[K[%s] %s\n", optionSelected, option.Label)
		lines++
	}

	selectOptions.lastRenderedLines = lines
	fmt.Print(globals.HideCursor)
}

func getInput() (byte, error) {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return 0, err
	}

	buffer := make([]byte, 16)
	n, err := os.Stdin.Read(buffer)

	term.Restore(fd, oldState)

	if err != nil {
		return 0, err
	}

	if n == 1 {
		switch buffer[0] {
		case globals.KeyEnter:
			return globals.KeyEnter, nil
		case globals.KeyCtrlC:
			return globals.KeyCtrlC, nil
		case 'q':
			return globals.KeyCtrlC, nil
		}
	}

	if n >= 3 && buffer[0] == 0x1B && buffer[1] == '[' {
		switch buffer[2] {
		case 'A':
			return globals.KeyUp, nil
		case 'B':
			return globals.KeyDown, nil
		}
	}

	return 0, fmt.Errorf("unknown key")
}
