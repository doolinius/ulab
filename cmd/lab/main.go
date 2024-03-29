package main

import (
	"fmt"
	"os"
	"os/user"
	"strconv"

	"github.com/doolinius/ulab"
	"github.com/pterm/pterm"
)

func main() {

	// Get user info and current lab status
	//var userStatus ulab.Status
	user, err := user.Current()
	if err != nil {
		pterm.Error.Printf("User information not found.\n")
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
	case "restart":
		inProgressCheck(userStatus)
		// make sure user in HOME directory
		if os.Getenv("PWD") != os.Getenv("HOME") {
			pterm.Warning.Println("Labs must be restarted in your HOME directory. Please 'cd' to your home directory and try again.")
			//fmt.Printf("Labs must be restarted in your HOME directory. Please 'cd' to your home directory and try again.\n\n")
			os.Exit(1)
		}
		// Are you sure prompt
		pterm.Warning.Printf("By restarting this lab, you will erase all current progress, including all steps and flags, and begin from the start.\n\n")
		answer, _ := pterm.DefaultInteractiveConfirm.Show()
		if !answer {
			os.Exit(1)
		} //else {
		//fmt.Printf("Restarting lab %s...", userStatus.CurrentLab)
		//}

		// delete lab data files
		_, err := os.Stat("lab_data")
		if !os.IsNotExist(err) {
			rmerr := os.RemoveAll("lab_data/")
			if rmerr != nil {
				fmt.Printf("Error removing lab data files: %v\n", rmerr)
			}
			//fmt.Printf("Lab data folder does not exist.")
			//os.Exit(1)
		}

		// start lab
		labStart(userStatus.CurrentLab, user, userStatus)
	case "start":
		if len(os.Args) != 3 {
			pterm.Error.Printf("A lab number must be supplied when starting a new lab.\n")
			printUsage()
			os.Exit(1)
		}
		//fmt.Printf("Starting Lab %s\n", os.Args[2])
		if userStatus.LabComplete(os.Args[2]) {
			pterm.Warning.Printf("This lab has already been submitted. If you start this lab again, you will ERASE your previous submission and begin a new one.\n\n")
			//fmt.Printf("This lab has already been submitted. If you start this lab again, you will ERASE your previous submission and begin a new one.\n\n")
			//answer := yesOrNo("Are you sure you want to start a new attempt? ")
			answer, _ := pterm.DefaultInteractiveConfirm.Show("Are you sure you want to start a new attempt?")
			if !answer {
				os.Exit(1)
			} else {
				pterm.Info.Printf("Erasing prior attempt and starting lab...\n\n")
			}
		}

		// make sure user in HOME directory
		if os.Getenv("PWD") != os.Getenv("HOME") {
			pterm.Error.Println("Labs should always be started in your HOME directory. Please 'cd' to your home directory and try again.")
			//fmt.Printf("Labs should always be started in your HOME directory. Please 'cd' to your home directory and try again.\n\n")
			os.Exit(1)
		}

		curLabNum, labInProgress := userStatus.InProgress()
		if labInProgress {
			//answer := ""
			pterm.Warning.Printf("Lab %s appears to be in progress.\nYou must submit this lab before starting a new one.\n", curLabNum)
			//fmt.Printf("Lab %s appears to be in progress. ", curLabNum)
			// printLabStatus(user)
			//fmt.Printf("\nYou must submit this lab before starting a new one.\n")
			submitPrompt := fmt.Sprintf("Would you like to submit Lab %s before starting the new one?", curLabNum)
			//answer = yesOrNo(submitPrompt)
			answer, _ := pterm.DefaultInteractiveConfirm.Show(submitPrompt)
			if answer {
				lab := ulab.OpenLabFile(curLabNum)
				labSubmit(userStatus, lab)
			} else {
				os.Exit(1)
			}
		}

		labStart(os.Args[2], user, userStatus)

	case "status":
		// TODO: necessary checks
		inProgressCheck(userStatus)
		//pwdCheck(userStatus)
		lab := ulab.OpenLabFile(userStatus.CurrentLab)
		lab.PrintSteps(userStatus)
	case "current":
		// TODO: necessary checks
		inProgressCheck(userStatus)
		pwdCheck(userStatus)
		lab := ulab.OpenLabFile(userStatus.CurrentLab)
		lab.PrintStep(userStatus.CurrentStep)
	case "check":
		// TODO: necessary checks
		inProgressCheck(userStatus)
		//pwdCheck(userStatus)
		//fmt.Printf("Checking current step...\n")
		pterm.Println()
		labCheck(userStatus)
	case "next":
		// do necessary checks
		inProgressCheck(userStatus)
		pwdCheck(userStatus)
		lab := ulab.OpenLabFile(userStatus.CurrentLab)
		labNext(lab, userStatus)
	case "flag":
		// TODO: Check arg numbers
		// TODO: necessary checks
		inProgressCheck(userStatus)
		pwdCheck(userStatus)
		if len(os.Args) != 3 {
			pterm.Error.Println("A flag number must be supplied when capturing a flag.")
			//fmt.Printf("A flag number must be supplied when capturing a flag.\n")
			//printUsage()
		} else {
			flagNum, err := strconv.Atoi(os.Args[2])
			if err != nil {
				pterm.Error.Println("The flag must be a valid four digit number.")
				//fmt.Printf("Flag must be a valid number")
			}
			lab := ulab.OpenLabFile(userStatus.CurrentLab)
			if lab.CheckFlag(flagNum) {
				if userStatus.AddFlag(lab.Number, flagNum) {
					pterm.Success.Printf("You've captured Lab %s Flag %d\n", lab.Number, flagNum)
					//fmt.Printf("SUCCESS! You've captured Lab %s Flag %d\n", lab.Number, flagNum)
				} else {
					pterm.Error.Printf("Flag %d has already been captured\n", flagNum)
				}
			} else if lab.CheckBonusFlag(flagNum) {
				if userStatus.AddBonusFlag(lab.Number, flagNum) {
					pterm.Success.Printf("You've captured Lab %s BONUS Flag %d\n", lab.Number, flagNum)
					//fmt.Printf("SUCCESS! You've captured Lab %s BONUS Flag %d\n", lab.Number, flagNum)
				} else {
					pterm.Error.Printf("Bonus Flag %d has already been captured\n", flagNum)
				}
			} else {
				pterm.Error.Printf("Invalid flag number (%d) for Lab %s!\n", flagNum, lab.Number)
				//fmt.Printf("Invalid flag number (%d) for Lab %s!\n", flagNum, lab.Number)
			}
		}
	case "submit":
		// TODO: necessary checks
		inProgressCheck(userStatus)
		lab := ulab.OpenLabFile(userStatus.CurrentLab)
		labSubmit(userStatus, lab)
	case "score":
		if len(os.Args) != 3 {
			pterm.Error.Println("A lab number must be supplied to view the score of a lab.")
			//fmt.Printf("A lab number must be supplied when starting a new lab.\n")
			printUsage()
		} else {
			// if the user has completed the lab
			if userStatus.LabComplete(os.Args[2]) {
				// show the report
				//fmt.Printf("Score Report for Lab %s\n\n", os.Args[2])
				userStatus.ScoreReport(os.Args[2])
			} else {
				// TODO: Check to see if it is a valid lab at all
				pterm.Warning.Printf("Lab %s has not been submitted yet.\n", os.Args[2])
				//fmt.Printf("Lab %s has not been submitted yet.\n", os.Args[2])
				os.Exit(1)
			}
		}
	case "help":
		printUsage()
	default:
		fmt.Printf("Unrecognized sub-command '%s'.\n", os.Args[2])
		printUsage()
	}
}

