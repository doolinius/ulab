package main

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/doolinius/ulab"
)

func main() {

	// Get user info and current lab status
	//var userStatus ulab.Status
	user, err := user.Current()
	if err != nil {
		fmt.Printf("User information not found.\n")
		os.Exit(1)
	}
	//fmt.Println("Data File: ", ulab.ULConfig.Data+"/"+user.Username+".json")
	userStatus := ulab.ReadStatusFile(ulab.ULConfig.Data + "/" + user.Username + ".json")
	//userStatus.init()

	// Check for subcommand
	if len(os.Args) == 1 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "start":
		if len(os.Args) != 3 {
			fmt.Printf("A lab number must be supplied when starting a new lab.\n")
			printUsage()
		} else {
			//fmt.Printf("Starting Lab %s\n", os.Args[2])
			labStart(os.Args[2], user, userStatus)
		}
	case "steps":
		// TODO: necessary checks
		if userStatus.CurrentLab == "" {
			fmt.Printf("There is not currently a lab in progress. Start a lab with:\n")
			fmt.Printf("\n\tlab start <lab number>\n\n")
			os.Exit(1)
		}
		lab := ulab.OpenLabFile(userStatus.CurrentLab)
		lab.PrintSteps(userStatus)
	case "current":
		// TODO: necessary checks
		if userStatus.CurrentLab == "" {
			fmt.Printf("There is not currently a lab in progress. Start a lab with:\n")
			fmt.Printf("\n\tlab start <lab number>\n\n")
			os.Exit(1)
		}
		lab := ulab.OpenLabFile(userStatus.CurrentLab)
		lab.PrintStep(userStatus.CurrentStep)
	case "check":
		// TODO: necessary checks
		if userStatus.CurrentLab == "" {
			fmt.Printf("There is not currently a lab in progress. Start a lab with:\n")
			fmt.Printf("\n\tlab start <lab number>\n\n")
			os.Exit(1)
		}
		fmt.Printf("Checking current step...\n")
		labCheck(userStatus)
	case "next":
		// do necessary checks
		if userStatus.CurrentLab == "" {
			fmt.Printf("There is not currently a lab in progress. Start a lab with:\n")
			fmt.Printf("\n\tlab start <lab number>\n\n")
			os.Exit(1)
		}
		lab := ulab.OpenLabFile(userStatus.CurrentLab)
		labNext(lab, userStatus)
	case "flag":
		// TODO: Check arg numbers
		// TODO: necessary checks
		if userStatus.CurrentLab == "" {
			fmt.Printf("There is not currently a lab in progress. Start a lab with:\n")
			fmt.Printf("\n\tlab start <lab number>\n\n")
			os.Exit(1)
		}
		if len(os.Args) != 3 {
			fmt.Printf("A flag number must be supplied when capturing a flag.\n")
			//printUsage()
		} else {
			flagNum, err := strconv.Atoi(os.Args[2])
			if err != nil {
				fmt.Printf("Flag must be a valid number")
			}
			lab := ulab.OpenLabFile(userStatus.CurrentLab)
			if lab.CheckFlag(flagNum) {
				fmt.Printf("SUCCESS! You've captured Lab %s Flag %d\n", lab.Number, flagNum)
				userStatus.AddFlag(lab.Number, flagNum)
			} else if lab.CheckBonusFlag(flagNum) {
				fmt.Printf("SUCCESS! You've captured Lab %s BONUS Flag %d\n", lab.Number, flagNum)
				userStatus.AddFlag(lab.Number, flagNum, true)
			} else {
				fmt.Printf("Invalid flag number (%d) for Lab %s!\n", flagNum, lab.Number)
			}
		}
	case "submit":
		// TODO: necessary checks
		if userStatus.CurrentLab == "" {
			fmt.Printf("There is not currently a lab in progress. Start a lab with:\n")
			fmt.Printf("\n\tlab start <lab number>\n\n")
			os.Exit(1)
		}
		lab := ulab.OpenLabFile(userStatus.CurrentLab)
		labSubmit(userStatus, lab)
	case "help":
		printUsage()
	default:
		fmt.Printf("Unrecognized sub-command.\n")
		printUsage()
	}
}

