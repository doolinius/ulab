package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("Welcome to gradelab")

	if len(os.Args) < 1 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {

	}
}

func printUsage() {

}
