# go-cli-db

I made this program to practice cli in Golang and to make the daily backups and restores more easy. 
Now it only supports mysql/mariadb.

In the current state it has 4 functions:

1. Create new backup config
2. Start backup
3. Restore dump
4. Change files path


## 1 - Create new backup config
When selecting this options you are going to to input:

    Configuration name
    Database host
    Database name
    Database username
    Database password
    Database ignore tables

Configuration name: Is the name we are going to use to name the .json file

Database host, name, username, password: Is the config for your database.

Ignore tables: You can put the name of tables you dont want to make the data dump, separated by comma like: user_log,product_log

After confirm, you can start de dump with this new config or go back to menu.

## 2 - Start backup

When selecting this option, you are going to select from previous config files to start a new dump
Its goin to appear a list with the name of the configs your previous saved like:

    [] mydb.json

After you confirm, its going to make the dump on the dump files dir, and asks if you want to restore in another db config file

## 3 - Restore dump

When selecting this option, you are going to select options from the dump files dir. Its going to apear a list with the name of the dump files like:
    
    [] mydb_t20260118.sql

After selecting a dump file, you are going to select the config db file like in "Start backup" option.
Its going to appear a confimation screen then its going to make the restore.

## 4 - Change files path

When selecting this options, you are going to input the files path for the config and backup.

When starting the program for the first time, it will be asked for two folders
The first is where we are going to save the .json files of database config
The second is where we are going to save the dump files

It will create the file in the users folder as in:
    
    ~/.go-cli.db/.env

With this variables:

    DB_BACKUP_DIR=
    DB_FILES_DIR=