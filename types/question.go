package types

import "gorm.io/gorm"

type Question struct {
	gorm.Model
	Name         string  `json:"name"`
	URL          string  `json:"url"`
	Difficulty   string  `json:"difficulty"`
	LastReviewed *string `json:"last_reviewed"`
	SRScore      int     `json:"sr_score"`
	Attempted    bool    `json:"attempted"`
}

type TodayQuestion struct {
	gorm.Model
	QuestionID uint   `json:"question_id"`
	Date       string `json:"date"`
}
