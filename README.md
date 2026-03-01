# go-cli-db

I made this program to practice cli in Golang and to make the daily backups and restores more easy. 
Now it only supports mysql/mariadb

In the current state it has 4 functions:

1. Create new backup config
2. Start backup
3. Restore dump
4. Change files path


## 1 - Create new backup config
When selecting this options you are going to to input

    Configuration name
    Database host
    Database name
    Database username
    Database password
    Database ignore tables

Configuration name: Is the name we are going to use to name the .json file

Database host, name, username, password: Is the config for your database

Ignore tables: You can put the name of tables you dont want to make the data dump, separated by comma like: user_log,product_log



When starting the program for the first time, it will be asked for two folders
The first is where we are going to save the .json files of database config
The second is where we are going to save the dump files

It will create the file in the users folder as in:
    
    ~/.go-cli.db/.env

With this variables:

    DB_BACKUP_DIR=
    DB_FILES_DIR=