package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	//read input dir name and name of output file

	var input string
	var output string

	{
		flag.StringVar(&input, "input", "", "directory that contains pst files")
		flag.StringVar(&output, "output", "", "file that contains all emails and attachments, use sql to query")

		flag.Parse()

		fmt.Println("input folder= " + input)
		fmt.Println("output file= " + output)

		if len(input) == 0 || len(output) == 0 {
			fmt.Println("\ninvalid in/out specified... Exiting")
			os.Exit(1)
		}
	}

	os.Remove(output)
	db, err := sql.Open("sqlite3", output)
	check(err)
	defer db.Close()

	sqlStmt := `create table email (path text, content text);
	            delete from email;

		    create table attachment (path text, content blob);
		    delete from attachment;
		 `

	_, err = db.Exec(sqlStmt)
	check(err)

	tx, err := db.Begin()
	check(err)

	stmtEmail, err := tx.Prepare("insert into email(path, content) values(?, ?)")
	check(err)

	stmtAtt, err := tx.Prepare("insert into attachment(path, content) values(?, ?)")
	check(err)

	defer stmtEmail.Close()
	defer stmtAtt.Close()

	//TODO file path shouldn't contain reduntant info. No "/home/..."
	//TODO all logic goes to visit closure. Is it ok?
	//TODO Add logic to parse email file and create more columns in database

	visit := func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {

			//this is naive check, but is acceptable for now
			if strings.Contains(info.Name(), "-") {

				fmt.Println(" ", p, "attachment")

				file, err := ioutil.ReadFile(p)
				check(err)

				_, err = stmtAtt.Exec(p, file)
				check(err)

			} else {

				fmt.Println(" ", p, "email")

				file, err := ioutil.ReadFile(p)
				check(err)

				_, err = stmtEmail.Exec(p, string(file))
				check(err)

			}
		}

		return nil
	}

	err = filepath.Walk(input, visit)
	check(err)

	tx.Commit()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
