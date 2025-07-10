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

var Command = &cobra.Command{
	Use:   "today",
	Short: "Suggests two DSA questions for today",
	Long:  `Suggests two DSA questions for today based on difficulty progression and smart review.`,
	Run:   todayCmd,
}

func todayCmd(cmd *cobra.Command, args []string) {
	todayQuestions, err := db.GetTodayQuestions()
	if err == nil && len(todayQuestions) > 0 {
		for idx, q := range todayQuestions {
			fmt.Printf("%d. %s (%s) - %s\n", idx+1, q.Name,
				strings.ToUpper(string(q.Difficulty[0]))+q.Difficulty[1:], q.URL)
		}
		return
	}

	easyQuestions, err := db.GetQuestionsByDifficulty("easy")
	if err != nil {
		color.Red("Error loading easy questions: %v", err)
		return
	}

	mediumQuestions, err := db.GetQuestionsByDifficulty("medium")
	if err != nil {
		color.Red("Error loading medium questions: %v", err)
		return
	}

	hardQuestions, err := db.GetQuestionsByDifficulty("hard")
	if err != nil {
		color.Red("Error loading hard questions: %v", err)
		return
	}

	allQuestions, err := db.GetAllQuestions()
	if err != nil {
		color.Red("Error loading all questions: %v", err)
		return
	}

	var questionsForToday []types.Question

	// Phase 1: Easy
	if !allAttempted(easyQuestions) {
		color.Green("Focusing on: Easy Questions")
		q1 := getFocusQuestion(easyQuestions)
		if q1 != nil {
			questionsForToday = append(questionsForToday, *q1)
			var remainingEasy []types.Question
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

			var attemptedPool []types.Question
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

			var attemptedPool []types.Question
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
		fmt.Println("No questions found")
		return
	}

	color.Cyan("Here are your questions for today:")

	for idx, q := range questionsForToday {
		fmt.Printf("%d. %s (%s) - %s\n", idx+1, q.Name,
			strings.ToUpper(string(q.Difficulty[0]))+q.Difficulty[1:], q.URL)
	}

	if err := db.InsertTodayQuestions(questionsForToday); err != nil {
		color.Red("Unable to fetch today's question.. Please try again..", err)
		return
	}
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

func getFocusQuestion(pool []types.Question) *types.Question {
	var unattempted []types.Question
	for _, q := range pool {
		if !q.Attempted {
			unattempted = append(unattempted, q)
		}
	}

	if len(unattempted) > 0 {
		return &unattempted[rand.Intn(len(unattempted))]
	}

	var attempted []types.Question
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
