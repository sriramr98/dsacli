package db

import (
	"dsacli/types"
	"gorm.io/gorm"
)

func GetQuestionsByDifficulty(difficulty string) ([]types.Question, error) {
	var question []types.Question
	res := gormDB.Find(&question, "difficulty = ?", difficulty).Order("id")
	if res.Error != nil {
		return nil, res.Error
	}

	return question, nil
}

func GetAllQuestions() ([]types.Question, error) {
	var questions []types.Question
	res := gormDB.Find(&questions)
	if res.Error != nil {
		return nil, res.Error
	}
	return questions, nil
}

func FindQuestionByID(id int) (types.Question, error) {
	q := types.Question{Model: gorm.Model{ID: uint(id)}}
	gormDB.First(&q)
	return q, nil
}

func UpdateQuestion(q types.Question) error {
	res := gormDB.Save(&q)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func InsertQuestions(questions []types.Question) error {
	res := gormDB.Create(questions)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
