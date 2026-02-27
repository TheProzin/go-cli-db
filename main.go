package main

import (
	"encoding/json"
	"fmt"
	"go-cli-db/database"
	"go-cli-db/globals"
	"go-cli-db/prompts"
	"go-cli-db/select_option"
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
	mainMenu = append(mainMenu, struct {
		Id    int
		Label string
	}{Id: 3, Label: "Restore dump"})
	mainMenu = append(mainMenu, struct {
		Id    int
		Label string
	}{Id: 9, Label: "Quit"})

loop:
	for {
		fmt.Print(globals.ResetTerminal)
		mainMenuOption, err := prompts.SelectPrompt("Choose your option:", mainMenu)

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: "+err.Error())
			break loop
		}

		switch mainMenuOption {
		case 1:
			CreateBackup()
		case 2:
			ChooseDbFileForBackup()
		case 3:
			ChooseDbFileForRestore()
		case 9:
			break loop
		}
	}
}

func CreateBackup() {

	var dbConfig database.DatabaseData
	dbConfig.ConfigName = prompts.StringPrompt("Configuration name:", false)
	dbConfig.DbHost = prompts.StringPrompt("Database Host:", false)
	dbConfig.DbName = prompts.StringPrompt("Database Name:", false)
	dbConfig.DbUsername = prompts.StringPrompt("Database Username:", false)
	dbConfig.DbPassword = prompts.PasswordPrompt("Database Password:")
	ignoreTables := prompts.StringPrompt("Database Ignore Tables (separated by comma):", true)
	dbConfig.IgnoreTables = strings.Split(ignoreTables, ",")
	fmt.Print(globals.ResetTerminal)

	jsonDatabase, err := json.Marshal(dbConfig)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err.Error())
		return
	}

	dbFilesDir := globals.GoDotEnvVariable("DB_FILES_DIR")
	if _, err := os.Stat(dbFilesDir); err != nil {
		err := os.Mkdir(dbFilesDir, 0750)
		if err != nil && !os.IsExist(err) {
			fmt.Fprintln(os.Stderr, "Error: ", err.Error())
			return
		}
	}

	err = os.WriteFile(dbFilesDir+dbConfig.ConfigName+".json", jsonDatabase, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err.Error())
		return
	}

	startBkp := prompts.YesOrNoPrompt("Start backup?", true)
	if startBkp {
		sqlDumpFileName, err := database.InitDbDump(dbConfig)
		if err != nil {
			return
		}
		restoreFile := prompts.YesOrNoPrompt("Do you want to restore the file?", true)
		if restoreFile {
			dbConfig, err := chooseDbConfig()
			if err != nil {
				return
			}

			database.InitDbRestore(sqlDumpFileName, dbConfig)
		}
	}
}

func ChooseDbFileForBackup() {

	dbConfig, err := chooseDbConfig()
	if err != nil {
		return
	}

	startBkp := prompts.YesOrNoPrompt("Start backup?", true)
	if startBkp {
		sqlDumpFileName, err := database.InitDbDump(dbConfig)
		if err != nil {
			return
		}
		restoreFile := prompts.YesOrNoPrompt("Do you want to restore the file?", true)
		if restoreFile {
			dbConfig, err := chooseDbConfig()
			if err != nil {
				return
			}

			database.InitDbRestore(sqlDumpFileName, dbConfig)
		}
	}
}

func ChooseDbFileForRestore() {

	dbConfig, err := chooseDbConfig()
	if err != nil {
		return
	}

	sqlDumpFileName, err := chooseDumpFileName()
	if err != nil {
		return
	}

	fmt.Fprintf(os.Stdout, "Restoring from: %s \nTo db: %s with host: %s\n", sqlDumpFileName, dbConfig.DbName, dbConfig.DbHost)
	startBkp := prompts.YesOrNoPrompt("Start backup?", true)
	if startBkp {
		database.InitDbRestore(sqlDumpFileName, dbConfig)
	}
}

func chooseDbConfig() (database.DatabaseData, error) {

	dbConfigFilesDir := globals.GoDotEnvVariable("DB_FILES_DIR")
	files, err := os.ReadDir(dbConfigFilesDir)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err.Error())
		return database.DatabaseData{}, err
	}

	if len(files) == 0 {
		fmt.Fprintln(os.Stderr, "Don't exist any config file")
		return database.DatabaseData{}, err
	}

	var filesOptions select_option.Options
	for i, file := range files {
		filesOptions = append(filesOptions, struct {
			Id    int
			Label string
		}{
			Id:    i,
			Label: file.Name(),
		})
	}

	fileDbSelected, err := prompts.SelectPrompt("Choose your db config:", filesOptions)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err.Error())
		return database.DatabaseData{}, err
	}

	configDbJson := files[fileDbSelected]

	file, err := os.Open(dbConfigFilesDir + configDbJson.Name())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err.Error())
		return database.DatabaseData{}, err
	}
	defer file.Close()

	var dbConfig database.DatabaseData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&dbConfig); err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err.Error())
		return database.DatabaseData{}, err
	}

	return dbConfig, nil
}

func chooseDumpFileName() (string, error) {

	dbBkpFilesDir := globals.GoDotEnvVariable("DB_BACKUP_DIR")
	filesFrom, err := os.ReadDir(dbBkpFilesDir)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err.Error())
		return "", err
	}

	if len(filesFrom) == 0 {
		fmt.Fprintln(os.Stderr, "Don't exist any dump file")
		return "", err
	}

	var filesFromOptions select_option.Options
	for i, file := range filesFrom {
		filesFromOptions = append(filesFromOptions, struct {
			Id    int
			Label string
		}{
			Id:    i,
			Label: file.Name(),
		})
	}

	fileWhereFromDbSelected, err := prompts.SelectPrompt("Choose your sql dump:", filesFromOptions)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err.Error())
		return "", err
	}

	sqlDumpFile := filesFrom[fileWhereFromDbSelected]

	return sqlDumpFile.Name(), nil
}
