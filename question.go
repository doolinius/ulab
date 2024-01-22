package ulab

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

/*
The Question object. Uses the same object for different types of questions.
*/
type Question struct {
	Type     string   `json:"type"`
	Text     string   `json:"text"`
	Options  []string `json:"options"`
	Correct  []int    `json:"correct"`
	Feedback string   `json:"feedback"`
}

func (q *Question) Ask() {
	// Print question prompt
	// Print options if necessary
	// Prompt for answer
	// Validate input
	// Check answer
	fmt.Println("			** QUESTION **")
	fmt.Printf("%s\n\n", q.Text)
	switch q.Type {
	case "TF":
		fmt.Printf("\ttrue or false?\n\n")
		answer := getTF()
		fmt.Printf("%s\n", answer)
		//q.CheckAnswer(answer)
	case "MC", "MS", "ORD":
		for i, option := range q.Options {
			fmt.Printf("\t%d - %s\n", i+1, option)
		}
	case "NUM":
	}
}

func getTF() string {
	var answer string
	fmt.Printf("	=> ")
	fmt.Scan(&answer)
	answer = strings.ToLower(answer)
	for answer != "true" && answer != "false" {
		fmt.Printf("Invalid response. Enter 'true' or 'false'")
		fmt.Printf("	=>  ")
		fmt.Scan(&answer)
		answer = strings.ToLower(answer)
	}
	return answer
}

func (q *Question) getAnswer() {
	switch q.Type {
	case "TF":

	case "MC":
		var answer int
		fmt.Printf("Enter the number of the correct response (1-%d):", len(q.Options))
		fmt.Printf("	=> ")
		fmt.Scan(&answer)
		for answer > 1 && answer > len(q.Options) {
			fmt.Printf("Invalid response. Enter a number between 1 and %d\n", len(q.Options))
			fmt.Printf("	=>  ")
			fmt.Scan(&answer)
		}
	case "SA":
		var answer string
		fmt.Printf("Type a short answer (will not be auto-graded)\n")
		fmt.Printf("	=> ")
		fmt.Scan(&answer)
	case "NUM":
		stdin := bufio.NewReader(os.Stdin)
		var answer int
		var err error
		fmt.Printf("Enter the answer as a number:\n")
		fmt.Printf("	=> ")
		_, err = fmt.Scanf("%d", &answer)
		for err != nil {
			stdin.ReadString('\n')
			fmt.Printf("Invalid input. Please enter a whole number.")
			fmt.Printf("	=> ")
			_, err = fmt.Scanf("%d", &answer)
		}
	}
}
