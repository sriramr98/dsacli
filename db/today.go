package db

import (
	"dsacli/types"
	"time"
)

func InsertTodayQuestions(questions []types.Question) error {
	today := time.Now()
	var todayQuestions []types.TodayQuestion

	for _, q := range questions {
		todayQuestions = append(todayQuestions, types.TodayQuestion{
			Date:       today.Format("2006-01-02"),
			QuestionID: q.Model.ID,
		})
	}

	res := gormDB.Create(&todayQuestions)
	return res.Error
}

func GetTodayQuestions() ([]types.Question, error) {
	var questions []types.Question
	today := time.Now().Format("2006-01-02")

	res := gormDB.Where("date = ?", today).Find(&questions)
	if res.Error != nil {
		return nil, res.Error
	}

	return questions, nil
}
