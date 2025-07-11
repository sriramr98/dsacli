package complete

import (
	"dsacli/common"
	"dsacli/db"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const DateFormat = "2006-01-02T15:04:05"

var Command = &cobra.Command{
	Use:   "complete",
	Short: "Mark a question as complete",
	Long:  `Mark a question as complete and provide feedback to update its SR score.`,
	Run:   completeCmd,
}

func completeCmd(cmd *cobra.Command, args []string) {
	todaysQuestions, err := db.GetTodayQuestions()
	if err != nil {
		color.Red("Error loading today's questions: %v", err)
		return
	}

	if len(todaysQuestions) == 0 {
		color.Red("No questions found for today. Start by running 'dsacli today' to get today's questions.")
		return
	}

	questionPrompts := make([]string, len(todaysQuestions))
	for i, q := range todaysQuestions {
		questionPrompts[i] = fmt.Sprintf("%s (ID: %d)", q.Name, q.ID)
	}

	idx, err := common.PromptSelect("Select a question", questionPrompts)
	if err != nil {
		color.Red("Error reading input: %v", err)
		return
	}
	questionID := todaysQuestions[idx].ID

	questionToUpdate, err := db.FindQuestionByID(questionID)
	if err != nil {
		color.Red("Error finding question: %v", err)
		return
	}

	color.Cyan("You are about to update the question: %s (ID: %d)", questionToUpdate.Name, questionToUpdate.ID)

	hintsNeeded, err := common.PromptInt("Did you need hints? (1=many hints, 5=no hints)", common.OneToFiveRatingValidator)
	if err != nil {
		color.Red("Error reading input: %v", err)
		return
	}

	timeTaken, err := common.PromptInt("How long did it take (in minutes)? (-1 if you couldn't solve without solution)", common.NumberValidator)
	if err != nil {
		color.Red("Error reading input: %v", err)
		return
	}

	optimalSolution, err := common.PromptInt("Was the solution optimal? (1=not optimal, 5=very optimal)", common.OneToFiveRatingValidator)
	if err != nil {
		color.Red("Error reading input: %v", err)
		return
	}

	anyBugs, err := common.PromptInt("Were there any bugs? (1=many bugs, 5=no bugs)", common.OneToFiveRatingValidator)
	if err != nil {
		color.Red("Error reading input: %v", err)
		return
	}

	srScore := CalculateScore(timeTaken, hintsNeeded, optimalSolution, anyBugs, questionToUpdate)

	now := time.Now()
	questionToUpdate.LastReviewed = &now
	questionToUpdate.Attempted = true
	questionToUpdate.SRScore = srScore

	if err := db.UpdateQuestion(questionToUpdate); err != nil {
		color.Red("Error updating question: %v", err)
		return
	}

	if err := db.MarkTodayQuestionCompleted(questionID); err != nil {
		color.Yellow("Warning: Could not mark today's question as completed: %v", err)
	}

	color.Green("\nSuccessfully updated! New SR Score for '%s' is %d.",
		questionToUpdate.Name, srScore)
}
