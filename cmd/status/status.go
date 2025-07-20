package status

import (
	"dsacli/common"
	"dsacli/db"
	"dsacli/types"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func GetCommand(db db.Database) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show the status of learning progression",
		Long:  "Shows the state of solved and unsolved problems",
		Run:   printStatus(db),
	}
}

func printStatus(db db.Database) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		attempted, err := db.GetAllAttemptedQuestions()
		if err != nil {
			cmd.Println("Error fetching attempted questions:", err)
			return
		}

		masteredQns := common.FilterSlice(attempted, func(q types.Question) bool {
			return q.Mastered
		})

		if len(masteredQns) == 0 {
			color.Yellow("No questions marked as mastered yet.")
		} else {
			color.Green("Mastered Questions (%d):", len(masteredQns))
			for _, q := range masteredQns {
				cmd.Printf("    - %s\n", q.Name)
			}
		}

		nonMasteredQns := common.FilterSlice(attempted, func(q types.Question) bool {
			return !q.Mastered
		})

		if len(nonMasteredQns) == 0 {
			color.Yellow("All attempted questions are mastered.")
		} else {
			color.Red("Non-Mastered Questions (%d):", len(nonMasteredQns))
			for _, q := range nonMasteredQns {
				cmd.Printf("    - %s\n", q.Name)
			}
		}

	}
}
