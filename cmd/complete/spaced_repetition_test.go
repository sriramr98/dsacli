package complete

import (
	"dsacli/types"
	"testing"
)

func TestCalculatePScore(t *testing.T) {
	tests := []struct {
		name       string
		timeTaken  int
		hintsUsed  int
		optimality int
		bugs       int
		expected   float64
	}{
		{
			name:       "Perfect performance",
			timeTaken:  25,  // ≤ 30 minutes → 0.4
			hintsUsed:  0,   // No hints → 0.3
			optimality: 5,   // Very optimal → 0.15
			bugs:       5,   // No bugs → 0.15
			expected:   1.0, // 0.4 + 0.3 + 0.15 + 0.15
		},
		{
			name:       "Unsolved problem",
			timeTaken:  -1,    // Unsolved → 0.0
			hintsUsed:  3,     // 3 hints → 0.3/4 = 0.075
			optimality: 1,     // Not optimal → 0.0
			bugs:       1,     // Many bugs → 0.0
			expected:   0.075, // 0.0 + 0.075 + 0.0 + 0.0
		},
		{
			name:       "Medium performance",
			timeTaken:  35,     // 30-45 minutes → 0.2
			hintsUsed:  1,      // 1 hint → 0.3/2 = 0.15
			optimality: 3,      // Middle → (3-1)/4 * 0.15 = 0.075
			bugs:       4,      // Few bugs → (4-1)/4 * 0.15 = 0.1125
			expected:   0.5375, // 0.2 + 0.15 + 0.075 + 0.1125
		},
		{
			name:       "Slow performance",
			timeTaken:  50,     // > 45 minutes → 0.1
			hintsUsed:  2,      // 2 hints → 0.3/3 = 0.1
			optimality: 2,      // Low optimal → (2-1)/4 * 0.15 = 0.0375
			bugs:       3,      // Some bugs → (3-1)/4 * 0.15 = 0.075
			expected:   0.3125, // 0.1 + 0.1 + 0.0375 + 0.075
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculatePScore(tt.timeTaken, tt.hintsUsed, tt.optimality, tt.bugs)
			if result != tt.expected {
				t.Errorf("CalculatePScore() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestProcessReview(t *testing.T) {
	t.Run("Instant mastery on first perfect attempt", func(t *testing.T) {
		question := &types.Question{
			AttemptCount:   0,
			Mastered:       false,
			EasinessFactor: 2.5,
			ReviewStreak:   0,
		}

		// Perfect performance on first try
		ProcessReview(question, 20, 0, 5, 5) // p-score = 1.0

		if question.AttemptCount != 1 {
			t.Errorf("Expected AttemptCount = 1, got %d", question.AttemptCount)
		}

		if !question.Mastered {
			t.Errorf("Expected HasProgressionMastery = true for instant mastery")
		}

		if question.ReviewStreak != 1 {
			t.Errorf("Expected ReviewStreak = 1, got %d", question.ReviewStreak)
		}

		if question.ReviewInterval != 1 {
			t.Errorf("Expected ReviewInterval = 1, got %d", question.ReviewInterval)
		}
	})

	t.Run("Proven mastery on second good attempt", func(t *testing.T) {
		question := &types.Question{
			AttemptCount:   1,
			Mastered:       false,
			LastPScore:     0.9, // Previous attempt was good
			EasinessFactor: 2.5,
			ReviewStreak:   1,
			ReviewInterval: 1,
		}

		// Good performance on second try
		ProcessReview(question, 25, 0, 5, 5) // p-score = 1.0

		if question.AttemptCount != 2 {
			t.Errorf("Expected AttemptCount = 2, got %d", question.AttemptCount)
		}

		if !question.Mastered {
			t.Errorf("Expected HasProgressionMastery = true for proven mastery")
		}

		if question.ReviewStreak != 2 {
			t.Errorf("Expected ReviewStreak = 2, got %d", question.ReviewStreak)
		}

		if question.ReviewInterval != 6 {
			t.Errorf("Expected ReviewInterval = 6, got %d", question.ReviewInterval)
		}
	})

	t.Run("Failed recall resets streak", func(t *testing.T) {
		question := &types.Question{
			AttemptCount:   2,
			Mastered:       false,
			EasinessFactor: 2.8,
			ReviewStreak:   3,
			ReviewInterval: 15,
		}

		// Poor performance
		ProcessReview(question, -1, 5, 1, 1) // p-score = 0.0375

		if question.AttemptCount != 3 {
			t.Errorf("Expected AttemptCount = 3, got %d", question.AttemptCount)
		}

		if question.ReviewStreak != 0 {
			t.Errorf("Expected ReviewStreak = 0, got %d", question.ReviewStreak)
		}

		if question.ReviewInterval != 1 {
			t.Errorf("Expected ReviewInterval = 1, got %d", question.ReviewInterval)
		}

		if question.EasinessFactor < 2.59 || question.EasinessFactor > 2.61 {
			t.Errorf("Expected EasinessFactor ≈ 2.6, got %f", question.EasinessFactor)
		}
	})

	t.Run("Easiness factor minimum bound", func(t *testing.T) {
		question := &types.Question{
			EasinessFactor: 1.4,
			ReviewStreak:   0,
		}

		// Poor performance that would lower EF below 1.3
		ProcessReview(question, -1, 5, 1, 1)

		if question.EasinessFactor != 1.3 {
			t.Errorf("Expected EasinessFactor = 1.3 (minimum), got %f", question.EasinessFactor)
		}
	})
}

func TestCheckProgressionGate(t *testing.T) {
	t.Run("Empty question list", func(t *testing.T) {
		questions := make([]types.Question, 0)
		result := CheckProgressionGate(questions)
		if result {
			t.Errorf("Expected false for empty question list")
		}
	})

	t.Run("Less than 50% mastery", func(t *testing.T) {
		questions := []types.Question{
			{Mastered: true},
			{Mastered: false},
			{Mastered: false},
			{Mastered: false},
		}
		result := CheckProgressionGate(questions)
		if result {
			t.Errorf("Expected false for 25%% mastery")
		}
	})

	t.Run("Exactly 50% mastery", func(t *testing.T) {
		questions := []types.Question{
			{Mastered: true},
			{Mastered: true},
			{Mastered: false},
			{Mastered: false},
		}
		result := CheckProgressionGate(questions)
		if result {
			t.Errorf("Expected false for exactly 50%% mastery")
		}
	})

	t.Run("More than 50% mastery", func(t *testing.T) {
		questions := []types.Question{
			{Mastered: true},
			{Mastered: true},
			{Mastered: true},
			{Mastered: false},
		}
		result := CheckProgressionGate(questions)
		if !result {
			t.Errorf("Expected true for 75%% mastery")
		}
	})

	t.Run("100% mastery", func(t *testing.T) {
		questions := []types.Question{
			{Mastered: true},
			{Mastered: true},
			{Mastered: true},
		}
		result := CheckProgressionGate(questions)
		if !result {
			t.Errorf("Expected true for 100%% mastery")
		}
	})
}
