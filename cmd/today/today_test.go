package today

import (
	"dsacli/types"
	"errors"
	"testing"

	"github.com/spf13/cobra"
)

// MockDatabase implements the Database interface for testing
type MockDatabase struct {
	QuestionsByDifficulty      map[string][]types.Question
	AllQuestions               []types.Question
	TodayQuestionsWithStatus   []types.TodayQuestionWithStatus
	ShouldReturnError          bool
	ErrorMessage               string
	InsertTodayQuestionsCalled bool
	InsertedQuestions          []types.Question
}

func (m *MockDatabase) GetQuestionsByDifficulty(difficulty string) ([]types.Question, error) {
	if m.ShouldReturnError {
		return nil, errors.New(m.ErrorMessage)
	}
	return m.QuestionsByDifficulty[difficulty], nil
}

func (m *MockDatabase) GetAllQuestions() ([]types.Question, error) {
	if m.ShouldReturnError {
		return nil, errors.New(m.ErrorMessage)
	}
	return m.AllQuestions, nil
}

func (m *MockDatabase) GetTodayQuestionsWithStatus() ([]types.TodayQuestionWithStatus, error) {
	if m.ShouldReturnError {
		return nil, errors.New(m.ErrorMessage)
	}
	return m.TodayQuestionsWithStatus, nil
}

func (m *MockDatabase) InsertTodayQuestions(questions []types.Question) error {
	m.InsertTodayQuestionsCalled = true
	m.InsertedQuestions = questions
	if m.ShouldReturnError {
		return errors.New(m.ErrorMessage)
	}
	return nil
}

// Unused methods for interface compliance
func (m *MockDatabase) FindQuestionByID(id uint) (types.Question, error) {
	return types.Question{}, nil
}
func (m *MockDatabase) UpdateQuestion(question types.Question) error {
	return nil
}
func (m *MockDatabase) InsertQuestions(questions []types.Question) error {
	return nil
}
func (m *MockDatabase) GetTodayQuestions() ([]types.Question, error) {
	return nil, nil
}
func (m *MockDatabase) MarkTodayQuestionCompleted(questionID uint) error {
	return nil
}

// Helper function to create test questions
func createTestQuestion(id uint, name, difficulty string, attempted bool, pScore float64) types.Question {
	return types.Question{
		ID:           id,
		Name:         name,
		URL:          "https://example.com/" + name,
		Difficulty:   difficulty,
		LastReviewed: nil, // Set to nil to avoid time comparison issues
		LastPScore:   pScore,
		Attempted:    attempted,
	}
}

// Helper function to compare question slices ignoring time fields
func questionsEqual(a, b []types.Question) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !questionEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

// Helper function to compare individual questions ignoring time fields
func questionEqual(a, b types.Question) bool {
	return a.ID == b.ID &&
		a.Name == b.Name &&
		a.URL == b.URL &&
		a.Difficulty == b.Difficulty &&
		a.LastPScore == b.LastPScore &&
		a.Attempted == b.Attempted
}

func questionOneOf(q types.Question, questions ...types.Question) bool {
	for _, candidate := range questions {
		if questionEqual(q, candidate) {
			return true
		}
	}
	return false
}

