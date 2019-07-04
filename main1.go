package main

import (
	"fmt"
	//"os"
	"flag"
)

func main() {
	var input string
	var output string
	flag.StringVar(&input, "input", "", "directory that contains pst files")
	flag.StringVar(&output, "output", "", "file that contains all emails and attachments, use sql to query")

	flag.Parse()

	fmt.Println("input folder= " + input)
	fmt.Println("output file= " + output)

	//var toWald Directory
	//take top dir from set, look what child dirs it has and put to set,
	//just files should be parsed and put to db
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
