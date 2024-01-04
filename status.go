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
	Number         string   `json:"number"`
	SubmissionCode string   `json:"submissionCode"`
	StartTime      string   `json:"startTime"`
	FinishTime     string   `json:"finishTime"`
	Status         string   `json:"status"`
	Steps          []string `json:"steps"`
	Flags          []int    `json:"flags"`
	BonusFlags     []int    `json:"bonusFlags"`
}

type Status struct {
	Username    string                `json:"username"`
	CurrentLab  string                `json:"currentLab"`
	CurrentStep int                   `json:"currentStep"`
	Results     map[string]*LabResult `json:"results"`
}

func (lr *LabResult) init() {
	var flagSlice []int
	var stepSlice []string
	for _, f := range lr.Flags {
		flagSlice = append(flagSlice, f)
	}
	for _, s := range lr.Steps {
		stepSlice = append(stepSlice, s)
	}
	lr.Flags = flagSlice
	lr.Steps = stepSlice
}

func (s *Status) init() {
	var resultsSlice []LabResult
	for _, lr := range s.Results {
		lr.init()
		resultsSlice = append(resultsSlice, *lr)
	}
}

func (s *Status) FlagStatus(l *Lab) (int, int) {
	result := s.getResults(s.CurrentLab)
	return len(result.Flags), len(l.Flags)
}

func (s *Status) StepStatus(l *Lab) (int, int) {
	result := s.getResults(s.CurrentLab)
	numSteps := 0
	for _, stepStatus := range result.Steps {
		if stepStatus == "success" {
			numSteps++
		}
	}
	return numSteps, len(l.Steps)
}

func (s *Status) ScoreReport(l *Lab) {
	score := 0
	result := s.getResults(s.CurrentLab)
	numFlags := len(result.Flags)
	numBonusFlags := len(result.BonusFlags)
	numSteps, numTotalSteps := s.StepStatus(l)

	score = numSteps + numFlags + numBonusFlags
	fmt.Printf("Steps completed: %d/%d\n", numSteps, numTotalSteps)
	fmt.Printf(" Flags captured: %d/%d\n", numFlags, len(l.Flags))
	fmt.Printf("    Bonus Flags: %d/%d\n", numBonusFlags, len(l.BonusFlags))
	fmt.Printf("    Total Score: %d/%d\n", score, len(l.Steps)+len(l.Flags))
}

func (s *Status) submissionCode(time string) string {
	text := s.Username + time
	h := fnv.New64()
	h.Write([]byte(text))
	sum := h.Sum64()
	return fmt.Sprintf("%d", sum)
}

func (s *Status) Submit(lab *Lab) {
	result := s.getResults(lab.Number)
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
	newResult.Steps = make([]string, len(lab.Steps))
	for i := 0; i < len(newResult.Steps); i++ {
		newResult.Steps[i] = "incomplete"
	}
	newResult.Flags = make([]int, 0)
	newResult.BonusFlags = make([]int, 0)
	s.Results[lab.Number] = &newResult
}

func (s *Status) getResults(labNum string) *LabResult {
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

func (s *Status) AddFlag(labNum string, flagNum int) {
	result := s.getResults(labNum)
	//fmt.Printf("Result: %v  Flagnum: %d\n", result, flagNum)
	if slices.Contains(result.Flags, flagNum) {
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
	result := s.getResults(labNum)
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

func ReadStatusFile(username string, s *Status) {
	jsonData, err := os.ReadFile(username + ".json")
	if err != nil {
		fmt.Printf("Could not open user file.")
		os.Exit(1)
	}
	json.Unmarshal(jsonData, s)
}
