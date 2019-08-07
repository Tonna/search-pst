package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func main() {

	var input string
	var output string
	var queryFile string
	var resultName string = time.Now().Format(time.RFC3339)

	fmt.Println(resultName)
	{
		flag.StringVar(&input, "input", "", "path to sqlite database file")
		flag.StringVar(&output, "output", "", "folder where subfolder with search results will be placed")
		flag.StringVar(&queryFile, "query", "", "path to file with sql query")
		flag.Parse()

		fmt.Println("input folder= " + input)
		fmt.Println("output file= " + output)
		fmt.Println("query file= " + queryFile)

		if len(input) == 0 || len(output) == 0 || len(queryFile) == 0 {
			fmt.Println("\ninvalid in/out specified... Exiting")
			os.Exit(1)
		}
	}

	path := output + "/" + resultName
	err := os.Mkdir(path, 0755)
	fmt.Println("created output file - " + path)

	//TODO assumption here that query will contain 'SELECT path, content'
	query, err := ioutil.ReadFile(queryFile)
	check(err)
	fmt.Println("query is - " + string(query))

	db, err := sql.Open("sqlite3", input)
	check(err)
	defer db.Close()

	rows, err := db.Query(string(query))
	check(err)
	defer rows.Close()

	for rows.Next() {

		var name string
		var content string
		err = rows.Scan(&name, &content)
		check(err)
		fmt.Println("name is " + name)
		fmt.Println("content is " + content)

		pathArray := strings.Split(name, "/")
		err := ioutil.WriteFile(path+"/"+pathArray[len(pathArray)-1], []byte(content), 0644)
		check(err)
	}

	check(err)
}

func check(e error) {

	if e != nil {
		panic(e)
	}
}
