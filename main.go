package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

// Configuration constants
const (
	AppName    = "dsa-cli"
	Version    = "1.0.0"
	DBFilename = "dsa_questions.db"
	DateFormat = "2006-01-02T15:04:05"
)

// Question represents a DSA question
type Question struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	URL          string  `json:"url"`
	Difficulty   string  `json:"difficulty"`
	LastReviewed *string `json:"last_reviewed"`
	SRScore      int     `json:"sr_score"`
	Attempted    bool    `json:"attempted"`
}

func getAppDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(homeDir, "."+AppName)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}
	return appDir, nil
}

func getDBPath() (string, error) {
	appDir, err := getAppDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(appDir, DBFilename), nil
}

func initDB() (*sql.DB, error) {
	dbPath, err := getDBPath()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create questions table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS questions (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		url TEXT NOT NULL,
		difficulty TEXT NOT NULL,
		last_reviewed TEXT,
		sr_score INTEGER DEFAULT 0,
		attempted BOOLEAN DEFAULT 0
	);`

	if _, err := db.Exec(createTableSQL); err != nil {
		db.Close()
		return nil, err
	}

	// Check if table is empty and insert initial data
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM questions").Scan(&count)
	if err != nil {
		db.Close()
		return nil, err
	}

	if count == 0 {
		if err := insertInitialQuestions(db); err != nil {
			db.Close()
			return nil, err
		}
	}

	return db, nil
}

func insertInitialQuestions(db *sql.DB) error {
	initialQuestions := []Question{
		// Easy
		{ID: 0, Name: "Two Sum", URL: "https://leetcode.com/problems/two-sum/", Difficulty: "easy", LastReviewed: nil, SRScore: 0, Attempted: false},
		{ID: 1, Name: "Contains Duplicate", URL: "https://leetcode.com/problems/contains-duplicate/", Difficulty: "easy", LastReviewed: nil, SRScore: 0, Attempted: false},
		{ID: 2, Name: "Merge Two Sorted Lists", URL: "https://leetcode.com/problems/merge-two-sorted-lists/", Difficulty: "easy", LastReviewed: nil, SRScore: 0, Attempted: false},
		{ID: 3, Name: "Valid Anagram", URL: "https://leetcode.com/problems/valid-anagram/", Difficulty: "easy", LastReviewed: nil, SRScore: 0, Attempted: false},
		// Medium
		{ID: 4, Name: "Valid Parentheses", URL: "https://leetcode.com/problems/valid-parentheses/", Difficulty: "medium", LastReviewed: nil, SRScore: 0, Attempted: false},
		{ID: 5, Name: "Invert Binary Tree", URL: "https://leetcode.com/problems/invert-binary-tree/", Difficulty: "medium", LastReviewed: nil, SRScore: 0, Attempted: false},
		{ID: 6, Name: "Linked List Cycle", URL: "https://leetcode.com/problems/linked-list-cycle/", Difficulty: "medium", LastReviewed: nil, SRScore: 0, Attempted: false},
		{ID: 7, Name: "Best Time to Buy and Sell Stock", URL: "https://leetcode.com/problems/best-time-to-buy-and-sell-stock/", Difficulty: "medium", LastReviewed: nil, SRScore: 0, Attempted: false},
		// Hard
		{ID: 8, Name: "Find Median from Data Stream", URL: "https://leetcode.com/problems/find-median-from-data-stream/", Difficulty: "hard", LastReviewed: nil, SRScore: 0, Attempted: false},
		{ID: 9, Name: "Word Ladder", URL: "https://leetcode.com/problems/word-ladder/", Difficulty: "hard", LastReviewed: nil, SRScore: 0, Attempted: false},
	}

	insertSQL := `INSERT INTO questions (id, name, url, difficulty, last_reviewed, sr_score, attempted) 
	              VALUES (?, ?, ?, ?, ?, ?, ?)`

	for _, q := range initialQuestions {
		_, err := db.Exec(insertSQL, q.ID, q.Name, q.URL, q.Difficulty, q.LastReviewed, q.SRScore, q.Attempted)
		if err != nil {
			return err
		}
	}

	return nil
}

func getAllQuestions(db *sql.DB) ([]Question, error) {
	rows, err := db.Query("SELECT id, name, url, difficulty, last_reviewed, sr_score, attempted FROM questions ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []Question
	for rows.Next() {
		var q Question
		err := rows.Scan(&q.ID, &q.Name, &q.URL, &q.Difficulty, &q.LastReviewed, &q.SRScore, &q.Attempted)
		if err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}

	return questions, rows.Err()
}

func getQuestionsByDifficulty(db *sql.DB, difficulty string) ([]Question, error) {
	rows, err := db.Query("SELECT id, name, url, difficulty, last_reviewed, sr_score, attempted FROM questions WHERE difficulty = ? ORDER BY id", difficulty)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []Question
	for rows.Next() {
		var q Question
		err := rows.Scan(&q.ID, &q.Name, &q.URL, &q.Difficulty, &q.LastReviewed, &q.SRScore, &q.Attempted)
		if err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}

	return questions, rows.Err()
}

func updateQuestion(db *sql.DB, q Question) error {
	updateSQL := `UPDATE questions SET name = ?, url = ?, difficulty = ?, last_reviewed = ?, sr_score = ?, attempted = ? WHERE id = ?`
	_, err := db.Exec(updateSQL, q.Name, q.URL, q.Difficulty, q.LastReviewed, q.SRScore, q.Attempted, q.ID)
	return err
}

func findQuestionByID(db *sql.DB, id int) (*Question, error) {
	var q Question
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

func allAttempted(questions []Question) bool {
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

func getFocusQuestion(pool []Question) *Question {
	var unattempted []Question
	for _, q := range pool {
		if !q.Attempted {
			unattempted = append(unattempted, q)
		}
	}

	if len(unattempted) > 0 {
		return &unattempted[rand.Intn(len(unattempted))]
	}

	var attempted []Question
	for _, q := range pool {
		if q.Attempted {
			attempted = append(attempted, q)
		}
	}

	if len(attempted) > 0 {
		maxSR := attempted[0]
		for _, q := range attempted {
			if q.SRScore > maxSR.SRScore {
				maxSR = q
			}
		}
		return &maxSR
	}

	return nil
}

func promptInt(question string) (int, error) {
	fmt.Print(question + " ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	input = strings.TrimSpace(input)
	return strconv.Atoi(input)
}

// Commands
func todayCmd(cmd *cobra.Command, args []string) {
	db, err := initDB()
	if err != nil {
		color.Red("Error initializing database: %v", err)
		return
	}
	defer db.Close()

	easyQuestions, err := getQuestionsByDifficulty(db, "easy")
	if err != nil {
		color.Red("Error loading easy questions: %v", err)
		return
	}

	mediumQuestions, err := getQuestionsByDifficulty(db, "medium")
	if err != nil {
		color.Red("Error loading medium questions: %v", err)
		return
	}

	hardQuestions, err := getQuestionsByDifficulty(db, "hard")
	if err != nil {
		color.Red("Error loading hard questions: %v", err)
		return
	}

	allQuestions, err := getAllQuestions(db)
	if err != nil {
		color.Red("Error loading all questions: %v", err)
		return
	}

	var questionsForToday []Question

	// Phase 1: Easy
	if !allAttempted(easyQuestions) {
		color.Green("Focusing on: Easy Questions")
		q1 := getFocusQuestion(easyQuestions)
		if q1 != nil {
			questionsForToday = append(questionsForToday, *q1)
			var remainingEasy []Question
			for _, q := range easyQuestions {
				if q.ID != q1.ID {
					remainingEasy = append(remainingEasy, q)
				}
			}
			if len(remainingEasy) > 0 {
				q2 := getFocusQuestion(remainingEasy)
				if q2 != nil {
					questionsForToday = append(questionsForToday, *q2)
				}
			}
		}
	} else if !allAttempted(mediumQuestions) {
		// Phase 2: Medium + Smart Review
		color.Yellow("Focusing on: Medium Questions (with Smart Review)")
		qFocus := getFocusQuestion(mediumQuestions)
		if qFocus != nil {
			questionsForToday = append(questionsForToday, *qFocus)

			var attemptedPool []Question
			for _, q := range append(easyQuestions, mediumQuestions...) {
				if q.Attempted && q.ID != qFocus.ID {
					attemptedPool = append(attemptedPool, q)
				}
			}

			if len(attemptedPool) > 0 {
				maxSR := attemptedPool[0]
				for _, q := range attemptedPool {
					if q.SRScore > maxSR.SRScore {
						maxSR = q
					}
				}
				questionsForToday = append(questionsForToday, maxSR)
			}
		}
	} else if !allAttempted(hardQuestions) {
		// Phase 3: Hard + Smart Review
		color.Red("Focusing on: Hard Questions (with Smart Review)")
		qFocus := getFocusQuestion(hardQuestions)
		if qFocus != nil {
			questionsForToday = append(questionsForToday, *qFocus)

			var attemptedPool []Question
			for _, q := range allQuestions {
				if q.Attempted && q.ID != qFocus.ID {
					attemptedPool = append(attemptedPool, q)
				}
			}

			if len(attemptedPool) > 0 {
				maxSR := attemptedPool[0]
				for _, q := range attemptedPool {
					if q.SRScore > maxSR.SRScore {
						maxSR = q
					}
				}
				questionsForToday = append(questionsForToday, maxSR)
			}
		}
	} else {
		// Phase 4: Full Review
		color.Magenta("Mastery Mode: Reviewing all questions!")
		sort.Slice(allQuestions, func(i, j int) bool {
			return allQuestions[i].SRScore > allQuestions[j].SRScore
		})
		if len(allQuestions) >= 2 {
			questionsForToday = allQuestions[:2]
		} else {
			questionsForToday = allQuestions
		}
	}

	if len(questionsForToday) == 0 {
		dbPath, _ := getDBPath()
		fmt.Printf("No questions found. Your database is located at: %s\n", dbPath)
		return
	}

	color.Cyan("Here are your questions for today:")

	// Remove duplicates
	uniqueQuestions := make(map[int]Question)
	for _, q := range questionsForToday {
		uniqueQuestions[q.ID] = q
	}

	i := 1
	for _, q := range uniqueQuestions {
		fmt.Printf("%d. %s (%s) - %s\n", i, q.Name,
			strings.ToUpper(string(q.Difficulty[0]))+q.Difficulty[1:], q.URL)
		i++
	}
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

	db, err := initDB()
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

	hintsNeeded, err := promptInt("Did you need hints? (1=many hints, 5=no hints)")
	if err != nil {
		color.Red("Error reading input: %v", err)
		return
	}

	timeTaken, err := promptInt("How long did it take (in minutes)? (-1 if you couldn't solve without solution)")
	if err != nil {
		color.Red("Error reading input: %v", err)
		return
	}

	optimalSolution, err := promptInt("Was the solution optimal? (1=not optimal, 5=very optimal)")
	if err != nil {
		color.Red("Error reading input: %v", err)
		return
	}

	anyBugs, err := promptInt("Were there any bugs? (1=many bugs, 5=no bugs)")
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

func versionCmd(cmd *cobra.Command, args []string) {
	fmt.Printf("dsacli version %s\n", Version)
}

func listCmd(cmd *cobra.Command, args []string) {
	db, err := initDB()
	if err != nil {
		color.Red("Error initializing database: %v", err)
		return
	}
	defer db.Close()

	questions, err := getAllQuestions(db)
	if err != nil {
		color.Red("Error loading questions: %v", err)
		return
	}

	color.Cyan("All DSA Questions:")
	fmt.Println()

	var easyQuestions, mediumQuestions, hardQuestions []Question
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

	printQuestionList := func(title string, questions []Question, titleColor *color.Color) {
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

func main() {
	var rootCmd = &cobra.Command{
		Use:   "dsacli",
		Short: "A CLI tool to practice DSA questions using spaced repetition",
		Long:  `A CLI tool to practice DSA questions using a spaced repetition algorithm with difficulty progression.`,
	}

	var todayCommand = &cobra.Command{
		Use:   "today",
		Short: "Suggests two DSA questions for today",
		Long:  `Suggests two DSA questions for today based on difficulty progression and smart review.`,
		Run:   todayCmd,
	}

	var completeCommand = &cobra.Command{
		Use:   "complete [question_id]",
		Short: "Mark a question as complete",
		Long:  `Mark a question as complete and provide feedback to update its SR score. Use the question ID from 'dsacli list'.`,
		Args:  cobra.ExactArgs(1),
		Run:   completeCmd,
	}

	var listCommand = &cobra.Command{
		Use:   "list",
		Short: "List all questions with their IDs",
		Long:  `List all questions with their IDs, completion status, and SR scores.`,
		Run:   listCmd,
	}

	var versionCommand = &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  `Show the current version of the DSA CLI tool.`,
		Run:   versionCmd,
	}

	rootCmd.AddCommand(todayCommand)
	rootCmd.AddCommand(completeCommand)
	rootCmd.AddCommand(listCommand)
	rootCmd.AddCommand(versionCommand)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
