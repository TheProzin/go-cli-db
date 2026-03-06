package globals

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
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