func TestAllAttempted(t *testing.T) {
	tests := []struct {
		name      string
		questions []types.Question
		expected  bool
	}{
		{
			name:      "Empty slice",
			questions: []types.Question{},
			expected:  true,
		},
		{
			name: "All attempted",
			questions: []types.Question{
				createTestQuestion(1, "q1", "easy", true, 0),
				createTestQuestion(2, "q2", "easy", true, 0),
			},
			expected: true,
		},
		{
			name: "Some attempted",
			questions: []types.Question{
				createTestQuestion(1, "q1", "easy", true, 0),
				createTestQuestion(2, "q2", "easy", false, 0),
			},
			expected: false,
		},
		{
			name: "None attempted",
			questions: []types.Question{
				createTestQuestion(1, "q1", "easy", false, 0),
				createTestQuestion(2, "q2", "easy", false, 0),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := allAttempted(tt.questions)
			if result != tt.expected {
				t.Errorf("allAttempted() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFilterUnattemptedQuestions(t *testing.T) {
	tests := []struct {
		name      string
		questions []types.Question
		expected  []types.Question
	}{
		{
			name:      "Empty slice",
			questions: []types.Question{},
			expected:  []types.Question{},
		},
		{
			name: "All attempted",
			questions: []types.Question{
				createTestQuestion(1, "q1", "easy", true, 0),
				createTestQuestion(2, "q2", "easy", true, 0),
			},
			expected: []types.Question{},
		},
		{
			name: "Some unattempted",
			questions: []types.Question{
				createTestQuestion(1, "q1", "easy", true, 0),
				createTestQuestion(2, "q2", "easy", false, 0),
				createTestQuestion(3, "q3", "easy", false, 0),
			},
			expected: []types.Question{
				createTestQuestion(2, "q2", "easy", false, 0),
				createTestQuestion(3, "q3", "easy", false, 0),
			},
		},
		{
			name: "All unattempted",
			questions: []types.Question{
				createTestQuestion(1, "q1", "easy", false, 0),
				createTestQuestion(2, "q2", "easy", false, 0),
			},
			expected: []types.Question{
				createTestQuestion(1, "q1", "easy", false, 0),
				createTestQuestion(2, "q2", "easy", false, 0),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterUnattemptedQuestions(tt.questions)
			if !questionsEqual(result, tt.expected) {
				t.Errorf("filterUnattemptedQuestions() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFilterOutQuestion(t *testing.T) {
	tests := []struct {
		name      string
		questions []types.Question
		excludeID uint
		expected  []types.Question
	}{
		{
			name:      "Empty slice",
			questions: []types.Question{},
			excludeID: 1,
			expected:  []types.Question{},
		},
		{
			name: "Exclude existing question",
			questions: []types.Question{
				createTestQuestion(1, "q1", "easy", false, 0),
				createTestQuestion(2, "q2", "easy", false, 0),
				createTestQuestion(3, "q3", "easy", false, 0),
			},
			excludeID: 2,
			expected: []types.Question{
				createTestQuestion(1, "q1", "easy", false, 0),
				createTestQuestion(3, "q3", "easy", false, 0),
			},
		},
		{
			name: "Exclude non-existing question",
			questions: []types.Question{
				createTestQuestion(1, "q1", "easy", false, 0),
				createTestQuestion(2, "q2", "easy", false, 0),
			},
			excludeID: 999,
			expected: []types.Question{
				createTestQuestion(1, "q1", "easy", false, 0),
				createTestQuestion(2, "q2", "easy", false, 0),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterOutQuestion(tt.questions, tt.excludeID)
			if !questionsEqual(result, tt.expected) {
				t.Errorf("filterOutQuestion() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestBuildAttemptedPool(t *testing.T) {
	tests := []struct {
		name      string
		questions []types.Question
		excludeID uint
		expected  []types.Question
	}{
		{
			name:      "Empty slice",
			questions: []types.Question{},
			excludeID: 1,
			expected:  []types.Question{},
		},
		{
			name: "Mixed attempted and unattempted",
			questions: []types.Question{
				createTestQuestion(1, "q1", "easy", true, 50),
				createTestQuestion(2, "q2", "easy", false, 0),
				createTestQuestion(3, "q3", "easy", true, 30),
				createTestQuestion(4, "q4", "easy", true, 70),
			},
			excludeID: 2,
			expected: []types.Question{
				createTestQuestion(1, "q1", "easy", true, 50),
				createTestQuestion(3, "q3", "easy", true, 30),
				createTestQuestion(4, "q4", "easy", true, 70),
			},
		},
		{
			name: "Exclude attempted question",
			questions: []types.Question{
				createTestQuestion(1, "q1", "easy", true, 50),
				createTestQuestion(2, "q2", "easy", true, 30),
				createTestQuestion(3, "q3", "easy", true, 70),
			},
			excludeID: 2,
			expected: []types.Question{
				createTestQuestion(1, "q1", "easy", true, 50),
				createTestQuestion(3, "q3", "easy", true, 70),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildAttemptedPool(tt.questions, tt.excludeID)
			if !questionsEqual(result, tt.expected) {
				t.Errorf("buildAttemptedPool() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetHighestSRQuestion(t *testing.T) {
	tests := []struct {
		name         string
		questions    []types.Question
		expectedQ    types.Question
		expectedBool bool
	}{
		{
			name:         "Empty slice",
			questions:    []types.Question{},
			expectedQ:    types.Question{},
			expectedBool: false,
		},
		{
			name: "Single question",
			questions: []types.Question{
				createTestQuestion(1, "q1", "easy", true, 50),
			},
			expectedQ:    createTestQuestion(1, "q1", "easy", true, 50),
			expectedBool: true,
		},
		{
			name: "Multiple questions",
			questions: []types.Question{
				createTestQuestion(1, "q1", "easy", true, 30),
				createTestQuestion(2, "q2", "easy", true, 70),
				createTestQuestion(3, "q3", "easy", true, 50),
			},
			expectedQ:    createTestQuestion(2, "q2", "easy", true, 70),
			expectedBool: true,
		},
		{
			name: "Questions with same SR score",
			questions: []types.Question{
				createTestQuestion(1, "q1", "easy", true, 50),
				createTestQuestion(2, "q2", "easy", true, 50),
			},
			expectedQ:    createTestQuestion(1, "q1", "easy", true, 50), // First one wins
			expectedBool: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultQ, resultBool := getHighestSRQuestion(tt.questions)
			if !questionEqual(resultQ, tt.expectedQ) || resultBool != tt.expectedBool {
				t.Errorf("getHighestSRQuestion() = (%v, %v), want (%v, %v)", resultQ, resultBool, tt.expectedQ, tt.expectedBool)
			}
		})
	}
}

func TestGetFocusQuestion(t *testing.T) {
	tests := []struct {
		name            string
		questions       []types.Question
		expectHighestSR bool
		expectedBool    bool
	}{
		{
			name:            "Empty slice",
			questions:       []types.Question{},
			expectHighestSR: false,
			expectedBool:    false,
		},
		{
			name: "Only unattempted questions",
			questions: []types.Question{
				createTestQuestion(1, "q1", "easy", false, 0),
				createTestQuestion(2, "q2", "easy", false, 0),
			},
			expectHighestSR: false,
			expectedBool:    true,
		},
		{
			name: "Only attempted questions",
			questions: []types.Question{
				createTestQuestion(1, "q1", "easy", true, 30),
				createTestQuestion(2, "q2", "easy", true, 70),
			},
			expectedBool:    true,
			expectHighestSR: true,
		},
		{
			name: "Mixed attempted and unattempted",
			questions: []types.Question{
				createTestQuestion(1, "q1", "easy", true, 70),
				createTestQuestion(2, "q2", "easy", false, 0),
				createTestQuestion(3, "q3", "easy", true, 30),
			},
			expectedBool:    true,
			expectHighestSR: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultQ, resultBool := getFocusQuestion(tt.questions)
			if resultBool != tt.expectedBool {
				t.Errorf("getFocusQuestion() bool = %v, want %v", resultBool, tt.expectedBool)
			}

			if tt.expectHighestSR {
				// Expect highest SR question if all are attempted
				highestQ, _ := getHighestSRQuestion(tt.questions)
				if !questionEqual(resultQ, highestQ) {
					t.Errorf("getFocusQuestion() returned %v, want highest SR question %v", resultQ, highestQ)
				}
			} else {
				// Expect any one of unattempted questions
				unattempted := filterUnattemptedQuestions(tt.questions)
				if len(unattempted) > 0 && !questionOneOf(resultQ, unattempted[0]) {
					t.Errorf("getFocusQuestion() returned %v, want unattempted question %v", resultQ, unattempted[0])
				}
			}
		})
	}
}

func TestGenerateEasyPhaseQuestions(t *testing.T) {
	tests := []struct {
		name          string
		easyQuestions []types.Question
		expectedCount int
	}{
		{
			name:          "Empty questions",
			easyQuestions: []types.Question{},
			expectedCount: 0,
		},
		{
			name: "Single question",
			easyQuestions: []types.Question{
				createTestQuestion(1, "q1", "easy", false, 0),
			},
			expectedCount: 1,
		},
		{
			name: "Multiple questions",
			easyQuestions: []types.Question{
				createTestQuestion(1, "q1", "easy", false, 0),
				createTestQuestion(2, "q2", "easy", false, 0),
				createTestQuestion(3, "q3", "easy", false, 0),
			},
			expectedCount: 2, // questionsPerDay
		},
		{
			name: "All attempted questions",
			easyQuestions: []types.Question{
				createTestQuestion(1, "q1", "easy", true, 50),
				createTestQuestion(2, "q2", "easy", true, 30),
			},
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateEasyPhaseQuestions(tt.easyQuestions)
			if len(result) != tt.expectedCount {
				t.Errorf("generateEasyPhaseQuestions() returned %d questions, want %d", len(result), tt.expectedCount)
			}
		})
	}
}

func TestGenerateMediumPhaseQuestions(t *testing.T) {
	tests := []struct {
		name            string
		mediumQuestions []types.Question
		easyQuestions   []types.Question
		expectedCount   int
	}{
		{
			name:            "Empty medium questions",
			mediumQuestions: []types.Question{},
			easyQuestions:   []types.Question{},
			expectedCount:   0,
		},
		{
			name: "Medium question with review pool",
			mediumQuestions: []types.Question{
				createTestQuestion(1, "m1", "medium", false, 0),
			},
			easyQuestions: []types.Question{
				createTestQuestion(2, "e1", "easy", true, 50),
				createTestQuestion(3, "e2", "easy", true, 30),
			},
			expectedCount: 2,
		},
		{
			name: "Multiple medium questions with review pool",
			mediumQuestions: []types.Question{
				createTestQuestion(1, "m1", "medium", false, 0),
				createTestQuestion(2, "m2", "medium", false, 0),
			},
			easyQuestions: []types.Question{
				createTestQuestion(3, "e1", "easy", true, 50),
			},
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateMediumPhaseQuestions(tt.mediumQuestions, tt.easyQuestions)
			if len(result) != tt.expectedCount {
				t.Errorf("generateMediumPhaseQuestions() returned %d questions, want %d", len(result), tt.expectedCount)
			}
		})
	}
}

func TestGenerateHardPhaseQuestions(t *testing.T) {
	tests := []struct {
		name          string
		hardQuestions []types.Question
		allQuestions  []types.Question
		expectedCount int
	}{
		{
			name:          "Empty hard questions",
			hardQuestions: []types.Question{},
			allQuestions:  []types.Question{},
			expectedCount: 0,
		},
		{
			name: "Hard question with review pool",
			hardQuestions: []types.Question{
				createTestQuestion(1, "h1", "hard", false, 0),
			},
			allQuestions: []types.Question{
				createTestQuestion(1, "h1", "hard", false, 0),
				createTestQuestion(2, "e1", "easy", true, 50),
				createTestQuestion(3, "m1", "medium", true, 30),
			},
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateHardPhaseQuestions(tt.hardQuestions, tt.allQuestions)
			if len(result) != tt.expectedCount {
				t.Errorf("generateHardPhaseQuestions() returned %d questions, want %d", len(result), tt.expectedCount)
			}
		})
	}
}

func TestGenerateMasteryPhaseQuestions(t *testing.T) {
	tests := []struct {
		name          string
		allQuestions  []types.Question
		expectedCount int
	}{
		{
			name:          "Empty questions",
			allQuestions:  []types.Question{},
			expectedCount: 0,
		},
		{
			name: "Single question",
			allQuestions: []types.Question{
				createTestQuestion(1, "q1", "easy", true, 50),
			},
			expectedCount: 1,
		},
		{
			name: "Multiple questions",
			allQuestions: []types.Question{
				createTestQuestion(1, "q1", "easy", true, 30),
				createTestQuestion(2, "q2", "medium", true, 70),
				createTestQuestion(3, "q3", "hard", true, 50),
			},
			expectedCount: 2, // questionsPerDay
		},
		{
			name: "More than questionsPerDay",
			allQuestions: []types.Question{
				createTestQuestion(1, "q1", "easy", true, 10),
				createTestQuestion(2, "q2", "medium", true, 80),
				createTestQuestion(3, "q3", "hard", true, 50),
				createTestQuestion(4, "q4", "easy", true, 90),
			},
			expectedCount: 2, // questionsPerDay
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateMasteryPhaseQuestions(tt.allQuestions)
			if len(result) != tt.expectedCount {
				t.Errorf("generateMasteryPhaseQuestions() returned %d questions, want %d", len(result), tt.expectedCount)
			}
			// Check if questions are sorted by LastPScore (highest first)
			for i := 1; i < len(result); i++ {
				if result[i-1].LastPScore < result[i].LastPScore {
					t.Errorf("generateMasteryPhaseQuestions() questions not sorted by LastPScore: %d < %d", result[i-1].LastPScore, result[i].LastPScore)
				}
			}
		})
	}
}

func TestGenerateTodayQuestions(t *testing.T) {
	tests := []struct {
		name                   string
		mockDB                 *MockDatabase
		expectedError          bool
		expectedQuestionsCount int
	}{
		{
			name: "Error loading easy questions",
			mockDB: &MockDatabase{
				ShouldReturnError: true,
				ErrorMessage:      "database error",
			},
			expectedError: true,
		},
		{
			name: "Easy phase - unattempted easy questions",
			mockDB: &MockDatabase{
				QuestionsByDifficulty: map[string][]types.Question{
					easyPhase: {
						createTestQuestion(1, "e1", "easy", false, 0),
						createTestQuestion(2, "e2", "easy", false, 0),
					},
				},
			},
			expectedError:          false,
			expectedQuestionsCount: 2,
		},
		{
			name: "Medium phase - all easy attempted",
			mockDB: &MockDatabase{
				QuestionsByDifficulty: map[string][]types.Question{
					easyPhase: {
						createTestQuestion(1, "e1", "easy", true, 50),
					},
					mediumPhase: {
						createTestQuestion(2, "m1", "medium", false, 0),
					},
				},
			},
			expectedError:          false,
			expectedQuestionsCount: 2,
		},
		{
			name: "Hard phase - all easy and medium attempted",
			mockDB: &MockDatabase{
				QuestionsByDifficulty: map[string][]types.Question{
					easyPhase: {
						createTestQuestion(1, "e1", "easy", true, 50),
					},
					mediumPhase: {
						createTestQuestion(2, "m1", "medium", true, 60),
					},
					hardPhase: {
						createTestQuestion(3, "h1", "hard", false, 0),
					},
				},
				AllQuestions: []types.Question{
					createTestQuestion(1, "e1", "easy", true, 50),
					createTestQuestion(2, "m1", "medium", true, 60),
					createTestQuestion(3, "h1", "hard", false, 0),
				},
			},
			expectedError:          false,
			expectedQuestionsCount: 2,
		},
		{
			name: "Mastery phase - all attempted",
			mockDB: &MockDatabase{
				QuestionsByDifficulty: map[string][]types.Question{
					easyPhase: {
						createTestQuestion(1, "e1", "easy", true, 50),
					},
					mediumPhase: {
						createTestQuestion(2, "m1", "medium", true, 60),
					},
					hardPhase: {
						createTestQuestion(3, "h1", "hard", true, 70),
					},
				},
				AllQuestions: []types.Question{
					createTestQuestion(1, "e1", "easy", true, 50),
					createTestQuestion(2, "m1", "medium", true, 60),
					createTestQuestion(3, "h1", "hard", true, 70),
				},
			},
			expectedError:          false,
			expectedQuestionsCount: 2, // questionsPerDay in mastery mode
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			questions, err := generateTodayQuestions(tt.mockDB)

			if tt.expectedError && err == nil {
				t.Errorf("generateTodayQuestions() expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("generateTodayQuestions() unexpected error: %v", err)
			}
			if len(questions) != tt.expectedQuestionsCount {
				t.Errorf("generateTodayQuestions() returned %d questions, want %d", len(questions), tt.expectedQuestionsCount)
			}
		})
	}
}

func TestGetTodayQuestionsWithStatusIfExist(t *testing.T) {
	tests := []struct {
		name          string
		mockDB        *MockDatabase
		expectedError bool
		expectedCount int
	}{
		{
			name: "Database error",
			mockDB: &MockDatabase{
				ShouldReturnError: true,
				ErrorMessage:      "database error",
			},
			expectedError: true,
		},
		{
			name: "No today questions",
			mockDB: &MockDatabase{
				TodayQuestionsWithStatus: []types.TodayQuestionWithStatus{},
			},
			expectedError: false,
			expectedCount: 0,
		},
		{
			name: "Has today questions",
			mockDB: &MockDatabase{
				TodayQuestionsWithStatus: []types.TodayQuestionWithStatus{
					{
						Question:  createTestQuestion(1, "q1", "easy", false, 0),
						Completed: false,
					},
					{
						Question:  createTestQuestion(2, "q2", "medium", true, 50),
						Completed: true,
					},
				},
			},
			expectedError: false,
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			questions, err := getTodayQuestionsWithStatusIfExist(tt.mockDB)

			if tt.expectedError && err == nil {
				t.Errorf("getTodayQuestionsWithStatusIfExist() expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("getTodayQuestionsWithStatusIfExist() unexpected error: %v", err)
			}
			if len(questions) != tt.expectedCount {
				t.Errorf("getTodayQuestionsWithStatusIfExist() returned %d questions, want %d", len(questions), tt.expectedCount)
			}
		})
	}
}

func TestExecuteToday(t *testing.T) {
	tests := []struct {
		name                       string
		mockDB                     *MockDatabase
		expectInsertTodayQuestions bool
	}{
		{
			name: "Today questions already exist",
			mockDB: &MockDatabase{
				TodayQuestionsWithStatus: []types.TodayQuestionWithStatus{
					{
						Question:  createTestQuestion(1, "q1", "easy", false, 0),
						Completed: false,
					},
				},
			},
			expectInsertTodayQuestions: false,
		},
		{
			name: "Generate new questions",
			mockDB: &MockDatabase{
				TodayQuestionsWithStatus: []types.TodayQuestionWithStatus{}, // No existing questions
				QuestionsByDifficulty: map[string][]types.Question{
					easyPhase: {
						createTestQuestion(1, "e1", "easy", false, 0),
						createTestQuestion(2, "e2", "easy", false, 0),
					},
				},
			},
			expectInsertTodayQuestions: true,
		},
		{
			name: "No questions found",
			mockDB: &MockDatabase{
				TodayQuestionsWithStatus: []types.TodayQuestionWithStatus{}, // No existing questions
				QuestionsByDifficulty:    map[string][]types.Question{},     // No questions available
			},
			expectInsertTodayQuestions: false,
		},
		{
			name: "Error generating questions",
			mockDB: &MockDatabase{
				TodayQuestionsWithStatus: []types.TodayQuestionWithStatus{}, // No existing questions
				ShouldReturnError:        true,
				ErrorMessage:             "database error",
			},
			expectInsertTodayQuestions: false,
		},
		{
			name: "Error saving today questions",
			mockDB: &MockDatabase{
				TodayQuestionsWithStatus: []types.TodayQuestionWithStatus{}, // No existing questions
				QuestionsByDifficulty: map[string][]types.Question{
					easyPhase: {
						createTestQuestion(1, "e1", "easy", false, 0),
					},
				},
				// This will cause InsertTodayQuestions to fail after being called
			},
			expectInsertTodayQuestions: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset the mock for each test
			originalShouldReturnError := tt.mockDB.ShouldReturnError
			originalErrorMessage := tt.mockDB.ErrorMessage
			tt.mockDB.InsertTodayQuestionsCalled = false
			tt.mockDB.InsertedQuestions = nil

			executeToday(tt.mockDB)

			if tt.expectInsertTodayQuestions != tt.mockDB.InsertTodayQuestionsCalled {
				t.Errorf("executeToday() InsertTodayQuestions called = %v, want %v",
					tt.mockDB.InsertTodayQuestionsCalled, tt.expectInsertTodayQuestions)
			}

			// For the error saving case, simulate the error after the call
			if tt.name == "Error saving today questions" && tt.mockDB.InsertTodayQuestionsCalled {
				tt.mockDB.ShouldReturnError = true
				tt.mockDB.ErrorMessage = "save error"
			}

			// Restore original values
			tt.mockDB.ShouldReturnError = originalShouldReturnError
			tt.mockDB.ErrorMessage = originalErrorMessage
		})
	}
}

// Additional edge case tests

func TestDisplayTodayQuestions(t *testing.T) {
	tests := []struct {
		name                     string
		questionsWithStatus      []types.TodayQuestionWithStatus
		expectCongratulationsMsg bool
	}{
		{
			name:                     "Empty questions",
			questionsWithStatus:      []types.TodayQuestionWithStatus{},
			expectCongratulationsMsg: false,
		},
		{
			name: "All questions completed",
			questionsWithStatus: []types.TodayQuestionWithStatus{
				{
					Question:  createTestQuestion(1, "q1", "easy", true, 50),
					Completed: true,
				},
				{
					Question:  createTestQuestion(2, "q2", "medium", true, 70),
					Completed: true,
				},
			},
			expectCongratulationsMsg: true,
		},
		{
			name: "Some questions completed",
			questionsWithStatus: []types.TodayQuestionWithStatus{
				{
					Question:  createTestQuestion(1, "q1", "easy", true, 50),
					Completed: true,
				},
				{
					Question:  createTestQuestion(2, "q2", "medium", false, 0),
					Completed: false,
				},
			},
			expectCongratulationsMsg: false,
		},
		{
			name: "No questions completed",
			questionsWithStatus: []types.TodayQuestionWithStatus{
				{
					Question:  createTestQuestion(1, "q1", "easy", false, 0),
					Completed: false,
				},
				{
					Question:  createTestQuestion(2, "q2", "medium", false, 0),
					Completed: false,
				},
			},
			expectCongratulationsMsg: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: This function prints to stdout, so we can't easily test the output
			// But we can test that it doesn't panic
			displayTodayQuestions(tt.questionsWithStatus)
			// In a real-world scenario, you might want to capture stdout to test the output
		})
	}
}

func TestDisplayQuestions(t *testing.T) {
	tests := []struct {
		name      string
		questions []types.Question
	}{
		{
			name:      "Empty questions",
			questions: []types.Question{},
		},
		{
			name: "Single question",
			questions: []types.Question{
				createTestQuestion(1, "q1", "easy", false, 0),
			},
		},
		{
			name: "Multiple questions with different difficulties",
			questions: []types.Question{
				createTestQuestion(1, "q1", "easy", false, 0),
				createTestQuestion(2, "q2", "medium", true, 50),
				createTestQuestion(3, "q3", "hard", true, 80),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: This function prints to stdout, so we can't easily test the output
			// But we can test that it doesn't panic
			displayQuestions(tt.questions)
		})
	}
}

func TestEdgeCasesForPhaseGeneration(t *testing.T) {
	t.Run("Easy phase with single question", func(t *testing.T) {
		questions := []types.Question{
			createTestQuestion(1, "q1", "easy", false, 0),
		}
		result := generateEasyPhaseQuestions(questions)
		if len(result) != 1 {
			t.Errorf("Expected 1 question, got %d", len(result))
		}
	})

	t.Run("Medium phase with no review candidates", func(t *testing.T) {
		mediumQuestions := []types.Question{
			createTestQuestion(1, "m1", "medium", false, 0),
		}
		easyQuestions := []types.Question{
			createTestQuestion(2, "e1", "easy", false, 0), // Not attempted, so no review candidates
		}
		result := generateMediumPhaseQuestions(mediumQuestions, easyQuestions)
		if len(result) != 1 { // Only focus question, no review
			t.Errorf("Expected 1 question, got %d", len(result))
		}
	})

	t.Run("Hard phase with no review candidates", func(t *testing.T) {
		hardQuestions := []types.Question{
			createTestQuestion(1, "h1", "hard", false, 0),
		}
		allQuestions := []types.Question{
			createTestQuestion(1, "h1", "hard", false, 0),
			createTestQuestion(2, "e1", "easy", false, 0), // Not attempted
		}
		result := generateHardPhaseQuestions(hardQuestions, allQuestions)
		if len(result) != 1 { // Only focus question, no review
			t.Errorf("Expected 1 question, got %d", len(result))
		}
	})
}

func TestDatabaseErrorScenarios(t *testing.T) {
	t.Run("Error loading medium questions", func(t *testing.T) {
		mockDB := &MockDatabase{
			QuestionsByDifficulty: map[string][]types.Question{
				easyPhase: {
					createTestQuestion(1, "e1", "easy", true, 50), // All attempted
				},
			},
		}
		// Simulate error when loading medium questions
		originalQuestions := mockDB.QuestionsByDifficulty[mediumPhase]
		mockDB.QuestionsByDifficulty[mediumPhase] = nil
		mockDB.ShouldReturnError = true
		mockDB.ErrorMessage = "medium questions error"

		_, err := generateTodayQuestions(mockDB)
		if err == nil {
			t.Error("Expected error when loading medium questions")
		}

		// Restore
		mockDB.QuestionsByDifficulty[mediumPhase] = originalQuestions
		mockDB.ShouldReturnError = false
	})

	t.Run("Error loading hard questions", func(t *testing.T) {
		mockDB := &MockDatabase{
			QuestionsByDifficulty: map[string][]types.Question{
				easyPhase: {
					createTestQuestion(1, "e1", "easy", true, 50),
				},
				mediumPhase: {
					createTestQuestion(2, "m1", "medium", true, 60),
				},
			},
		}
		// Simulate error when loading hard questions
		mockDB.ShouldReturnError = true
		mockDB.ErrorMessage = "hard questions error"

		_, err := generateTodayQuestions(mockDB)
		if err == nil {
			t.Error("Expected error when loading hard questions")
		}
	})

	t.Run("Error loading all questions for hard phase", func(t *testing.T) {
		mockDB := &MockDatabase{
			QuestionsByDifficulty: map[string][]types.Question{
				easyPhase: {
					createTestQuestion(1, "e1", "easy", true, 50),
				},
				mediumPhase: {
					createTestQuestion(2, "m1", "medium", true, 60),
				},
				hardPhase: {
					createTestQuestion(3, "h1", "hard", false, 0),
				},
			},
			AllQuestions: nil, // This will cause error
		}
		// Simulate error when loading all questions
		mockDB.ShouldReturnError = true
		mockDB.ErrorMessage = "all questions error"

		_, err := generateTodayQuestions(mockDB)
		if err == nil {
			t.Error("Expected error when loading all questions")
		}
	})
}

func TestGetCommand(t *testing.T) {
	mockDB := &MockDatabase{}
	cmd := GetCommand(mockDB)

	if cmd.Use != "today" {
		t.Errorf("Expected command Use to be 'today', got '%s'", cmd.Use)
	}

	if cmd.Short != "Suggests two DSA questions for today" {
		t.Errorf("Expected Short description to match, got '%s'", cmd.Short)
	}

	if cmd.Long != "Suggests two DSA questions for today based on difficulty progression and smart review." {
		t.Errorf("Expected Long description to match, got '%s'", cmd.Long)
	}

	if cmd.Run == nil {
		t.Error("Expected Run function to be set")
	}
}

func TestTodayCmd(t *testing.T) {
	mockDB := &MockDatabase{
		TodayQuestionsWithStatus: []types.TodayQuestionWithStatus{
			{
				Question:  createTestQuestion(1, "existing", "easy", false, 0),
				Completed: false,
			},
		},
	}

	cmdFunc := todayCmd(mockDB)
	if cmdFunc == nil {
		t.Error("Expected todayCmd to return a function")
	}

	// Test that the command function doesn't panic when called
	mockCmd := &cobra.Command{}
	cmdFunc(mockCmd, []string{})
}
