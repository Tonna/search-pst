package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	//read input dir name and name of output file

	var input string
	var output string
	flag.StringVar(&input, "input", "", "directory that contains pst files")
	flag.StringVar(&output, "output", "", "file that contains all emails and attachments, use sql to query")

	flag.Parse()

	fmt.Println("input folder= " + input)
	fmt.Println("output file= " + output)

	if len(input) == 0 || len(output) == 0 {
		fmt.Println("\ninvalid in/out specified... Exiting")
		os.Exit(1)
	}

	visit := func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if strings.Contains(info.Name(), "-") {
				fmt.Println(" ", p, "attachment")
			} else {
				fmt.Println(" ", p, "email")
			}
		}
		//fmt.Println(" ", p, info.IsDir(), info.Name())
		return nil
	}
	err := filepath.Walk(input, visit)
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