func labCheck(s *ulab.Status) {
	l := ulab.OpenLabFile(s.CurrentLab)
	// check status of current step
	if l.Check(s.CurrentStep) {
		s.Complete(l.Number, s.CurrentStep)
		s.SetPWD()
		s.Save()
		// If there is a Question to ask
		q := l.Steps[s.CurrentStep].Question
		if q.Type != "" {
			qNum := fmt.Sprintf("q%d", s.CurrentStep)
			if q.Ask() {
				pterm.Success.Println(q.Feedback)
				//fmt.Printf("Correct! %s\n", q.Feedback)
				s.AddQuestionResult(l.Number, qNum, "correct")
			} else {
				pterm.Error.Println("Sorry, that is incorrect")
				//fmt.Printf("Sorry, that is incorrect.\n")
				s.AddQuestionResult(l.Number, qNum, "incorrect")
			}
			s.Save()
		}
		// only prompt to move to the next step if the steps are not
		// completed
		if !s.StepsCompleted(l.Number) {
			//answer := yesOrNo("Would you like to move on to the next step?")
			//if answer == "yes" {
			//fmt.Printf("Let's move on to the next Step!\n\n")
			labNext(l, s)
			//}
		} else {
			pterm.Success.Printf("You have completed the steps for this lab. You may now submit it for grading with 'lab submit'\n")
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
    status         Lists the steps of current lab, as well as completion status
    current        Shows the details of the current step in the lab
    check          Checks the success of the current step in a lab
    flag <flag #>  Captures a numeric flag discovered in the lab
    submit         Submits the lab for grading
    restart        Restarts a lab in progress
    score <lab id> Show the final score of a submitted lab
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
	message := ""
	if success == stepsTotal {
		//fmt.Printf("Great job, you've finished all of the steps for this lab!\n")
	} else {
		message += fmt.Sprintf("You have %d steps left to complete.\n", stepsTotal-success)
		done = false
	}
	flagsCaptured := len(result.Flags)
	flagsTotal := result.TotalFlags
	if flagsCaptured == flagsTotal {
		//fmt.Printf("Nicely done, you've captured all of the flags for this lab!\n")
	} else {
		message += fmt.Sprintf("You have %d flags left to capture.\n", flagsTotal-flagsCaptured)
		done = false
	}
	if !done {
		pterm.Warning.Println(message, "This lab is not completed. Would you still like to submit?")
		// prompt "Are you sure?"
		//answer := yesOrNo("This lab is not completed. Would you still like to submit?")
		answer, _ := pterm.DefaultInteractiveConfirm.Show()
		if !answer {
			os.Exit(1)
		}
	}
	s.Submit(l)
	//s.ScoreReport(l.Number)
	pterm.DefaultSection.Printf("%s\n", l.SubmitMessage)
	//fmt.Printf("%s\n", l.SubmitMessage)
	// copy command history file
}

/*
func yesOrNo(prompt string) string {
	fmt.Printf("%s ", prompt)
	answer := ""
	for answer != "yes" && answer != "no" {
		fmt.Printf("yes/no =>  ")
		fmt.Scan(&answer)
		answer = strings.ToLower(answer)
	}
	return (answer)
}*/

func labStart(labNum string, u *user.User, s *ulab.Status) {
	//fmt.Printf("Attempting to start Lab %s\n", labNum)
	// check to see if lab number exists

	// Go on with starting lab
	// Open lab file

	lab := ulab.OpenLabFile(labNum)

	// Mark first step as in-progress for user
	s.CurrentLab = lab.Number
	s.CurrentStep = 0
	// Add new LabResult field to Status
	s.NewLab(lab)
	s.SetPWD()
	s.Save()

	// Print greeting
	primary := pterm.NewStyle(pterm.FgBlack, pterm.BgCyan)
	pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgCyan)).WithTextStyle(pterm.NewStyle(pterm.FgBlack)).Printf("Welcome to Lab %s - %s", lab.Number, lab.Name)
	//fmt.Printf("\nWelcome to Lab %s - %s\n\n", lab.Number, lab.Name)
	//pterm.Println()

	// check for data files and extract if necessary
	lab.Extract()

	descPar := pterm.DefaultSection.Sprint(lab.Description)
	//descParCenter := pterm.DefaultBox.WithTitle("Description").Sprint(descPar)
	pterm.DefaultCenter.Println(descPar)
	//fmt.Printf("%s\n\n", lab.Description)

	// Print flag info
	pterm.DefaultCenter.Println("This lab has " + primary.Sprintf("%d Flags", len(lab.Flags)) + " and " + primary.Sprintf("%d Bonus Flags", len(lab.BonusFlags)))
	//fmt.Printf("\tThis lab has %d Flags and %d Bonus Flags.\n\n", len(lab.Flags), len(lab.BonusFlags))

	// Get first step
	firstStep := lab.Steps[0]
	firstStep.PrintTasks(0)

}

/*
func printLabStatus(s *ulab.Status) {
	fmt.Printf("Printing lab status")
}*/

func pwdCheck(s *ulab.Status) {
	if os.Getenv("PWD") != s.LastPWD {
		pterm.Warning.Printf("You are not in the correct directory to continue this lab\nPlease 'cd' to %s\n", s.LastPWD)
		//fmt.Printf("You are not in the correct directory to continue this lab\n\n")
		//fmt.Printf("Please 'cd' to %s\n", s.LastPWD)
		os.Exit(2)
	}
}

func inProgressCheck(s *ulab.Status) {
	if s.CurrentLab == "" {
		pterm.Warning.Println("There is not currently a lab in progress. Start a lab with:\n\tlab start <lab number>")
		//fmt.Printf("There is not currently a lab in progress. Start a lab with:\n")
		//fmt.Printf("\n\tlab start <lab number>\n\n")
		os.Exit(1)
	}
}
