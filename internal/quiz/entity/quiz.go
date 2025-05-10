package entity

type AnswersQuizSubmit struct {
	Number int    `json:"number"`
	Answer string `json:"answer"`
}

type QuizSubmit struct {
	Answers []AnswersQuizSubmit `json:"answers"`
}

type ListQuizSubmission struct {
	ID          int    `json:"id"`
	SetID       int    `json:"set_id"`
	Correct     int    `json:"correct"`
	Grade       int    `json:"grade"`
	Lesson      string `json:"lesson"`
	Class       string `json:"class"`
	SubmittedAt int64  `json:"submitted_at"`
}

type ListQuizSubmissionAdmin struct {
	ID          int    `json:"id"`
	SetID       int    `json:"set_id"`
	Name        string `json:"name"`
	Correct     int    `json:"correct"`
	Grade       int    `json:"grade"`
	Lesson      string `json:"lesson"`
	Class       string `json:"class"`
	SubmittedAt int64  `json:"submitted_at"`
}

type QuizExp struct {
	ID          int          `json:"id"`
	Grade       string       `json:"grade"`
	SubmittedAt int64        `json:"submitted_at"`
	Correct     int          `json:"correct"`
	Wrong       int          `json:"wrong"`
	AttemptNo   int          `json:"attempt_no"`
	Answers     []QuizExpObj `json:"answers"`
}

type QuizExpObj struct {
	Number          int    `json:"number"`
	UserCode        string `json:"user_code"`
	ActualCode      string `json:"actual_code"`
	UserContent     string `json:"user_content"`
	ActualContent   string `json:"actual_content"`
	QuestionContent string `json:"question_content"`
	IsCorrect       bool   `json:"is_correct"`
	Explanation     string `json:"explanation"`
	Reason          string `json:"reason"`
	Format          string `format`
}
