package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/doolinius/ulab"
)

func main() {

	c := ulab.ReadConfigFile()

	if len(os.Args) < 2 || len(os.Args) > 3 {
		printUsage()
		os.Exit(1)
	}

	labNumArg := flag.String("n", "", "the lab number (e.g. '3-1')")
	usernameArg := flag.String("u", "", "a username")
	submissionArg := flag.String("s", "", "the submission code for a particular assignment submission")

	flag.Parse()

	switch os.Args[1] {
	case "-n":
		fmt.Printf("Lab Number: %s\n", *labNumArg)
		files, err := os.ReadDir(c.Data)
		check(err)
		// for each file:
		for _, entry := range files {
			fullPath := c.Data + "/" + entry.Name()
			fmt.Printf("  %s\n", fullPath)
			userstatus := ulab.ReadStatusFile(fullPath)
			userstatus.FullResults()
		}
		//		open file
		//		print score for student
	case "-u":
		fmt.Printf("Username: %s\n", *usernameArg)
		// open student file
		// print report
	case "-s":
		fmt.Printf("Submission Code: %s\n", *submissionArg)
		//
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
