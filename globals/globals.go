package globals

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

const (
	ShowCursor    = "\033[?25h"
	HideCursor    = "\033[?25l"
	ClearLine     = "\r\033[K"
	ClearScreen   = "\033[2J"
	ResetFormat   = "\033[0m"
	ResetTerminal = "\033[2J\033[H"
	PutOnStart    = "\033[1;1H"
	KeyUp         = byte(65)
	KeyDown       = byte(66)
	KeyEscape     = byte(27)
	KeyEnter      = byte(13)
	KeyCtrlC      = byte(3)
)

func GoDotEnvVariable(key string) string {

	homeDir, _ := os.UserHomeDir()
	envPath := filepath.Join(homeDir, ".go-cli-db", ".env")
	err := godotenv.Load(envPath)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
