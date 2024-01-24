package ulab

import (
	"encoding/json"
	"fmt"
	"testing"
)

var qJSON = `
{
	"type":"TF",
	"text":"Wendy is beautiful",
	"options":[],
	"correct":[1],
	"feedback":"Yes, Wendy is indeed beautiful."
}`

func TestQuestionJSON(t *testing.T) {
	var q Question
	json.Unmarshal([]byte(qJSON), &q)
	if q.Type != "TF" {
		t.Errorf("Did not read JSON")
	}
}

func TestCheck(t *testing.T) {
	var q Question
	json.Unmarshal([]byte(qJSON), &q)
	if q.Check(1) {
		fmt.Printf("Correct! %s\n", q.Feedback)
	} else {
		fmt.Printf("Incorrect!\n")
	}
}
