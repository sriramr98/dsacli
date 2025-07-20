package complete

import (
	"dsacli/types"
	"math"
)

// CalculatePScore computes the performance score based on user feedback
// Returns a float between 0.0 and 1.0
func CalculatePScore(timeTaken, hintsUsed, optimality, bugs int) float64 {
	// Maximum 40% weightage for time taken - Score gets converted to 0.0-0.4 scale
	var timeScore float64
	switch {
	case timeTaken == -1: // Unsolved
		timeScore = 0.0
	case timeTaken <= 30:
		timeScore = 0.4
	case timeTaken <= 45:
		timeScore = 0.2
	default: // timeTaken > 45
		timeScore = 0.1
	}

	// Hint score calculation (convert 1-5 scale to 0.0-0.3)
	var hintScore float64
	if hintsUsed == 0 {
		hintScore = 0.3
	} else {
		hintScore = 0.3 / float64(hintsUsed+1)
	}

	// Optimality score calculation (convert 1-5 scale to 0.0-0.15)
	optimalityScore := float64(optimality-1) / 4.0 * 0.15

	// Bug score calculation (convert 1-5 scale to 0.0-0.15)
	bugScore := float64(bugs-1) / 4.0 * 0.15

	return timeScore + hintScore + optimalityScore + bugScore
}

// ProcessReview updates the question based on spaced repetition logic after user completes a problem
func ProcessReview(question *types.Question, timeTaken, hintsUsed, optimality, bugs int) {
	// Increment attempt count
	question.AttemptCount++

	// Calculate current p-score
	currentPScore := CalculatePScore(timeTaken, hintsUsed, optimality, bugs)

	// Check for Progression Mastery (if not already achieved)
	if !question.Mastered {
		// Instant Mastery: first attempt with p_score >= 0.95
		if question.AttemptCount == 1 && currentPScore >= 0.95 {
			question.Mastered = true
		}

		// Proven Mastery: current and previous attempts both >= 0.85
		if question.AttemptCount > 1 && currentPScore >= 0.85 && question.LastPScore >= 0.85 {
			question.Mastered = true
		}
	}

	// Update Review Schedule
	if currentPScore >= 0.6 { // Successful Recall
		// This means that the user successfully recalled the question. A streak of two successful recalls will lead to mastery.
		question.ReviewStreak++

		// Reasoning for this is present in SPACED_REPETITION.md
		newEF := question.EasinessFactor + (0.1 - (0.85-currentPScore)*(0.08+(0.85-currentPScore)*0.02))
		if newEF < 1.3 {
			newEF = 1.3
		}
		question.EasinessFactor = newEF

		// Calculate review interval
		switch question.ReviewStreak {
		case 1:
			question.ReviewInterval = 1
		case 2:
			question.ReviewInterval = 6
		default: // ReviewStreak > 2
			previousInterval := float64(question.ReviewInterval)
			question.ReviewInterval = int(math.Round(previousInterval * newEF))
		}
	} else { // Failed Recall (p_score < 0.6)
		// Reset review streak
		question.ReviewStreak = 0

		// Set review interval to 1 day
		question.ReviewInterval = 1

		// Penalize easiness factor
		newEF := question.EasinessFactor - 0.2
		if newEF < 1.3 {
			newEF = 1.3
		}
		question.EasinessFactor = newEF
	}

	// Store current p-score as last p-score for next attempt
	question.LastPScore = currentPScore
}

// CheckProgressionGate determines if a difficulty tier is unlocked
// Returns true if > 50% of questions in the category have progression mastery
func CheckProgressionGate(categoryQuestions []types.Question) bool {
	if len(categoryQuestions) == 0 {
		return false
	}

	masteredCount := 0
	for _, q := range categoryQuestions {
		if q.Mastered {
			masteredCount++
		}
	}

	masteryPercentage := float64(masteredCount) / float64(len(categoryQuestions)) * 100
	return masteryPercentage > 50.0
}
