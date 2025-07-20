package today

import (
	"dsacli/common"
	"dsacli/db"
	"dsacli/types"
	"fmt"
	"math/rand"
	"os/exec"
	"runtime"
	"sort"
	"strings"

	"github.com/fatih/color"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

const (
	questionsPerDay = 2
	easyPhase       = "easy"
	mediumPhase     = "medium"
	hardPhase       = "hard"
)

var More = false

func GetCommand(db db.Database) *cobra.Command {
	Command := &cobra.Command{
		Use:   "today",
		Short: "Suggests two DSA questions for today",
		Long:  `Suggests two DSA questions for today based on difficulty progression and smart review.`,
		Run:   todayCmd(db),
	}

	Command.Flags().BoolVarP(&More, "more", "m", false, "Show more questions (after completing today's questions)")

	return Command
}

func todayCmd(db db.Database) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		executeToday(db)
	}
}

func executeToday(db db.Database) {
	// Check if today's questions already exist
	questionsWithStatus, err := getTodayQuestionsWithStatusIfExist(db)
	if err == nil && len(questionsWithStatus) > 0 {
		// If all completed and more flag is set, generate new questions
		if !allCompleted(questionsWithStatus) {
			displayTodayQuestions(questionsWithStatus)
			return
		} else if !More {
			// If all questions are completed and more flag is not set, notify user.
			color.Green("You have already completed today's questions!")
			color.Yellow("Use --more flag to see more questions.")
			return
		} else {
			color.Cyan("You have completed today's questions. Generating new ones...")
		}
	}

	questionIdsToIgnore := make([]uint, 0)
	// If more flag is set, retrieve today's questions to avoid duplicates
	if More {
		for _, qws := range questionsWithStatus {
			questionIdsToIgnore = append(questionIdsToIgnore, qws.Question.ID)
		}
	}

	// Generate new questions for today
	questions, err := generateTodayQuestions(db, questionIdsToIgnore)
	if err != nil {
		color.Red("Error generating today's questions: %v", err)
		return
	}

	if len(questions) == 0 {
		fmt.Println("No questions found")
		return
	}

	if err := db.InsertTodayQuestions(questions); err != nil {
		color.Red("Unable to save today's questions: %v", err)
	}

	color.Cyan("Here are your questions for today:")
	displayQuestions(questions)

}

// getTodayQuestionsWithStatusIfExist retrieves today's questions with completion status if they already exist
func getTodayQuestionsWithStatusIfExist(db db.Database) ([]types.TodayQuestionWithStatus, error) {
	return db.GetTodayQuestionsWithStatus()
}

// generateTodayQuestions generates new questions based on difficulty progression
func generateTodayQuestions(db db.Database, questionsToIgnore []uint) ([]types.Question, error) {
	// Load questions by difficulty
	easyQuestions, err := db.GetQuestionsByDifficulty(easyPhase)
	if err != nil {
		return nil, fmt.Errorf("failed to load easy questions: %w", err)
	}
	if !allAttempted(easyQuestions) {
		questions := generateEasyPhaseQuestions(easyQuestions)
		if len(questions) > 0 {
			return questions, nil
		}
	}

	mediumQuestions, err := db.GetQuestionsByDifficulty(mediumPhase)
	if err != nil {
		return nil, fmt.Errorf("failed to load medium questions: %w", err)
	}
	if !allAttempted(mediumQuestions) {
		return generateMediumPhaseQuestions(mediumQuestions, easyQuestions), nil
	}

	hardQuestions, err := db.GetQuestionsByDifficulty(hardPhase)
	if err != nil {
		return nil, fmt.Errorf("failed to load hard questions: %w", err)
	}
	allQuestions, err := db.GetAllQuestions()
	if err != nil {
		return nil, fmt.Errorf("failed to load all questions: %w", err)
	}
	if !allAttempted(hardQuestions) {
		return generateHardPhaseQuestions(hardQuestions, allQuestions), nil
	}

	return generateMasteryPhaseQuestions(allQuestions), nil
}

