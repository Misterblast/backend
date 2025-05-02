package entity

type SetQuestion struct {
	Number      int    `json:"number" validate:"required,min=1"`
	Type        string `json:"type" validate:"required,oneof=c1_faktual c1_konseptual c1_prosedural c1_metakognitif c2_faktual c2_konseptual c2_prosedural c2_metakognitif c3_faktual c3_konseptual c3_prosedural c3_metakognitif c4_faktual c4_konseptual c4_prosedural c4_metakognitif c5_faktual c5_konseptual c5_prosedural c5_metakognitif c6_faktual c6_konseptual c6_prosedural c6_metakognitif"`
	Format      string `json:"format" validate:"required,oneof=mm t/f sa mc4 mcx essay"`
	Content     string `json:"content" validate:"required"`
	Explanation string `json:"explanation" validate:"required"`
	SetID       int32  `json:"set_id" validate:"required"`
}

type EditQuestion struct {
	Number      int    `json:"number" validate:"required,min=1"`
	Type        string `json:"type" validate:"required,oneof=c1_faktual c1_konseptual c1_prosedural c1_metakognitif c2_faktual c2_konseptual c2_prosedural c2_metakognitif c3_faktual c3_konseptual c3_prosedural c3_metakognitif c4_faktual c4_konseptual c4_prosedural c4_metakognitif c5_faktual c5_konseptual c5_prosedural c5_metakognitif c6_faktual c6_konseptual c6_prosedural c6_metakognitif"`
	Format      string `json:"format" validate:"required,oneof=mm t/f sa mc4 mcx essay"`
	Content     string `json:"content" validate:"required"`
	IsQuiz      bool   `json:"is_quiz"`
	Explanation string `json:"explanation" validate:"required"`
	SetID       int32  `json:"set_id" validate:"required"`
}

type ListQuestionExample struct {
	ID          int32  `json:"id"`
	Number      int    `json:"number"`
	Type        string `json:"type"`
	Format      string `json:"format"`
	Content     string `json:"content"`
	Explanation string `json:"explanation"`
	SetID       int32  `json:"set_id"`
}

type DetailQuestionExample struct {
	ID          int32        `json:"id"`
	Number      int          `json:"number"`
	Type        string       `json:"type"`
	Format      string       `json:"format"`
	Content     string       `json:"content"`
	Explanation string       `json:"explanation"`
	SetID       int32        `json:"set_id"`
	Answers     []ListAnswer `json:"answers"`
}

type ListQuestionQuiz struct {
	ID      int32        `json:"id"`
	Number  int          `json:"number"`
	Type    string       `json:"type"`
	Format  string       `json:"format"`
	Content string       `json:"content"`
	SetID   int32        `json:"set_id"`
	Answers []ListAnswer `json:"answers"`
}

type ListQuestionAdmin struct {
	ID          int32  `json:"id"`
	Number      int    `json:"number"`
	Type        string `json:"type"`
	Format      string `json:"format"`
	Content     string `json:"content"`
	IsQuiz      bool   `json:"is_quiz"`
	SetID       int32  `json:"set_id"`
	SetName     string `json:"set_name"`
	LessonName  string `json:"lesson_name"`
	ClassName   string `json:"class_name"`
	Explanation string `json:"explanation"`
}
