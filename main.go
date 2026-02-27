package main

import (
	"encoding/json"
	"fmt"
	"go-cli-db/database"
	"go-cli-db/globals"
	"go-cli-db/prompts"
	"go-cli-db/select_option"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {

	fmt.Print(globals.ResetTerminal)
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

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan,
		syscall.SIGINT,  // CTRL+C
		syscall.SIGTERM, // Terminação
	)

	go func() {
		<-sigChan
		fmt.Println("\n\nPrograma interrompido. Limpando...")

		fmt.Print(globals.ShowCursor)
		fmt.Print(globals.ResetFormat)
		fmt.Print(globals.ResetTerminal)

		os.Exit(0)
	}()
	ShowMainMenu()
}

func ShowMainMenu() {
	var mainMenu select_option.Options

	mainMenu = append(mainMenu, struct {
		Id    int
		Label string
	}{Id: 1, Label: "Create new backup config"})
	mainMenu = append(mainMenu, struct {
		Id    int
		Label string
	}{Id: 2, Label: "Start backup"})

	mainMenuOption, err := prompts.SelectPrompt("Choose your option:", mainMenu)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: "+err.Error())
	}

	if mainMenuOption == 1 {
		CreateBackup()
	}
}

func CreateBackup() {

	var Database database.DatabaseData
	Database.Host = prompts.StringPrompt("Database Host:", false)
	Database.Name = prompts.StringPrompt("Database Name:", false)
	Database.Username = prompts.StringPrompt("Database Username:", false)
	Database.Password = prompts.PasswordPrompt("Database Password:")
	ignoreTables := prompts.StringPrompt("Database Ignore Tables (separated by comma):", true)
	Database.IgnoreTables = strings.Split(ignoreTables, ",")
	fmt.Print(globals.ResetTerminal)

	jsonDatabase, err := json.Marshal(Database)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err.Error())
	}

	dbFilesDir := globals.GoDotEnvVariable("DB_FILES_DIR")
	if _, err := os.Stat(dbFilesDir); err != nil {
		err := os.Mkdir(dbFilesDir, 0750)
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}
	}
	err = os.WriteFile(dbFilesDir+Database.Name+".json", jsonDatabase, 0644)

	database.InitDbDump(Database)
}
