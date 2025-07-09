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

func FindQuestionByID(db *sql.DB, id int) (*types.Question, error) {
	var q types.Question
	err := db.QueryRow("SELECT id, name, url, difficulty, last_reviewed, sr_score, attempted FROM questions WHERE id = ?", id).
		Scan(&q.ID, &q.Name, &q.URL, &q.Difficulty, &q.LastReviewed, &q.SRScore, &q.Attempted)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &q, nil
}

func UpdateQuestion(db *sql.DB, q types.Question) error {
	updateSQL := `UPDATE questions SET name = ?, url = ?, difficulty = ?, last_reviewed = ?, sr_score = ?, attempted = ? WHERE id = ?`
	_, err := db.Exec(updateSQL, q.Name, q.URL, q.Difficulty, q.LastReviewed, q.SRScore, q.Attempted, q.ID)
	return err
}

func InsertQuestions(db *sql.DB, questions []types.Question) error {
	insertSQL := `INSERT INTO questions (id, name, url, difficulty, last_reviewed, sr_score, attempted) VALUES (?, ?, ?, ?, ?, ?, ?)`
	stmt, err := db.Prepare(insertSQL)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, q := range questions {
		_, err := stmt.Exec(q.ID, q.Name, q.URL, q.Difficulty, q.LastReviewed, q.SRScore, q.Attempted)
		if err != nil {
			return err
		}
	}

	return nil
}
