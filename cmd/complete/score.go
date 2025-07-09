package complete

import (
	"dsacli/types"
	"math"
	"time"
)

func CalculateScore(timeTaken int, hintsNeeded int, optimalSolution int, anyBugs int, question types.Question) int {
	timeRank := calculateTimeRank(timeTaken)

	// Score is the average rating of these 4 parameters
	score := (float64(hintsNeeded) + timeRank + float64(optimalSolution) + float64(anyBugs)) / 4

	now := time.Now()
	cDate := 10000.0
	if question.LastReviewed != nil {
		lastReviewedDate, err := time.Parse(DateFormat, *question.LastReviewed)
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

	cTime := calculateTimeFactor(timeTaken)
	srScore := int(math.Round((cDate + cTime) * cSolution))

	if question.SRScore == 0 {
		return srScore
	} else {
		return int(float64(question.SRScore)*0.7 + float64(srScore)*0.3)
	}
}

func calculateTimeRank(timeTaken int) float64 {
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

	return timeRank
}

func calculateTimeFactor(timeTaken int) float64 {
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

	return cTime
}
