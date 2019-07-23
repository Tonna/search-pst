package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
	"time"
)

func main() {

	var input string
	var output string
	var resultName string = time.Now().Format(time.RFC3339)

	fmt.Println(resultName)
	{
		flag.StringVar(&input, "input", "", "path to sqlite database file")
		flag.StringVar(&output, "output", "", "folder where subfolder with search results will be placed")

		flag.Parse()

		fmt.Println("input folder= " + input)
		fmt.Println("output file= " + output)

		if len(input) == 0 || len(output) == 0 {
			fmt.Println("\ninvalid in/out specified... Exiting")
			os.Exit(1)
		}
	}

	path := output + "/" + resultName
	err := os.Mkdir(path, 0755)

	db, err := sql.Open("sqlite3", input)
	check(err)
	defer db.Close()

	query := "select path, content from email limit 1;"

	rows, err := db.Query(query)
	check(err)
	defer rows.Close()

	for rows.Next() {

		var name string
		var content string
		err = rows.Scan(&name, &content)
		check(err)
		fmt.Println("name is " + name)
		fmt.Println("content is " + content)

		err := ioutil.WriteFile(path+"/"+"1", []byte(content), 0644)
		check(err)
	}

	check(err)
}

func check(e error) {

	if e != nil {
		panic(e)
	}
}
