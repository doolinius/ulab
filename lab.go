package ulab

import (
	"encoding/json"
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
	Datafiles   bool   `json:"dataFiles"`
	Steps       []Step `json:"steps"`
	Flags       []int  `json:"flags"`
	BonusFlags  []int  `json:"bonusFlags"`
}

func (l *Lab) Extract() {
	if !l.Datafiles {
		fmt.Printf("No data files are required for this lab.")
	} else {
		fmt.Printf("\nThis lab requires data files.\n")
		fmt.Printf("Extracting data files...")
		datafilePath := ULConfig.Root + "/labs/" + l.Number + "/data.zip"
		out, err := exec.Command("/usr/bin/unzip", datafilePath).Output()
		if err != nil {
			fmt.Printf("Error extracting data files: %v\n", err)
		}
	}
}

func (l *Lab) PrintSteps(s *Status) {
	results := s.GetResults(l.Number)
	fmt.Printf("Lab %s Steps\n\n", l.Number)
	fmt.Printf("   Flags: %d	Bonus Flags: %d\n", len(l.Flags), len(l.BonusFlags))
	fmt.Printf("   Steps:\n")
	for i, step := range l.Steps {
		stepStatus := "incomplete"
		if i < len(results.Steps) {
			stepStatus = results.Steps[i]
		}
		fmt.Printf("      %d. %s (%s)\n", i+1, step.Text, stepStatus)
	}
}

func (l *Lab) Check(step int) bool {
	test := l.Steps[step].Test
	switch test.TestType {
	case "script":
		// execute check command
		cmdPath := ULConfig.Root + "/labs/" + l.Number + "/scripts/" + test.Command
		fmt.Printf("Running script %s...\n", cmdPath)
		_, err := exec.Command("sh", cmdPath).Output()
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
		fmt.Printf("Running command %s...\n", test.Command)
		_, err := exec.Command("sh", "-c", test.Command).Output()
		switch {
		case (err != nil && test.Condition == true), (err == nil && test.Condition == false):
			fmt.Printf("%v\n", err)
			fmt.Printf("\n%s\n\n", l.Steps[step].RetryMessage)
			return false
		case (err == nil && test.Condition == true), (err != nil && test.Condition == false):
			fmt.Printf("\n%s\n\n", l.Steps[step].SuccessMessage)
			return true
		}
		/*
		if err != nil {
			fmt.Printf("%v\n", err)
			fmt.Printf("\n%s\n\n", l.Steps[step].RetryMessage)
			return false
		} else {
			fmt.Printf("\n%s\n\n", l.Steps[step].SuccessMessage)
			return true
		}
		*/
	default:
		fmt.Printf("Not a recognized test type.")
		os.Exit(1)
		return false
	}
	return false
}

func (l *Lab) PrintStep(stepNum int) {
	l.Steps[stepNum].PrintTasks(stepNum)
}

func (s *Step) PrintTasks(stepNum int) {
	fmt.Printf("Step %d: %s\n\n", stepNum+1, s.Text)
	fmt.Println("Perform the following tasks/commands:")
	for _, task := range s.Tasks {
		fmt.Printf("\t%s\n", task)
	}
	fmt.Printf("\tlab check\n\n")
	if s.Tips != "" {
		fmt.Printf("TIPS: %s\n", s.Tips)
	}
}

func (l *Lab) CheckFlag(flagNum int) bool {
	for _, flag := range l.Flags {
		if flag == flagNum {
			return true
		}
	}
	return false
}

func (l *Lab) CheckBonusFlag(flagNum int) bool {
	for _, flag := range l.BonusFlags {
		if flag == flagNum {
			return true
		}
	}
	return false
}

func OpenLabFile(labNum string) *Lab {
	filepath := ULConfig.Root + "/labs/" + labNum + "/lab.json"
	//fmt.Printf("Opening lab file %s\n", filepath)
	labJson, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Could not read data file for Lab %s: ", labNum)
		os.Exit(1)
		return nil
	} else {
		lab := Lab{}
		json.Unmarshal(labJson, &lab)
		return &lab
	}
}
