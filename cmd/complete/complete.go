package complete

import (
	"dsacli/common"
	"dsacli/db"
	"fmt"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const DateFormat = "2006-01-02T15:04:05"

var Command = &cobra.Command{
	Use:   "complete [question_id]",
	Short: "Mark a question as complete",
	Long:  `Mark a question as complete and provide feedback to update its SR score. Use the question ID from 'dsacli list'.`,
	Args:  cobra.ExactArgs(1),
	Run:   completeCmd,
}

func completeCmd(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		color.Red("Error: Please provide the question ID.")
		color.White("Usage: dsacli complete <question_id>")
		color.White("Use 'dsacli list' to see all question IDs.")
		return
	}

	questionID, err := strconv.Atoi(args[0])
	if err != nil {
		color.Red("Error: Invalid question ID. Please provide a number.")
		return
	}

	questionToUpdate, err := db.FindQuestionByID(questionID)
	if err != nil {
		color.Red("Error finding question: %v", err)
		return
	}

	fmt.Printf("Updating score for: %s\n", color.New(color.Bold).Sprintf(questionToUpdate.Name))

	hintsNeeded, err := common.PromptInt("Did you need hints? (1=many hints, 5=no hints)")
	if err != nil {
		color.Red("Error reading input: %v", err)
		return
	}

	timeTaken, err := common.PromptInt("How long did it take (in minutes)? (-1 if you couldn't solve without solution)")
	if err != nil {
		color.Red("Error reading input: %v", err)
		return
	}

	optimalSolution, err := common.PromptInt("Was the solution optimal? (1=not optimal, 5=very optimal)")
	if err != nil {
		color.Red("Error reading input: %v", err)
		return
	}

	anyBugs, err := common.PromptInt("Were there any bugs? (1=many bugs, 5=no bugs)")
	if err != nil {
		color.Red("Error reading input: %v", err)
		return
	}

	srScore := CalculateScore(timeTaken, hintsNeeded, optimalSolution, anyBugs, questionToUpdate)

	nowStr := time.Now().Format(DateFormat)
	questionToUpdate.LastReviewed = &nowStr
	questionToUpdate.Attempted = true
	questionToUpdate.SRScore = srScore

	if err := db.UpdateQuestion(questionToUpdate); err != nil {
		color.Red("Error updating question: %v", err)
		return
	}

	color.Green("\nSuccessfully updated! New SR Score for '%s' is %d.",
		questionToUpdate.Name, srScore)
}
