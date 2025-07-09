package complete

import (
	"database/sql"
	"dsacli/common"
	"dsacli/db"
	"dsacli/types"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const DateFormat = "2006-01-02T15:04:05"

var CompleteCommand = &cobra.Command{
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

	db, err := db.GetDB()
	if err != nil {
		color.Red("Error initializing database: %v", err)
		return
	}
	defer db.Close()

	questionToUpdate, err := findQuestionByID(db, questionID)
	if err != nil {
		color.Red("Error finding question: %v", err)
		return
	}

	if questionToUpdate == nil {
		color.Red("Error: Could not find question with ID %d.", questionID)
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

	var timeRank float64
	if timeTaken == -1 {
		timeRank = 0
	} else if timeTaken > 45 {
		timeRank = 2
	} else if timeTaken >= 30 && timeTaken <= 45 {
		timeRank = 3.5
	} else {
		timeRank = 5
	}

	score := (float64(hintsNeeded) + timeRank + float64(optimalSolution) + float64(anyBugs)) / 4

	now := time.Now()
	cDate := 10000.0
	if questionToUpdate.LastReviewed != nil {
		lastReviewedDate, err := time.Parse(DateFormat, *questionToUpdate.LastReviewed)
		if err == nil {
			delta := now.Sub(lastReviewedDate)
			cDate = delta.Minutes()
		}
	}

	var cSolution float64
	if score == 5 {
		cSolution = 0.5
	} else {
		cSolution = (5 - score) + 1
	}

	var cTime float64
	if timeTaken == -1 {
		cTime = 60 * 400
	} else if timeTaken < 25 {
		cTime = float64(timeTaken) * 100
	} else if timeTaken < 35 {
		cTime = float64(timeTaken) * 200
	} else if timeTaken < 45 {
		cTime = float64(timeTaken) * 300
	} else {
		cTime = float64(timeTaken) * 400
	}

	srScore := int(math.Round((cDate + cTime) * cSolution))

	if questionToUpdate.SRScore == 0 {
		questionToUpdate.SRScore = srScore
	} else {
		questionToUpdate.SRScore = int(float64(questionToUpdate.SRScore)*0.7 + float64(srScore)*0.3)
	}
	nowStr := now.Format(DateFormat)
	questionToUpdate.LastReviewed = &nowStr
	questionToUpdate.Attempted = true

	if err := updateQuestion(db, *questionToUpdate); err != nil {
		color.Red("Error updating question: %v", err)
		return
	}

	color.Green("\nSuccessfully updated! New SR Score for '%s' is %d.",
		questionToUpdate.Name, srScore)
}

func findQuestionByID(db *sql.DB, id int) (*types.Question, error) {
	var q types.Question
	err := db.QueryRow("SELECT id, name, url, difficulty, last_reviewed, sr_score, attempted FROM questions WHERE id = ?", id).
		Scan(&q.ID, &q.Name, &q.URL, &q.Difficulty, &q.LastReviewed, &q.SRScore, &q.Attempted)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &q, nil
}

func updateQuestion(db *sql.DB, q types.Question) error {
	updateSQL := `UPDATE questions SET name = ?, url = ?, difficulty = ?, last_reviewed = ?, sr_score = ?, attempted = ? WHERE id = ?`
	_, err := db.Exec(updateSQL, q.Name, q.URL, q.Difficulty, q.LastReviewed, q.SRScore, q.Attempted, q.ID)
	return err
}
