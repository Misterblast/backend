package entity

type Question struct {
	ID      int32  `json:"id"`
	Number  int    `json:"number"`
	Type    string `json:"type"`
	Content string `json:"content"`
	IsQuiz  bool   `json:"is_quiz"`
	SetID   int32  `json:"set_id"`
}
