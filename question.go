package ulab

import (
	"strconv"

	"github.com/pterm/pterm"
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

func (q *Question) Ask() bool {
	var answer int
	// Print question prompt
	// Print options if necessary
	// Prompt for answer
	// Validate input
	// Check answer
	pterm.DefaultCenter.Println(pterm.LightCyan("** QUESTION **"))
	//fmt.Printf("			** QUESTION **\n\n")
	//pterm.Printf("%s %s", pterm.LightYellow("-> "), q.Text)
	switch q.Type {
	case "TF":
		answer = q.getTF()
		//q.Check(answer)
	case "MC":
		answer = q.getOption()
	case "MS":
	case "NUM":
		answer = q.getNum()
	}
	return q.Check(answer)
}

func (q *Question) getTF() int {

	/*
		fmt.Printf("\ttrue or false?\n\n")
		fmt.Printf("	=> ")
		fmt.Scan(&answer)
		answer = strings.ToLower(answer)
		for answer != "true" && answer != "false" {
			fmt.Printf("Invalid response. Enter 'true' or 'false'\n\n")
			fmt.Printf("	=> ")
			fmt.Scan(&answer)
			answer = strings.ToLower(answer)
		}
	*/

	var options []string
	options = append(options, "True")
	options = append(options, "False")

	answer, _ := pterm.DefaultInteractiveSelect.WithFilter(false).WithOptions(options).Show(q.Text)

	if answer == "True" {
		return 1
	} else {
		return 0
	}
}

func (q *Question) getOption() int {
	selectedOption, _ := pterm.DefaultInteractiveSelect.WithFilter(false).WithOptions(q.Options).Show(q.Text)

	return getInt(selectedOption, q.Options)
}

func (q *Question) Check(answer int) bool {
	switch q.Type {
	case "TF", "NUM", "MC":
		return answer == q.Correct[0]
	case "MS", "ORD":
		break
	}
	return false
}

/*
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
*/

func getInt(selected string, list []string) int {
	for i, option := range list {
		if option == selected {
			return i + 1
		}
	}
	return -1
}

func (q *Question) getNum() int {

	for {
		result, _ := pterm.DefaultInteractiveTextInput.Show(q.Text)
		i, err := strconv.Atoi(result)
		if err != nil {
			// ... handle error
			pterm.Error.Println("Invalid numeric input. Try again.")
		} else {
			return i
		}
	}

	/*
		stdin := bufio.NewReader(os.Stdin)
		var num int

		for {
			fmt.Printf("	=> ")
			_, err := fmt.Scanf("%d\n", &num)
			if err != nil {
				fmt.Printf("%v\n", err)
				stdin.ReadString('\n')
			} else if num > max || num < 1 {
				fmt.Printf("Invalid choice.\n")
			} else {
				return num
			}
		}
	*/
}
