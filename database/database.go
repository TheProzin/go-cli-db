package database

import (
	"fmt"
	"go-cli-db/globals"
	"log"
	"os"
	"os/exec"
	"time"
)

type DatabaseData struct {
	Password     string   `json:"db_password"`
	Username     string   `json:"db_username"`
	Host         string   `json:"db_host"`
	Name         string   `json:"db_name"`
	IgnoreTables []string `json:"ignore_tables"`
}

func InitDbDump(dbConfig DatabaseData) {
	time := time.Now().Format("20060102150405")
	outputFile := dbConfig.Name + "__t" + time + ".sql"
	args := []string{
		"--max_allowed_packet=512M",
		"--skip-lock-tables",
		"--skip-add-locks",
		"-h", dbConfig.Host,
		"-u", dbConfig.Username,
		"-p" + dbConfig.Password,
		"--routines",
		"--triggers",
		"--set-charset",
		"--complete-insert",
		"--single-transaction",
	}

	for _, table := range dbConfig.IgnoreTables {
		args = append(args, "--ignore-table="+dbConfig.Name+"."+table)
	}

	args = append(args, dbConfig.Name)

	cmd := exec.Command("mysqldump", args...)

	dbFilesDir := globals.GoDotEnvVariable("DB_BACKUP_DIR")
	if _, err := os.Stat(dbFilesDir); err != nil {
		err := os.Mkdir(dbFilesDir, 0750)
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}
	}

	outFile, err := os.Create(dbFilesDir + outputFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating file:", err)
		return
	}
	defer outFile.Close()

	cmd.Stdout = outFile
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
}
