package db

import (
	"database/sql"
	"dsacli/types"
)

func GetQuestionsByDifficulty(db *sql.DB, difficulty string) ([]types.Question, error) {
	rows, err := db.Query("SELECT id, name, url, difficulty, last_reviewed, sr_score, attempted FROM questions WHERE difficulty = ? ORDER BY id", difficulty)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []types.Question
	for rows.Next() {
		var q types.Question
		err := rows.Scan(&q.ID, &q.Name, &q.URL, &q.Difficulty, &q.LastReviewed, &q.SRScore, &q.Attempted)
		if err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}

	return questions, rows.Err()
}

func GetAllQuestions(db *sql.DB) ([]types.Question, error) {
	rows, err := db.Query("SELECT id, name, url, difficulty, last_reviewed, sr_score, attempted FROM questions ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []types.Question
	for rows.Next() {
		var q types.Question
		err := rows.Scan(&q.ID, &q.Name, &q.URL, &q.Difficulty, &q.LastReviewed, &q.SRScore, &q.Attempted)
		if err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}

	return questions, rows.Err()
}
