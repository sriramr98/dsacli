package db

import (
	"dsacli/types"
	_ "embed"
)

const (
	DBFilename string = "dsacli.db"
	AppName    string = "dsacli"
)

type Database interface {
	GetQuestionsByDifficulty(difficulty string) ([]types.Question, error)
	GetAllQuestions() ([]types.Question, error)
	FindQuestionByID(id uint) (types.Question, error)
	UpdateQuestion(question types.Question) error
	InsertQuestions(questions []types.Question) error
	GetTodayQuestions() ([]types.Question, []types.TodayQuestion, error)
	InsertTodayQuestions(questions []types.Question) error
	GetTodayQuestionsWithStatus() ([]types.TodayQuestionWithStatus, error)
	MarkTodayQuestionCompleted(questionID uint) error
	GetAllAttemptedQuestions() ([]types.Question, error)
}
