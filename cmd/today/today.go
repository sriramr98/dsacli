package today

import (
	"dsacli/db"
	"dsacli/types"
	"fmt"
	"math/rand"
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

// Command represents the today command
var Command = &cobra.Command{
	Use:   "today",
	Short: "Suggests two DSA questions for today",
	Long:  `Suggests two DSA questions for today based on difficulty progression and smart review.`,
	Run:   todayCmd,
}

// todayCmd handles the today command execution
func todayCmd(cmd *cobra.Command, args []string) {
	// Check if today's questions already exist
	if questionsWithStatus, err := getTodayQuestionsWithStatusIfExist(); err == nil && len(questionsWithStatus) > 0 {
		displayTodayQuestions(questionsWithStatus)
		return
	}

	// Generate new questions for today
	questions, err := generateTodayQuestions()
	if err != nil {
		color.Red("Error generating today's questions: %v", err)
		return
	}

	if len(questions) == 0 {
		fmt.Println("No questions found")
		return
	}

	// Display and save today's questions
	color.Cyan("Here are your questions for today:")
	displayQuestions(questions)

	if err := db.InsertTodayQuestions(questions); err != nil {
		color.Red("Unable to save today's questions: %v", err)
	}
}

// getTodayQuestionsWithStatusIfExist retrieves today's questions with completion status if they already exist
func getTodayQuestionsWithStatusIfExist() ([]types.TodayQuestionWithStatus, error) {
	return db.GetTodayQuestionsWithStatus()
}

// generateTodayQuestions generates new questions based on difficulty progression
func generateTodayQuestions() ([]types.Question, error) {
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
		return allQuestions[i].SRScore > allQuestions[j].SRScore
	})

	if len(allQuestions) >= questionsPerDay {
		return allQuestions[:questionsPerDay]
	}

	return allQuestions
}

func displayQuestions(questions []types.Question) {
	for idx, q := range questions {
		difficultyFormatted := strings.ToUpper(string(q.Difficulty[0])) + q.Difficulty[1:]
		fmt.Printf("%d. %s (%s) - %s ❌\n", idx+1, q.Name, difficultyFormatted, q.URL)
	}
}

// displayTodayQuestions displays today's questions with their completion status
func displayTodayQuestions(questionsWithStatus []types.TodayQuestionWithStatus) {
	allCompleted := true

	for idx, qws := range questionsWithStatus {
		q := qws.Question
		difficultyFormatted := strings.ToUpper(string(q.Difficulty[0])) + q.Difficulty[1:]

		statusIcon := "❌" // Not completed
		if qws.Completed {
			statusIcon = "✅" // Completed
		} else {
			allCompleted = false
		}

		fmt.Printf("%d. %s (%s) - %s %s\n", idx+1, q.Name, difficultyFormatted, q.URL, statusIcon)
	}

	// If all questions are completed, show congratulations message
	if allCompleted && len(questionsWithStatus) > 0 {
		fmt.Println()
		color.Green("🎉 Congratulations! You've completed all your questions for today!")
		color.Yellow("Great job on staying consistent with your DSA practice!")
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
		if q.SRScore > maxSR.SRScore {
			maxSR = q
		}
	}

	return maxSR, false
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
