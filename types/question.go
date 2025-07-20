package types

import "time"

type Question struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	Name         string     `json:"name"`
	URL          string     `json:"url"`
	Difficulty   string     `json:"difficulty"`
	LastReviewed *time.Time `json:"last_reviewed"`
	Attempted    bool       `json:"attempted"`

	// Spaced Repetition Algorithm fields
	ReviewInterval int     `json:"review_interval" gorm:"default:0"`   // days until next review
	EasinessFactor float64 `json:"easiness_factor" gorm:"default:2.5"` // interval growth factor
	ReviewStreak   int     `json:"review_streak" gorm:"default:0"`     // consecutive successful recalls
	Mastered       bool    `json:"mastered" gorm:"default:false"`      // progression gate flag
	AttemptCount   int     `json:"attempt_count" gorm:"default:0"`     // total attempts
	LastPScore     float64 `json:"last_p_score" gorm:"default:0"`      // previous attempt's p-score
}

type TodayQuestion struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	QuestionID uint   `json:"question_id"`
	Date       string `json:"date"`
	Completed  bool   `json:"completed" gorm:"default:false"`
}

// TodayQuestionWithStatus represents a question for today with its completion status
type TodayQuestionWithStatus struct {
	Question  Question `json:"question"`
	Completed bool     `json:"completed"`
}
