package select_option

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

const (
	ShowCursor = "\033[?25h"
	HideCursor = "\033[?25l"
	ClearLine  = "\r\033[K"
	KeyUp      = byte(65)
	KeyDown    = byte(66)
	KeyEscape  = byte(27)
	KeyEnter   = byte(13)
)

type SelectOptions struct {
	Label           string
	Options         map[int]string
	PointerPosition int
}

func NewSelectMenu(options map[int]string) *SelectOptions {
	selectOptions := SelectOptions{
		Options: options,
	}

	return &selectOptions
}

func (selectOptions *SelectOptions) RenderMenu() {
	fmt.Fprintf(os.Stderr, "%s:", selectOptions.Label)
	fmt.Printf(HideCursor)
	fmt.Println()
	for optionId, option := range selectOptions.Options {
		optionSelected := ""
		if optionId == selectOptions.PointerPosition {
			optionSelected = "X"
		}
		fmt.Printf("[%s] %s", optionSelected, option)
		// fmt.Println(ClearLine)
	}
}

func (selectOptions *SelectOptions) DisplayMenu() int {
	selectOptions.RenderMenu()

	for {
		keyCode := getInput()

		switch keyCode {
		case KeyUp:
			selectOptions.PointerPosition = ((selectOptions.PointerPosition + len(selectOptions.Options)) + 1) % len(selectOptions.Options)
			selectOptions.RenderMenu()
		case KeyDown:
			selectOptions.PointerPosition = ((selectOptions.PointerPosition + len(selectOptions.Options)) - 1) % len(selectOptions.Options)
			selectOptions.RenderMenu()
		case KeyEnter:
			return selectOptions.PointerPosition
		}
	}

}

func getInput() byte {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// fmt.Println("Pressione teclas (q para sair):")

	buffer := make([]byte, 16)

	for {
		n, err := os.Stdin.Read(buffer)
		if err != nil {
			panic(err)
		}

		// Tecla 'q' para sair
		if n == 1 && buffer[0] == 'q' {
			// fmt.Println("\nSaindo...")
			break
		}

		// Interpretar sequências
		if n == 1 {
			// Caractere normal
			// b := buffer[0]
			// fmt.Printf("ASCII/Unicode: %d - Caractere: %q\n", b, rune(b))
		} else {
			// Sequência de escape
			// fmt.Printf("Sequência de escape - Bytes: %v\n", buffer[:n])

			// Detectar teclas especiais
			if n >= 3 && buffer[0] == 0x1B && buffer[1] == '[' {
				switch buffer[2] {
				case 'A':
					return KeyUp
				case 'B':
					return KeyDown
				}
			}
		}
	}

	return KeyEscape
}
