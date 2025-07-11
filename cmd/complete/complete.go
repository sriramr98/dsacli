package complete

import (
	"dsacli/common"
	"dsacli/db"
	"dsacli/types"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	// DateFormat is the standard date format used throughout the application
	DateFormat = "2006-01-02T15:04:05"

	// Rating bounds for user feedback
	MinRating = 1
	MaxRating = 5

	// Special value for unsolved questions
	UnsolvedTimeValue = -1
)

// Command is the cobra command for marking questions as complete
var Command = &cobra.Command{
	Use:   "complete",
	Short: "Mark a question as complete",
	Long:  `Mark a question as complete and provide feedback to update its SR score.`,
	Run:   completeCmd,
}

// CompletionFeedback holds the user's feedback about completing a question
type CompletionFeedback struct {
	HintsNeeded     int
	TimeTaken       int
	OptimalSolution int
	AnyBugs         int
}

// completeCmd is the cobra command handler
func completeCmd(cmd *cobra.Command, args []string) {
	if err := executeComplete(); err != nil {
		color.Red("Error: %v", err)
		return
	}
}

// executeComplete contains the main business logic for completing questions
func executeComplete() error {
	// Get today's questions
	todaysQuestions, err := db.GetTodayQuestions()
	if err != nil {
		return fmt.Errorf("loading today's questions: %w", err)
	}

	if len(todaysQuestions) == 0 {
		color.Red("No questions found for today. Start by running 'dsacli today' to get today's questions.")
		return nil // This is not an error, just a message to the user
	}

	// Let user select a question
	questionID, err := selectQuestion(todaysQuestions)
	if err != nil {
		return fmt.Errorf("selecting question: %w", err)
	}

	// Get the question details
	questionToUpdate, err := db.FindQuestionByID(questionID)
	if err != nil {
		return fmt.Errorf("finding question: %w", err)
	}

	// Display question info
	color.Cyan("You are about to update the question: %s (ID: %d)", questionToUpdate.Name, questionToUpdate.ID)

	// Collect user feedback
	feedback, err := collectFeedback()
	if err != nil {
		return fmt.Errorf("collecting feedback: %w", err)
	}

	// Update question with feedback
	if err := updateQuestionWithFeedback(&questionToUpdate, feedback); err != nil {
		return fmt.Errorf("updating question: %w", err)
	}

	// Save to database
	if err := db.UpdateQuestion(questionToUpdate); err != nil {
		return fmt.Errorf("saving question: %w", err)
	}

	// Mark as completed for today
	if err := db.MarkTodayQuestionCompleted(questionID); err != nil {
		color.Yellow("Warning: Could not mark today's question as completed: %v", err)
	}

	color.Green("\nSuccessfully updated! New SR Score for '%s' is %d.",
		questionToUpdate.Name, questionToUpdate.SRScore)

	return nil
}

// selectQuestion prompts the user to select a question from today's questions
func selectQuestion(questions []types.Question) (uint, error) {
	questionPrompts := make([]string, len(questions))
	for i, q := range questions {
		questionPrompts[i] = fmt.Sprintf("%s (ID: %d)", q.Name, q.ID)
	}

	idx, err := common.PromptSelect("Select a question", questionPrompts)
	if err != nil {
		return 0, fmt.Errorf("reading input: %w", err)
	}

	return questions[idx].ID, nil
}

// collectFeedback prompts the user for feedback about the completed question
func collectFeedback() (CompletionFeedback, error) {
	feedback := CompletionFeedback{}

	var err error
	feedback.HintsNeeded, err = common.PromptInt("Did you need hints? (1=many hints, 5=no hints)", common.OneToFiveRatingValidator)
	if err != nil {
		return feedback, fmt.Errorf("reading hints input: %w", err)
	}

	feedback.TimeTaken, err = common.PromptInt("How long did it take (in minutes)? (-1 if you couldn't solve without solution)", common.NumberValidator)
	if err != nil {
		return feedback, fmt.Errorf("reading time input: %w", err)
	}

	feedback.OptimalSolution, err = common.PromptInt("Was the solution optimal? (1=not optimal, 5=very optimal)", common.OneToFiveRatingValidator)
	if err != nil {
		return feedback, fmt.Errorf("reading optimality input: %w", err)
	}

	feedback.AnyBugs, err = common.PromptInt("Were there any bugs? (1=many bugs, 5=no bugs)", common.OneToFiveRatingValidator)
	if err != nil {
		return feedback, fmt.Errorf("reading bugs input: %w", err)
	}

	return feedback, nil
}

// updateQuestionWithFeedback updates the question with the user's feedback
func updateQuestionWithFeedback(question *types.Question, feedback CompletionFeedback) error {
	// Calculate new SR score
	srScore := CalculateScore(feedback.TimeTaken, feedback.HintsNeeded, feedback.OptimalSolution, feedback.AnyBugs, *question)

	// Update question fields
	now := time.Now()
	question.LastReviewed = &now
	question.Attempted = true
	question.SRScore = srScore

	return nil
}
