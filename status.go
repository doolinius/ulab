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

func (s *Status) QuickScore(labNum string) {
	lr := s.GetResults(labNum)
	score, totalScore := lr.Score()
	scoreStr := fmt.Sprintf("%d/%d", score, totalScore)
	fmt.Printf("%16s: %7s    %s\n", s.Username, scoreStr, lr.SubmissionCode)
}

func (lr *LabResult) TotalScore() int {
	return lr.TotalSteps + lr.TotalFlags
}

func (lr *LabResult) Score() (int, int) {
	numSteps, _ := lr.StepStatus()
	score := numSteps + len(lr.Flags) + len(lr.BonusFlags)
	return score, lr.TotalScore()
}

func (lr *LabResult) StepStatus() (int, int) {
	//result := s.GetResults(s.CurrentLab)
	//l := OpenLabFile(lr.Number)
	numSteps := 0
	for _, stepStatus := range lr.Steps {
		if stepStatus == "success" {
			numSteps++
		}
	}
	return numSteps, lr.TotalSteps
}

func (s *Status) StepsCompleted(labNum string) bool {
	numSteps, totalSteps := s.Results[labNum].StepStatus()
	return numSteps == totalSteps
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
		score, total := result.Score()
		fmt.Printf("    %s: %d/%d\n", result.Number, score, total)
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
	result.Status = "completed"
	result.FinishTime = fmt.Sprintf("%d-%02d-%02dT%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute())
	result.SubmissionCode = s.submissionCode(result.FinishTime)
	fmt.Printf("\n\tSubmission Code: %s\n\n", result.SubmissionCode)
	s.CurrentLab = ""
	s.CurrentStep = -1
	s.Save()
}

func (s *Status) LabComplete(labNum string) bool {
	result := s.GetResults(labNum)
	if result != nil {
		return result.Status == "completed"
	}else{
		return false
	}
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

func (s *Status) AddFlag(labNum string, flagNum int) {
	result := s.GetResults(labNum)
	if slices.Contains(result.Flags, flagNum) {
		fmt.Printf("Flag %d has already been captured\n", flagNum)
		return
	}
	result.Flags = append(result.Flags, flagNum)
	s.Save()
}

func (s *Status) AddBonusFlag(labNum string, flagNum int) {
	result := s.GetResults(labNum)
	if slices.Contains(result.BonusFlags, flagNum) {
		fmt.Printf("Bonus Flag %d has already been captured\n", flagNum)
		return
	}
	result.BonusFlags = append(result.BonusFlags, flagNum)
	s.Save()
}

func (s *Status) FlagStatus(l *Lab) (int, int) {
	result := s.GetResults(s.CurrentLab)
	return len(result.Flags), len(l.Flags)
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
	//fmt.Printf("%v\n", result.Steps)
	//fmt.Printf("Lab %s Step %d, setting to \"success\"\n", labNum, stepNum)
	result.Steps[stepNum] = "success"
	//fmt.Printf("%v\n", result.Steps)
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
	f, err := os.Create(ULConfig.Data + "/" + s.Username + ".json")
	f.Write(json)
	//os.WriteFile(s.Username+".json", json, fs.)
}

func ReadStatusFile(fileName string) *Status {
	var s Status
	jsonData, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Could not open user status file: %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(jsonData, &s)
	return &s
}
