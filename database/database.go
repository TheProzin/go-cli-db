package database

import (
	"fmt"
	"go-cli-db/globals"
	"os"
	"os/exec"
	"time"
)

type DatabaseData struct {
	ConfigName   string   `json:"config_name"`
	DbUsername   string   `json:"db_username"`
	DbPassword   string   `json:"db_password"`
	DbHost       string   `json:"db_host"`
	DbName       string   `json:"db_name"`
	IgnoreTables []string `json:"ignore_tables"`
}

func InitDbDump(dbConfig DatabaseData) (string, error) {
	time := time.Now().Format("20060102150405")
	outputFile := dbConfig.DbName + "__t" + time + ".sql"
	args := []string{
		"--max_allowed_packet=512M",
		"--skip-lock-tables",
		"--skip-add-locks",
		"-h", dbConfig.DbHost,
		"-u", dbConfig.DbUsername,
		"-p" + dbConfig.DbPassword,
		"--routines",
		"--triggers",
		"--set-charset",
		"--complete-insert",
		"--single-transaction",
	}

	for _, table := range dbConfig.IgnoreTables {
		args = append(args, "--ignore-table="+dbConfig.DbName+"."+table)
	}

	args = append(args, dbConfig.DbName)

	cmd := exec.Command("mysqldump", args...)

	dbBkpFilesDir := globals.GoDotEnvVariable("DB_BACKUP_DIR")
	if _, err := os.Stat(dbBkpFilesDir); err != nil {
		err := os.Mkdir(dbBkpFilesDir, 0750)
		if err != nil && !os.IsExist(err) {
			fmt.Fprintln(os.Stderr, "Error creating file:", err)
			return "", err
		}
	}

	outFile, err := os.Create(dbBkpFilesDir + outputFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating file:", err)
		return "", err
	}
	defer outFile.Close()

	cmd.Stdout = outFile
	cmd.Stderr = os.Stderr

	fmt.Fprintln(os.Stdout, "Running mysqldump...")
	err = cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return "", err
	}

	fmt.Fprintln(os.Stdout, "Success on making backup")
	return outputFile, nil
}

func InitDbRestore(dumpFileName string, dbConfig DatabaseData) {
	args := []string{
		"--max_allowed_packet=1024M",
		"--quick",
		"-h", dbConfig.DbHost,
		"-u", dbConfig.DbUsername,
		"-p" + dbConfig.DbPassword,
		dbConfig.DbName,
	}

	cmd := exec.Command("mysql", args...)

	dbBkpFilesDir := globals.GoDotEnvVariable("DB_BACKUP_DIR")
	inFile, err := os.Open(dbBkpFilesDir + dumpFileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error opening dump file:", err)
		return
	}
	defer inFile.Close()

	cmd.Stdin = inFile
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Fprintln(os.Stdout, "Running mysql restore...")
	err = cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	} else {
		fmt.Fprintln(os.Stdout, "Success on restore")
	}
}
