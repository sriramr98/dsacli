package list

import (
	"dsacli/db"
	"dsacli/types"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var ListCommand = &cobra.Command{
	Use:   "list",
	Short: "List all questions with their IDs",
	Long:  `List all questions with their IDs, completion status, and SR scores.`,
	Run:   listCmd,
}

func listCmd(cmd *cobra.Command, args []string) {
	sqlDB, err := db.GetDB()
	if err != nil {
		color.Red("Error initializing database: %v", err)
		return
	}
	defer sqlDB.Close()

	questions, err := db.GetAllQuestions(sqlDB)
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
		if len(questions) > 0 {
			titleColor.Printf("\n%s:\n", title)
			for _, q := range questions {
				status := "❌"
				if q.Attempted {
					status = "✅"
				}
				fmt.Printf("  %s ID:%d - %s (SR Score: %d)\n", status, q.ID, q.Name, q.SRScore)
			}
		}
	}

	printQuestionList("Easy Questions", easyQuestions, color.New(color.FgGreen))
	printQuestionList("Medium Questions", mediumQuestions, color.New(color.FgYellow))
	printQuestionList("Hard Questions", hardQuestions, color.New(color.FgRed))

	fmt.Printf("\nTotal Questions: %d\n", len(questions))
	fmt.Println("\nTo mark a question as complete, use: dsacli complete <question_id>")
}
