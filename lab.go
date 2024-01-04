package ulab

import (
	"fmt"
	"os"
	"os/exec"
)

type StepTest struct {
	TestType  string `json:"testType"`
	Command   string `json:"command"`
	Condition bool   `json:"condition"` // do we test for true/pass or false/fail?
}

type Step struct {
	Text           string   `json:"text"`
	Tasks          []string `json:"tasks"`
	Tips           string   `json:"tips"`
	Test           StepTest `json:"test"`
	SuccessMessage string   `json:"successMessage"`
	RetryMessage   string   `json:"retryMessage"`
}

type Lab struct {
	Number      string `json:"number"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Steps       []Step `json:"steps"`
	Flags       []int  `json:"flags"`
	BonusFlags  []int  `json:"bonusFlags"`
}

func (l *Lab) printSteps(s *Status) {
	results := s.getResults(l.Number)
	fmt.Printf("Lab %s Steps\n\n", l.Number)
	fmt.Printf("   Flags: %d	Bonus Flags: %d\n", len(l.Flags), len(l.BonusFlags))
	fmt.Printf("   Steps:\n")
	for i, step := range l.Steps {
		stepStatus := "incomplete"
		if len(results.Steps) < i {
			stepStatus = results.Steps[i]
		}
		fmt.Printf("      %d. %s (%s)\n", i+1, step.Text, stepStatus)
	}
}

func (l *Lab) check(step int) bool {
	test := l.Steps[step].Test
	switch test.TestType {
	case "script":
		// execute check command
		_, err := exec.Command("sh", test.Command).Output()
		// fun new switch statement
		switch {
		case (err != nil && test.Condition == true), (err == nil && test.Condition == false):
			fmt.Printf("%v\n", err)
			fmt.Printf("\n%s\n\n", l.Steps[step].RetryMessage)
			return false
		case (err == nil && test.Condition == true), (err != nil && test.Condition == false):
			fmt.Printf("\n%s\n\n", l.Steps[step].SuccessMessage)
			return true
		}
	case "command":
		// execute check command
		_, err := exec.Command("sh", "-c", test.Command).Output()
		if err != nil {
			fmt.Printf("%v\n", err)
			fmt.Printf("\n%s\n\n", l.Steps[step].RetryMessage)
			return false
		} else {
			fmt.Printf("\n%s\n\n", l.Steps[step].SuccessMessage)
			return true
		}
	default:
		fmt.Printf("Not a recognized test type.")
		os.Exit(1)
		return false
	}
	return false
}

func (l *Lab) printStep(stepNum int) {
	l.Steps[stepNum].printTasks(stepNum)
}

func (s *Step) printTasks(stepNum int) {
	fmt.Printf("Step %d: %s\n\n", stepNum, s.Text)
	fmt.Println("Perform the following tasks/commands:")
	for _, task := range s.Tasks {
		fmt.Printf("\t%s\n", task)
	}
	fmt.Printf("\tlab check\n\n")
	if s.Tips != "" {
		fmt.Printf("TIPS: %s\n", s.Tips)
	}
}

func (l *Lab) checkFlag(flagNum int) bool {
	for _, flag := range l.Flags {
		if flag == flagNum {
			return true
		}
	}
	return false
}