// generateEasyPhaseQuestions generates questions for the easy phase
func generateEasyPhaseQuestions(easyQuestions []types.Question) []types.Question {
	color.Green("Focusing on: Easy Questions")

	var questions []types.Question

	// Get first question
	if q1, hasQuestion := getFocusQuestion(easyQuestions); hasQuestion {
		questions = append(questions, q1)

		// Get second question from remaining easy questions
		remaining := filterOutQuestion(easyQuestions, q1.ID)
		if q2, hasQuestion := getFocusQuestion(remaining); hasQuestion {
			questions = append(questions, q2)
		}
	}

	return questions
}

// generateMediumPhaseQuestions generates questions for the medium phase with smart review
func generateMediumPhaseQuestions(mediumQuestions, easyQuestions []types.Question) []types.Question {
	color.Yellow("Focusing on: Medium Questions (with Smart Review)")

	var questions []types.Question

	// Get focus question from medium difficulty
	if qFocus, hasQuestion := getFocusQuestion(mediumQuestions); hasQuestion {
		questions = append(questions, qFocus)

		// Get review question from attempted easy/medium questions
		attemptedPool := buildAttemptedPool(append(easyQuestions, mediumQuestions...), qFocus.ID)
		if qReview, hasQuestion := getHighestSRQuestion(attemptedPool); hasQuestion {
			questions = append(questions, qReview)
		}
	}

	return questions
}

// generateHardPhaseQuestions generates questions for the hard phase with smart review
func generateHardPhaseQuestions(hardQuestions, allQuestions []types.Question) []types.Question {
	color.Red("Focusing on: Hard Questions (with Smart Review)")

	var questions []types.Question

	// Get focus question from hard difficulty
	if qFocus, hasQuestion := getFocusQuestion(hardQuestions); hasQuestion {
		questions = append(questions, qFocus)

		// Get review question from all attempted questions
		attemptedPool := buildAttemptedPool(allQuestions, qFocus.ID)
		if qReview, hasQuestion := getHighestSRQuestion(attemptedPool); hasQuestion {
			questions = append(questions, qReview)
		}
	}

	return questions
}

// generateMasteryPhaseQuestions generates questions for the mastery phase
func generateMasteryPhaseQuestions(allQuestions []types.Question) []types.Question {
	color.Magenta("Mastery Mode: Reviewing all questions!")

	// Sort by SR score (highest first)
	sort.Slice(allQuestions, func(i, j int) bool {
		return allQuestions[i].LastPScore > allQuestions[j].LastPScore
	})

	if len(allQuestions) >= questionsPerDay {
		return allQuestions[:questionsPerDay]
	}

	return allQuestions
}

func displayQuestions(questions []types.Question) {
	var prompts []string = make([]string, len(questions))
	for idx, q := range questions {
		difficultyFormatted := strings.ToUpper(string(q.Difficulty[0])) + q.Difficulty[1:]
		prompts[idx] = fmt.Sprintf("%d. %s (%s)", idx+1, q.Name, difficultyFormatted)
	}
	idx, err := common.PromptSelect("Select a question to open in browser", prompts)
	if err != nil {
		color.Red("Error selecting question: %v", err)
		return
	}

	question := questions[idx]
	color.Cyan("Opening question: %s (%s)", question.Name, question.URL)
	// Open question.URL in the default browser
	if err := openBrowser(question.URL); err != nil {
		color.Red("Error opening browser: %v", err)
		return
	}
}

