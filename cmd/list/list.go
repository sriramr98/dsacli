package list

import (
	"dsacli/db"
	"dsacli/types"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var shortForm = true
var longForm = false

var Command = &cobra.Command{
	Use:   "list",
	Short: "List all questions with their IDs",
	Long:  `List all questions with their IDs, completion status, and SR scores.`,
	Run:   listCmd,
}

func init() {
	Command.Flags().BoolVarP(&shortForm, "short", "s", true, "Prints category wise stats only")
	Command.Flags().BoolVarP(&longForm, "long", "l", false, "Prints all questions with IDs, completion status, and SR scores")
}

func listCmd(cmd *cobra.Command, args []string) {
	if longForm {
		shortForm = false
	}

	if !longForm && !shortForm {
		shortForm = true
	}

	questions, err := db.GetAllQuestions()
	if err != nil {
		color.Red("Error loading questions: %v", err)
		return
	}

	color.Cyan("All DSA Questions:")
	fmt.Println()

	var easyQuestions, mediumQuestions, hardQuestions []types.Question
	for _, q := range questions {
		switch q.Difficulty {
		case "easy":
			easyQuestions = append(easyQuestions, q)
		case "medium":
			mediumQuestions = append(mediumQuestions, q)
		case "hard":
			hardQuestions = append(hardQuestions, q)
		}
	}

	printQuestionList := func(title string, questions []types.Question, titleColor *color.Color) {
		titleColor.Printf("\n%s: (%d)\n", title, len(questions))

		totalAttempted := 0
		longFormLogs := make([]string, 0)
		for _, q := range questions {
			status := "❌"
			if q.Attempted {
				totalAttempted++
				status = "✅"
			}

			longFormLogs = append(longFormLogs, fmt.Sprintf("  %s ID:%d - %s (SR Score: %d)\n", status, q.ID, q.Name, q.SRScore))
		}

		color.White("	- Total Attempted: %d\n", totalAttempted)
		color.White("	- Total UnAttempted: %d\n", len(questions)-totalAttempted)
		if shortForm {
			return
		}

		for _, log := range longFormLogs {
			fmt.Print(log)
		}
	}

	printQuestionList("Easy Questions", easyQuestions, color.New(color.FgGreen))
	printQuestionList("Medium Questions", mediumQuestions, color.New(color.FgYellow))
	printQuestionList("Hard Questions", hardQuestions, color.New(color.FgRed))

	fmt.Printf("\nTotal Questions: %d\n", len(questions))
	fmt.Println("\nTo mark a question as complete, use: dsacli complete <question_id>")
}
