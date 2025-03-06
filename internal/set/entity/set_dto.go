package entity

type SetSet struct {
	Name     string `json:"name" validate:"required,min=2,max=20"`
	IsQuiz   bool   `json:"is_quiz"`
	LessonID int32  `json:"lesson_id" validate:"required"`
	ClassID  int32  `json:"class_id" validate:"required"`
}

type ListSet struct {
	ID     int32  `json:"id"`
	Name   string `json:"name"`
	Lesson string `json:"lesson"`
	Class  string `json:"class"`
	IsQuiz bool   `json:"is_quiz"`
}
