package complete

import (
	"dsacli/types"
	"math"
	"time"
)

const (
	// Time thresholds for scoring (in minutes)
	FastSolutionThreshold   = 25
	MediumSolutionThreshold = 35 // Changed from 30 to 35
	SlowSolutionThreshold   = 45

	// Score multipliers
	FastTimeMultiplier     = 100
	MediumTimeMultiplier   = 200
	SlowTimeMultiplier     = 300
	VerySlowTimeMultiplier = 400
	UnsolvedTimeMultiplier = 60 * 400

	// Score weights
	PreviousScoreWeight = 0.7
	CurrentScoreWeight  = 0.3

	// Default values
	DefaultReviewInterval = 10000.0 // minutes
	PerfectScore          = 5
	OptimalMultiplier     = 0.5
)

// CalculateScore computes the spaced repetition score based on user feedback
func CalculateScore(timeTaken, hintsNeeded, optimalSolution, anyBugs int, question types.Question) int {
	timeRank := calculateTimeRank(timeTaken)

	// Score is the average rating of these 4 parameters
	averageScore := (float64(hintsNeeded) + timeRank + float64(optimalSolution) + float64(anyBugs)) / 4

	reviewInterval := calculateReviewInterval(question.LastReviewed)
	solutionMultiplier := calculateSolutionMultiplier(averageScore)
	timeFactor := calculateTimeFactor(timeTaken)

	newScore := int(math.Round((reviewInterval + timeFactor) * solutionMultiplier))

	// If this is the first attempt, return the new score
	if question.LastPScore == 0 {
		return newScore
	}

	// Otherwise, blend with previous score
	return int(float64(question.LastPScore)*PreviousScoreWeight + float64(newScore)*CurrentScoreWeight)
}

// calculateTimeRank converts time taken into a 1-5 rating scale
func calculateTimeRank(timeTaken int) float64 {
	switch {
	case timeTaken == UnsolvedTimeValue:
		return 0 // Couldn't solve
	case timeTaken > SlowSolutionThreshold:
		return 2 // Very slow
	case timeTaken >= MediumSolutionThreshold:
		return 3.5 // Acceptable
	default:
		return 5 // Fast
	}
}

// calculateReviewInterval calculates the time since last review
func calculateReviewInterval(lastReviewed *time.Time) float64 {
	if lastReviewed == nil {
		return DefaultReviewInterval
	}

	delta := time.Since(*lastReviewed)
	return delta.Minutes()
}

// calculateSolutionMultiplier determines the multiplier based on solution quality
func calculateSolutionMultiplier(averageScore float64) float64 {
	if averageScore == PerfectScore {
		return OptimalMultiplier
	}
	return (PerfectScore - averageScore) + 1
}

// calculateTimeFactor converts time taken into a factor for SR score calculation
func calculateTimeFactor(timeTaken int) float64 {
	switch {
	case timeTaken == UnsolvedTimeValue:
		return UnsolvedTimeMultiplier
	case timeTaken < FastSolutionThreshold:
		return float64(timeTaken) * FastTimeMultiplier
	case timeTaken < MediumSolutionThreshold:
		return float64(timeTaken) * MediumTimeMultiplier
	case timeTaken < SlowSolutionThreshold:
		return float64(timeTaken) * SlowTimeMultiplier
	default:
		return float64(timeTaken) * VerySlowTimeMultiplier
	}
}
