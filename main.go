package main

import (
	"encoding/json"
	"fmt"
	"go-cli-db/database"
	"go-cli-db/globals"
	"go-cli-db/prompts"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\n\nExit program.")
		os.Exit(0)
	}()

	err := VerifyEnvFile()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err.Error())
	}
	ShowMainMenu()
}

func VerifyEnvFile() error {
	homeDir, _ := os.UserHomeDir()
	envPath := filepath.Join(homeDir, ".go-cli-db", ".env")

	if _, err := os.Stat(envPath); err != nil {
		fmt.Fprintln(os.Stdout, "Describe the path like '/tmp/db/'")
		dbConfigPath := prompts.StringPrompt("Path to database configuration files:", false)
		dbDumpPath := prompts.StringPrompt("Path to database dump files:", false)

		err := os.MkdirAll(filepath.Join(homeDir, ".go-cli-db"), 0755)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: ", err.Error())
			return err
		}

		content := fmt.Sprintf("DB_BACKUP_DIR=%s\nDB_FILES_DIR=%s\n", dbConfigPath, dbDumpPath)
		err = os.WriteFile(filepath.Join(homeDir, ".go-cli-db", ".env"), []byte(content), 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: ", err.Error())
			return err
		}
	}
	return nil
}

func ShowMainMenu() {
	mainMenu := prompts.Options{
		{Id: 1, Label: "Create new backup config"},
		{Id: 2, Label: "Start backup"},
		{Id: 3, Label: "Restore dump"},
		{Id: 4, Label: "Change files path"},
		{Id: 9, Label: "Quit"},
	}

	for {
		mainMenuOption, err := prompts.SelectPrompt("Choose your option:", mainMenu)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: "+err.Error())
			break
		}

		switch mainMenuOption {
		case 1:
			CreateBackup()
		case 2:
			ChooseDbFileForBackup()
		case 3:
			ChooseDbFileForRestore()
		case 4:
			ChangeFilePaths()
		case 9:
			return
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

	if prompts.YesOrNoPrompt("Start backup?", true) {
		sqlDumpFileName, err := database.InitDbDump(dbConfig)
		if err != nil {
			return
		}
		if prompts.YesOrNoPrompt("Do you want to restore the file?", true) {
			dbConfig, err := chooseDbConfig()
			if err != nil {
				return
			}
			fmt.Fprintf(os.Stdout, "Restoring from: %s \nTo db: %s with host: %s\n", sqlDumpFileName, dbConfig.DbName, dbConfig.DbHost)
			if prompts.YesOrNoPrompt("Start restore?", true) {
				database.InitDbRestore(sqlDumpFileName, dbConfig)
			}
		}
	}
}

func ChooseDbFileForBackup() {
	dbConfig, err := chooseDbConfig()
	if err != nil {
		return
	}

	if prompts.YesOrNoPrompt("Start backup?", true) {
		sqlDumpFileName, err := database.InitDbDump(dbConfig)
		if err != nil {
			return
		}
		if prompts.YesOrNoPrompt("Do you want to restore the file?", true) {
			dbConfig, err := chooseDbConfig()
			if err != nil {
				return
			}
			fmt.Fprintf(os.Stdout, "Restoring from: %s \nTo db: %s with host: %s\n", sqlDumpFileName, dbConfig.DbName, dbConfig.DbHost)
			if prompts.YesOrNoPrompt("Start restore?", true) {
				database.InitDbRestore(sqlDumpFileName, dbConfig)
			}
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
	if prompts.YesOrNoPrompt("Start restore?", true) {
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
		fmt.Fprintln(os.Stderr, "No config files found")
		return database.DatabaseData{}, fmt.Errorf("no config files found")
	}

	var filesOptions prompts.Options
	for i, file := range files {
		filesOptions = append(filesOptions, prompts.Option{Id: i, Label: file.Name()})
	}
	filesOptions = append(filesOptions, prompts.Option{Id: -1, Label: "Go back"})

	fileDbSelected, err := prompts.SelectPrompt("Choose your db config:", filesOptions)
	if err != nil || fileDbSelected == -1 {
		return database.DatabaseData{}, fmt.Errorf("go back")
	}

	configDbJson := files[fileDbSelected]
	file, err := os.Open(dbConfigFilesDir + configDbJson.Name())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err.Error())
		return database.DatabaseData{}, err
	}
	defer file.Close()

	var dbConfig database.DatabaseData
	if err := json.NewDecoder(file).Decode(&dbConfig); err != nil {
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
		fmt.Fprintln(os.Stderr, "No dump files found")
		return "", fmt.Errorf("no dump files found")
	}

	var filesFromOptions prompts.Options
	for i, file := range filesFrom {
		filesFromOptions = append(filesFromOptions, prompts.Option{Id: i, Label: file.Name()})
	}
	filesFromOptions = append(filesFromOptions, prompts.Option{Id: -1, Label: "Go back"})

	fileSelected, err := prompts.SelectPrompt("Choose your sql dump:", filesFromOptions)
	if err != nil || fileSelected == -1 {
		return "", fmt.Errorf("go back")
	}

	return filesFrom[fileSelected].Name(), nil
}

func ChangeFilePaths() {
	fmt.Fprintln(os.Stdout, "Describe the path like '/tmp/db/'")
	dbConfigPath := prompts.StringPrompt("Path to database configuration files:", false)
	dbDumpPath := prompts.StringPrompt("Path to database dump files:", false)

	homeDir, _ := os.UserHomeDir()
	err := os.MkdirAll(filepath.Join(homeDir, ".go-cli-db"), 0755)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err.Error())
		return
	}

	content := fmt.Sprintf("DB_BACKUP_DIR=%s\nDB_FILES_DIR=%s\n", dbConfigPath, dbDumpPath)
	err = os.WriteFile(filepath.Join(homeDir, ".go-cli-db", ".env"), []byte(content), 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err.Error())
		return
	}
}
