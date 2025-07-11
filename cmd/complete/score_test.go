package complete

import (
	"dsacli/types"
	"testing"
	"time"
)

func TestCalculateScore(t *testing.T) {
	tests := []struct {
		name            string
		timeTaken       int
		hintsNeeded     int
		optimalSolution int
		anyBugs         int
		question        types.Question
		expectPositive  bool
	}{
		{
			name:            "Perfect solution",
			timeTaken:       20,
			hintsNeeded:     5,
			optimalSolution: 5,
			anyBugs:         5,
			question:        types.Question{SRScore: 0},
			expectPositive:  true,
		},
		{
			name:            "Unsolved question",
			timeTaken:       UnsolvedTimeValue,
			hintsNeeded:     1,
			optimalSolution: 1,
			anyBugs:         1,
			question:        types.Question{SRScore: 0},
			expectPositive:  true,
		},
		{
			name:            "With previous score",
			timeTaken:       30,
			hintsNeeded:     3,
			optimalSolution: 4,
			anyBugs:         4,
			question:        types.Question{SRScore: 100},
			expectPositive:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := CalculateScore(tt.timeTaken, tt.hintsNeeded, tt.optimalSolution, tt.anyBugs, tt.question)

			if tt.expectPositive && score <= 0 {
				t.Errorf("Expected positive score, got %d", score)
			}
		})
	}
}

func TestCalculateTimeRank(t *testing.T) {
	tests := []struct {
		name      string
		timeTaken int
		expected  float64
	}{
		{"Unsolved", UnsolvedTimeValue, 0},
		{"Very fast", 15, 5},
		{"Fast", 24, 5},
		{"Medium", 35, 3.5},
		{"Slow", 50, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rank := calculateTimeRank(tt.timeTaken)
			if rank != tt.expected {
				t.Errorf("Expected rank %f, got %f", tt.expected, rank)
			}
		})
	}
}

func TestCalculateReviewInterval(t *testing.T) {
	t.Run("No previous review", func(t *testing.T) {
		interval := calculateReviewInterval(nil)
		if interval != DefaultReviewInterval {
			t.Errorf("Expected default interval %f, got %f", DefaultReviewInterval, interval)
		}
	})

	t.Run("With previous review", func(t *testing.T) {
		pastTime := time.Now().Add(-time.Hour)
		interval := calculateReviewInterval(&pastTime)
		if interval <= 0 {
			t.Errorf("Expected positive interval, got %f", interval)
		}
	})
}

func TestCalculateSolutionMultiplier(t *testing.T) {
	tests := []struct {
		name         string
		averageScore float64
		expected     float64
	}{
		{"Perfect score", PerfectScore, OptimalMultiplier},
		{"Good score", 4.0, 2.0},
		{"Average score", 3.0, 3.0},
		{"Poor score", 1.0, 5.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			multiplier := calculateSolutionMultiplier(tt.averageScore)
			if multiplier != tt.expected {
				t.Errorf("Expected multiplier %f, got %f", tt.expected, multiplier)
			}
		})
	}
}

func TestCalculateTimeFactor(t *testing.T) {
	tests := []struct {
		name      string
		timeTaken int
		expected  float64
	}{
		{"Unsolved", UnsolvedTimeValue, UnsolvedTimeMultiplier},
		{"Very fast", 20, 20 * FastTimeMultiplier},
		{"Medium", 30, 30 * MediumTimeMultiplier},
		{"Slow", 40, 40 * SlowTimeMultiplier},
		{"Very slow", 60, 60 * VerySlowTimeMultiplier},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factor := calculateTimeFactor(tt.timeTaken)
			if factor != tt.expected {
				t.Errorf("Expected factor %f, got %f", tt.expected, factor)
			}
		})
	}
}

func TestCalculateScore_Integration(t *testing.T) {
	// Test with a realistic scenario
	lastReviewed := time.Now().Add(-24 * time.Hour)
	question := types.Question{
		ID:           1,
		Name:         "Two Sum",
		SRScore:      1000,
		LastReviewed: &lastReviewed,
	}

	// Good performance
	score := CalculateScore(20, 5, 5, 5, question)
	if score <= 0 {
		t.Errorf("Expected positive score for good performance, got %d", score)
	}

	// Poor performance
	poorScore := CalculateScore(UnsolvedTimeValue, 1, 1, 1, question)
	if poorScore <= score {
		t.Errorf("Expected poor performance score (%d) to be higher than good performance score (%d)", poorScore, score)
	}
}

func BenchmarkCalculateScore(b *testing.B) {
	question := types.Question{
		ID:      1,
		SRScore: 1000,
	}

	for i := 0; i < b.N; i++ {
		CalculateScore(25, 4, 4, 4, question)
	}
}
