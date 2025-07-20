package db

import (
	"dsacli/types"
	"time"
)

func (d SQLDatabase) InsertTodayQuestions(questions []types.Question) error {
	today := time.Now()
	var todayQuestions []types.TodayQuestion

	for _, q := range questions {
		todayQuestions = append(todayQuestions, types.TodayQuestion{
			Date:       today.Format("2006-01-02"),
			QuestionID: q.ID,
			Completed:  false,
		})
	}

	res := d.db.Create(&todayQuestions)
	return res.Error
}

func (d SQLDatabase) GetTodayQuestions() ([]types.Question, []types.TodayQuestion, error) {
	var questions []types.TodayQuestion
	today := time.Now().Format("2006-01-02")

	res := d.db.Where("date = ?", today).Find(&questions)
	if res.Error != nil {
		return nil, nil, res.Error
	}

	var questionIDs []uint
	for _, tq := range questions {
		questionIDs = append(questionIDs, tq.QuestionID)
	}

	var result []types.Question
	res = d.db.Where("id IN ?", questionIDs).Find(&result)
	if res.Error != nil {
		return nil, nil, res.Error
	}

	return result, questions, nil
}

// GetTodayQuestionsWithStatus returns today's questions along with their completion status
func (d SQLDatabase) GetTodayQuestionsWithStatus() ([]types.TodayQuestionWithStatus, error) {
	var todayQuestions []types.TodayQuestion
	today := time.Now().Format("2006-01-02")

	res := d.db.Where("date = ?", today).Find(&todayQuestions)
	if res.Error != nil {
		return nil, res.Error
	}

	if len(todayQuestions) == 0 {
		return nil, nil
	}

	var questionIDs []uint
	for _, tq := range todayQuestions {
		questionIDs = append(questionIDs, tq.QuestionID)
	}

	var questions []types.Question
	res = d.db.Where("id IN ?", questionIDs).Find(&questions)
	if res.Error != nil {
		return nil, res.Error
	}

	// Create a map for quick lookup of completion status
	completionMap := make(map[uint]bool)
	for _, tq := range todayQuestions {
		completionMap[tq.QuestionID] = tq.Completed
	}

	// Build the result with completion status
	var result []types.TodayQuestionWithStatus
	for _, q := range questions {
		result = append(result, types.TodayQuestionWithStatus{
			Question:  q,
			Completed: completionMap[q.ID],
		})
	}

	return result, nil
}

// MarkTodayQuestionCompleted marks a specific question as completed for today
func (d SQLDatabase) MarkTodayQuestionCompleted(questionID uint) error {
	today := time.Now().Format("2006-01-02")

	res := d.db.Model(&types.TodayQuestion{}).
		Where("date = ? AND question_id = ?", today, questionID).
		Update("completed", true)

	return res.Error
}