// displayTodayQuestions displays today's questions with their completion status
func displayTodayQuestions(questionsWithStatus []types.TodayQuestionWithStatus) {
	allCompleted := true
	totalCompleted := 0
	prompts := make([]string, 0)

	displayedQns := make([]types.Question, 0)
	for idx, qws := range questionsWithStatus {
		q := qws.Question
		difficultyFormatted := strings.ToUpper(string(q.Difficulty[0])) + q.Difficulty[1:]

		if !qws.Completed {
			prompts = append(prompts, fmt.Sprintf("%d. %s (%s) - %s", idx+1, q.Name, difficultyFormatted, q.URL))
			displayedQns = append(displayedQns, qws.Question)
			allCompleted = false
		} else {
			totalCompleted += 1
		}
	}

	// If all questions are completed, show congratulations message
	if allCompleted && len(questionsWithStatus) > 0 {
		fmt.Println()
		color.Green("🎉 Congratulations! You've completed all your questions for today!")
		color.Yellow("Great job on staying consistent with your DSA practice!")
	} else {
		color.Cyan("You have completed %d out of %d questions for today.", totalCompleted, len(questionsWithStatus))
		idx, err := common.PromptSelect("Select a question to open in browser", prompts)
		if err != nil {
			color.Red("Error selecting question: %v", err)
			return
		}

		question := displayedQns[idx]
		color.Cyan("Opening question: %s (%s)", question.Name, question.URL)
		// Open question.URL in the default browser
		if err := openBrowser(question.URL); err != nil {
			color.Red("Error opening browser: %v", err)
			return
		}
	}
}

func filterOutQuestion(questions []types.Question, excludeID uint) []types.Question {
	var filtered []types.Question
	for _, q := range questions {
		if q.ID != excludeID {
			filtered = append(filtered, q)
		}
	}
	return filtered
}

func buildAttemptedPool(questions []types.Question, excludeID uint) []types.Question {
	var attemptedPool []types.Question
	for _, q := range questions {
		if q.Attempted && q.ID != excludeID {
			attemptedPool = append(attemptedPool, q)
		}
	}
	return attemptedPool
}

func getHighestSRQuestion(questions []types.Question) (types.Question, bool) {
	if len(questions) == 0 {
		return types.Question{}, false
	}

	maxSR := questions[0]
	for _, q := range questions {
		if q.LastPScore > maxSR.LastPScore {
			maxSR = q
		}
	}

	return maxSR, true
}

func allAttempted(questions []types.Question) bool {
	if len(questions) == 0 {
		return true
	}

	for _, q := range questions {
		if !q.Attempted {
			return false
		}
	}

	return true
}

// getFocusQuestion returns the best question to focus on from the given pool
// It prioritizes unattempted questions first, then questions with highest SR score
func getFocusQuestion(pool []types.Question) (types.Question, bool) {
	if len(pool) == 0 {
		return types.Question{}, false
	}

	// First, try to get an unattempted question
	unattempted := filterUnattemptedQuestions(pool)
	if len(unattempted) > 0 {
		return unattempted[rand.Intn(len(unattempted))], true
	}

	// If all are attempted, get the one with highest SR score
	return getHighestSRQuestion(pool)
}

// filterUnattemptedQuestions returns only the questions that haven't been attempted
func filterUnattemptedQuestions(questions []types.Question) []types.Question {
	var unattempted []types.Question
	for _, q := range questions {
		if !q.Attempted {
			unattempted = append(unattempted, q)
		}
	}
	return unattempted
}

func allCompleted(questionsWithStatus []types.TodayQuestionWithStatus) bool {
	if len(questionsWithStatus) == 0 {
		return false
	}

	for _, qws := range questionsWithStatus {
		if !qws.Completed {
			return false
		}
	}

	return true
}

func openBrowser(url string) error {
	var err error
	// Cross-platform browser opening
	switch runtime.GOOS {
	case "darwin":
		// macOS
		_, err = exec.Command("open", url).Output()
	case "windows": //TODO: Test on Windows
		_, err = exec.Command("start", url).Output()
	default: //TODO: Test on Linux
		_, err = exec.Command("xdg-open", url).Output()
	}
	return err
}
