package entity

type Lesson struct {
	ID   int32  `json:"id"`
	Name string `json:"name" validate:"required,min=2,max=20"`
}
