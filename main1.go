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

		    create table attachment (path text, content blob, email_id integer);
		    delete from attachment;

		 `

	indexStmt := `	    CREATE INDEX index_attachment_email_id ON attachment (email_id );
			    CREATE INDEX index_attachment_path ON attachment (path );
			    CREATE INDEX index_email_path ON email (path );
		    `

	_, err = db.Exec(sqlStmt)
	check(err)

	tx, err := db.Begin()
	check(err)

	stmtEmail, err := tx.Prepare("insert into email(path, content) values(?, ?)")
	check(err)

	//setting email_id with latest email rowid we do get advantage
	//not running update query later
	//assuming file system will give us files in order - link right email and attachment
	stmtAtt, err := tx.Prepare("insert into attachment(path, content, email_id) values(?, ?, (select max(rowid) from email))")
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

			//very lazy way to create with primary key
			//by removing slashes and whitespaces from file path
			clearName := strings.Replace(p[len(input)+1:], "/", "_", -1)
			clearName = strings.Replace(clearName, " ", "", -1)

			//this is naive check, but is acceptable for now
			if strings.Contains(info.Name(), "-") {

				fmt.Println(" ", p, "attachment")

				file, err := ioutil.ReadFile(p)
				check(err)

				_, err = stmtAtt.Exec(clearName, file)
				check(err)

			} else {

				fmt.Println(" ", p, "email")

				file, err := ioutil.ReadFile(p)
				check(err)

				_, err = stmtEmail.Exec(clearName, string(file))
				check(err)

			}
		}

		return nil
	}

	err = filepath.Walk(input, visit)
	check(err)

	tx.Commit()

	_, err = db.Exec(indexStmt)
	check(err)

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
