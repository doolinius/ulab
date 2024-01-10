package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/doolinius/ulab"
)

func main() {

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
		files, err := os.ReadDir(ulab.ULConfig.Data)
		check(err)
		// for each file:
		for _, entry := range files {
			fullPath := ulab.ULConfig.Data + "/" + entry.Name()
			//fmt.Printf("  %s\n", fullPath)
			userstatus := ulab.ReadStatusFile(fullPath)
			userstatus.QuickScore(*labNumArg)
		}
		//		open file
		//		print score for student
	case "-u":
		fmt.Printf("Username: %s\n", *usernameArg)
		// open student file
		fullPath := ulab.ULConfig.Data + "/" + *usernameArg + ".json"
		//fmt.Printf("  %s\n", fullPath)
		userstatus := ulab.ReadStatusFile(fullPath)
		// print report
		userstatus.FullResults()
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
