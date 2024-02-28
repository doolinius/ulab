package ulab

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

type StepTest struct {
	TestType  string `json:"testType"`
	Command   string `json:"command"`
	Value     string `json:"value"`
	Condition bool   `json:"condition"` // do we test for true/pass or false/fail?
}

type Step struct {
	ShortText      string   `json:"shortText"`
	Text           string   `json:"text"`
	Tasks          []string `json:"tasks"`
	Tips           string   `json:"tips"`
	Test           StepTest `json:"test"`
	SuccessMessage string   `json:"successMessage"`
	RetryMessage   string   `json:"retryMessage"`
	Question       Question `json:"question"`
}

type Lab struct {
	Number        string `json:"number"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Datafiles     bool   `json:"dataFiles"`
	Steps         []Step `json:"steps"`
	SubmitMessage string `json:"submitMessage"`
	Flags         []int  `json:"flags"`
	BonusFlags    []int  `json:"bonusFlags"`
}

func (l *Lab) Extract() {
	if !l.Datafiles {
		//fmt.Printf("No data files are required for this lab.")
	} else {
		//fmt.Printf("\nThis lab requires data files.\n")
		//fmt.Printf("Extracting data files...")
		datafilePath := ULConfig.Root + "/labs/" + l.Number + "/data.zip"
		_, err := exec.Command("/usr/bin/unzip", datafilePath).Output()
		if err != nil {
			fmt.Printf("Error extracting data files: %v\n", err)
		}
	}
}

func (l *Lab) PrintSteps(s *Status) {
	results := s.GetResults(l.Number)
	pterm.DefaultHeader.WithFullWidth().Printf("Lab %s Status", l.Number)
	//fmt.Printf("Lab %s Steps\n\n", l.Number)
	pterm.DefaultCenter.Printf("Flags: %d/%d\nBonus Flags: %d/%d", len(results.Flags), len(l.Flags), len(results.BonusFlags), len(l.BonusFlags))
	//fmt.Printf("   Flags: %d	Bonus Flags: %d\n", len(l.Flags), len(l.BonusFlags))
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
	case "checkvar":
		val := os.Getenv(test.Command)
		if val == test.Value && test.Condition == true {
			pterm.Success.Println(l.Steps[step].SuccessMessage)
			//fmt.Printf("\n%s\n\n", l.Steps[step].SuccessMessage)
			return true
		} else {
			pterm.Error.Println(l.Steps[step].RetryMessage)
			//fmt.Printf("\n%s\n\n", l.Steps[step].RetryMessage)
			return false
		}
	case "script":
		// execute check command
		cmdPath := ULConfig.Root + "/labs/" + l.Number + "/scripts/" + test.Command
		//fmt.Printf("Running script %s...\n", cmdPath)
		_, err := exec.Command("sh", cmdPath).Output()
		// fun new switch statement
		switch {
		case (err != nil && test.Condition == true), (err == nil && test.Condition == false):
			//fmt.Printf("%v\n", err)
			pterm.Error.Println(l.Steps[step].RetryMessage)
			//fmt.Printf("\n%s\n\n", l.Steps[step].RetryMessage)
			return false
		case (err == nil && test.Condition == true), (err != nil && test.Condition == false):
			pterm.Success.Println(l.Steps[step].SuccessMessage)
			//fmt.Printf("\n%s\n\n", l.Steps[step].SuccessMessage)
			return true
		}
	case "command":
		// execute check command
		//fmt.Printf("Running command %s...\n", test.Command)
		_, err := exec.Command("bash", "-c", test.Command).Output()
		switch {
		case (err != nil && test.Condition == true), (err == nil && test.Condition == false):
			//fmt.Printf("%v\n", err)
			pterm.Error.Println(l.Steps[step].RetryMessage)
			//fmt.Printf("\n%s\n\n", l.Steps[step].RetryMessage)
			return false
		case (err == nil && test.Condition == true), (err != nil && test.Condition == false):
			pterm.Success.Println(l.Steps[step].SuccessMessage)
			//fmt.Printf("\n%s\n\n", l.Steps[step].SuccessMessage)
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
	primary := pterm.NewStyle(pterm.FgBlack, pterm.BgCyan)
	secondary := pterm.NewStyle(pterm.FgBlack, pterm.BgLightYellow)
	pterm.Println(primary.Sprintf("\n Step %d: ", stepNum+1) + pterm.Sprintf(" %s\n", s.Text))
	//fmt.Printf("Step %d: %s\n\n", stepNum+1, s.Text)
	fmt.Println("Perform the following tasks/commands:")
	bulletListItems := "" /*pterm.BulletListItem{
		{Level: 0, Text: "Level 0"}, // Level 0 item
		{Level: 1, Text: "Level 1"}, // Level 1 item
		{Level: 2, Text: "Level 2"}, // Level 2 item
	}*/
	for _, task := range s.Tasks {
		bulletListItems += "    " + task + "\n"
		//fmt.Printf("\t%s\n", task)
	}
	bulletListItems += "    lab check"
	putils.BulletListFromString(bulletListItems, " ").Render()
	//fmt.Printf("\tlab check\n\n")
	if s.Tips != "" {
		pterm.Println(secondary.Sprintf("  TIPS:  ") + pterm.Sprintf(" %s", s.Tips))
		//fmt.Printf("TIPS: %s\n", s.Tips)
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
