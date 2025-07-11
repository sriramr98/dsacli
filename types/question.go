package types

import "time"

type Question struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	Name         string     `json:"name"`
	URL          string     `json:"url"`
	Difficulty   string     `json:"difficulty"`
	LastReviewed *time.Time `json:"last_reviewed"`
	SRScore      int        `json:"sr_score"`
	Attempted    bool       `json:"attempted"`
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
