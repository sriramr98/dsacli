package db

import (
	"dsacli/types"
)

func (d SQLDatabase) GetQuestionsByDifficulty(difficulty string) ([]types.Question, error) {
	var question []types.Question
	res := d.db.Find(&question, "difficulty = ?", difficulty).Order("id")
	if res.Error != nil {
		return nil, res.Error
	}

	return question, nil
}

func (d SQLDatabase) GetAllQuestions() ([]types.Question, error) {
	var questions []types.Question
	res := d.db.Find(&questions)
	if res.Error != nil {
		return nil, res.Error
	}
	return questions, nil
}

func (d SQLDatabase) FindQuestionByID(id uint) (types.Question, error) {
	q := types.Question{ID: id}
	d.db.First(&q)
	return q, nil
}

func (d SQLDatabase) UpdateQuestion(q types.Question) error {
	res := d.db.Save(&q)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (d SQLDatabase) InsertQuestions(questions []types.Question) error {
	res := d.db.Create(questions)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (d SQLDatabase) GetAllAttemptedQuestions() ([]types.Question, error) {
	var questions []types.Question
	res := d.db.Where("attempted = ?", true).Find(&questions)
	if res.Error != nil {
		return nil, res.Error
	}
	return questions, nil
}