func labCheck(s *ulab.Status) {
	l := ulab.OpenLabFile(s.CurrentLab)
	// check status of current step
	if l.Check(s.CurrentStep) {
		s.Complete(l.Number, s.CurrentStep)
		answer := yesOrNo("Would you like to move on to the next step?")
		if answer == "yes" {
			labNext(l, s)
		}
	}
}

func labNext(l *ulab.Lab, s *ulab.Status) {
	s.CurrentStep++
	s.Save()
	step := l.Steps[s.CurrentStep]
	step.PrintTasks(s.CurrentStep)
}

func printUsage() {
	out := `Usage: lab <subcommand> <argument>

Available subcommands:
    start <lab id> Begins a new lab (there can be only one active lab)
    steps          Lists the steps of the current lab
    current        Shows the details of the current step in the lab
    check          Checks the success of the current step in a lab
    next           Moves to the next step in the lab
    status         Shows your current progress in the lab
    flag <flag #>  Captures a numeric flag discovered in the lab
    submit         Submits the lab for grading
    help           Show this usage screen
`
	fmt.Println(out)
}

func labSubmit(s *ulab.Status, l *ulab.Lab) {
	// check for unfinished steps
	// print those
	done := true // if there any unfinished steps, or uncaptured flags
	result := s.GetResults(l.Number)
	success, stepsTotal := result.StepStatus()
	if success == stepsTotal {
		fmt.Printf("Great job, you've finished all of the steps for this lab!\n")
	} else {
		fmt.Printf("You have %d steps left to complete.\n", stepsTotal-success)
		done = false
	}
	flagsCaptured := len(result.Flags)
	flagsTotal := result.TotalFlags
	if flagsCaptured == flagsTotal {
		fmt.Printf("Nicely done, you've finished all of the steps for this lab!\n")
	} else {
		fmt.Printf("You have %d steps left to complete.\n", flagsTotal-flagsCaptured)
		done = false
	}
	if !done {
		// prompt "Are you sure?"
		answer := yesOrNo("This lab is not completed. Would you still like to submit?")
		if answer == "no" {
			os.Exit(1)
		}
	}
	s.Submit(l)
	// copy command history file
}

func yesOrNo(prompt string) string {
	fmt.Printf("%s ", prompt)
	answer := ""
	for answer != "yes" && answer != "no" {
		fmt.Printf("yes/no =>  ")
		fmt.Scan(&answer)
		answer = strings.ToLower(answer)
	}
	return (answer)
}

func labStart(labNum string, u *user.User, s *ulab.Status) {
	//fmt.Printf("Attempting to start Lab %s\n", labNum)
	// check to see if lab number exists
	curLabNum, labInProgress := s.InProgress()
	if labInProgress {
		answer := ""
		fmt.Printf("Lab %s appears to be in progress. ", curLabNum)
		// printLabStatus(user)
		fmt.Printf("\nYou must submit this lab before starting a new one.\n")
		submitPrompt := fmt.Sprintf("Would you like to submit Lab %s before starting the new one?", curLabNum)
		answer = yesOrNo(submitPrompt)
		if answer == "yes" {
			lab := ulab.OpenLabFile(curLabNum)
			labSubmit(s, lab)
		} else {
			os.Exit(1)
		}
	}

	// Go on with starting lab
	// Open lab file

	lab := ulab.OpenLabFile(labNum)

	// check for data files and extract if necessary
	lab.Extract()

	// Mark first step as in-progress for user
	s.CurrentLab = lab.Number
	s.CurrentStep = 0
	// Add new LabResult field to Status
	s.NewLab(lab)
	s.Save()

	// Print greeting
	fmt.Printf("\nWelcome to Lab %s - %s\n\n", lab.Number, lab.Name)
	fmt.Printf("%s\n\n", lab.Description)

	// Print flag info
	fmt.Printf("\tThis lab has %d Flags and %d Bonus Flags.\n\n", len(lab.Flags), len(lab.BonusFlags))

	// Get first step
	firstStep := lab.Steps[0]
	firstStep.PrintTasks(0)

}

func printLabStatus(s *ulab.Status) {
	fmt.Printf("Printing lab status")
}
