package ulab

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"os"
	"slices"
	"time"
)

type LabResult struct {
	Number          string   `json:"number"`
	SubmissionCode  string   `json:"submissionCode"`
	StartTime       string   `json:"startTime"`
	FinishTime      string   `json:"finishTime"`
	Status          string   `json:"status"`
	Steps           []string `json:"steps"`
	TotalSteps      int      `json:"totalSteps"`
	Flags           []int    `json:"flags"`
	TotalFlags      int      `json:"totalFlags"`
	BonusFlags      []int    `json:"bonusFlags"`
	TotalBonusFlags int      `json:"totalBonusFlags"`
}

type Status struct {
	Username    string                `json:"username"`
	CurrentLab  string                `json:"currentLab"`
	CurrentStep int                   `json:"currentStep"`
	Results     map[string]*LabResult `json:"results"`
}

func (lr *LabResult) QuickScore() {
	fmt.Printf("%s Score: %d/%d", lr.Number, lr.Score(), lr.TotalScore())
}

func (lr *LabResult) TotalScore() int {
	return lr.TotalSteps + lr.TotalFlags
}

func (lr *LabResult) Score() int {
	numSteps, _ := lr.StepStatus()
	score := numSteps + len(lr.Flags) + len(lr.BonusFlags)
	return score
}

func (s *Status) FlagStatus(l *Lab) (int, int) {
	result := s.GetResults(s.CurrentLab)
	return len(result.Flags), len(l.Flags)
}

func (lr *LabResult) StepStatus() (int, int) {
	//result := s.GetResults(s.CurrentLab)
	l := OpenLabFile(lr.Number)
	numSteps := 0
	for _, stepStatus := range lr.Steps {
		if stepStatus == "success" {
			numSteps++
		}
	}
	return numSteps, len(l.Steps)
}

func (s *Status) ScoreReport(l *Lab) {
	score := 0
	result := s.GetResults(s.CurrentLab)
	numFlags := len(result.Flags)
	numBonusFlags := len(result.BonusFlags)
	numSteps, numTotalSteps := result.StepStatus()

	score = numSteps + numFlags + numBonusFlags
	fmt.Printf("Steps completed: %d/%d\n", numSteps, numTotalSteps)
	fmt.Printf(" Flags captured: %d/%d\n", numFlags, result.TotalFlags)
	fmt.Printf("    Bonus Flags: %d/%d\n", numBonusFlags, result.TotalBonusFlags)
	fmt.Printf("    Total Score: %d/%d\n", score, numTotalSteps+result.TotalBonusFlags)
}

func (s *Status) FullResults() {
	for _, result := range s.Results {
		result.QuickScore()
	}
}

func (s *Status) submissionCode(time string) string {
	text := s.Username + time
	h := fnv.New64()
	h.Write([]byte(text))
	sum := h.Sum64()
	return fmt.Sprintf("%d", sum)
}

func (s *Status) Submit(lab *Lab) {
	result := s.GetResults(lab.Number)
	t := time.Now()
	result.FinishTime = fmt.Sprintf("%d-%02d-%02dT%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute())
	result.SubmissionCode = s.submissionCode(result.FinishTime)
	fmt.Printf("Submission Code: %s\n", result.SubmissionCode)
	s.CurrentLab = ""
	s.CurrentStep = -1
	s.Save()
}

func (s *Status) NewLab(lab *Lab) {
	var newResult = LabResult{}
	t := time.Now()
	newResult.Number = lab.Number
	newResult.SubmissionCode = ""
	newResult.StartTime = fmt.Sprintf("%d-%02d-%02dT%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute())
	newResult.Status = "inProgress"
	newResult.TotalSteps = len(lab.Steps)
	newResult.Steps = make([]string, newResult.TotalSteps)
	for i := 0; i < len(newResult.Steps); i++ {
		newResult.Steps[i] = "incomplete"
	}
	newResult.Flags = make([]int, 0)
	newResult.TotalFlags = len(lab.Flags)
	newResult.BonusFlags = make([]int, 0)
	newResult.TotalBonusFlags = len(lab.BonusFlags)
	s.Results[lab.Number] = &newResult
}

func (s *Status) GetResults(labNum string) *LabResult {
	_, keyExists := s.Results[labNum]
	if keyExists {
		return s.Results[labNum]
	}
	/*
		for i, result := range s.Results {
			if result.Number == labNum {
				return &s.Results[i]
			}
		}
	*/
	// if the labNum is not found in the results
	return nil
}

func (s *Status) AddFlag(labNum string, flagNum int, bOpt ...bool) {
	bonus := false
	if len(bOpt) > 1 {
		bonus = bOpt[0]
	}
	result := s.GetResults(labNum)
	flagList := result.Flags
	if bonus {
		flagList = result.BonusFlags
	}
	//fmt.Printf("Result: %v  Flagnum: %d\n", result, flagNum)
	if slices.Contains(flagList, flagNum) {
		fmt.Printf("Flag %d has already been captured\n", flagNum)
		return
	}
	result.Flags = append(result.Flags, flagNum)
	//fmt.Printf("Result Flags: %v\n", result.Flags)
	//fmt.Printf("Result Flags: %v\n", s.Results[1].Flags)
	s.Save()
}

func (s *Status) InProgress() (string, bool) {
	if s.CurrentLab == "" {
		return "", false
	} else {
		return s.CurrentLab, true
	}
}

func (s *Status) Complete(labNum string, stepNum int) {
	result := s.GetResults(labNum)
	result.Steps[stepNum] = "success"
	/*
		for _, result := range s.Results {
			if result.Number == labNum {
				result.Steps[stepNum] = "success"
				return
			}
		}
	*/
}

func (s *Status) Save() {
	json, err := json.MarshalIndent(s, "", "	")
	if err != nil {
		fmt.Printf("Error Marshaling Status Data")
		os.Exit(1)
	}
	f, err := os.Create(s.Username + ".json")
	f.Write(json)
	//os.WriteFile(s.Username+".json", json, fs.)
}

func ReadStatusFile(fileName string) *Status {
	var s Status
	jsonData, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Could not open user status file.")
		os.Exit(1)
	}
	json.Unmarshal(jsonData, &s)
	return &s
}
