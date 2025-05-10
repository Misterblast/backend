package repo

import (
	"context"

	questionEntity "github.com/ghulammuzz/misterblast/internal/question/entity"
)

func (r *questionRepository) ListQuestionTypes(ctx context.Context) ([]questionEntity.QuestionType, error) {
	levels := []struct {
		code string
		id   string
		en   string
	}{
		{"c1", "C1", "C1"},
		{"c2", "C2", "C2"},
		{"c3", "C3", "C3"},
		{"c4", "C4", "C4"},
		{"c5", "C5", "C5"},
		{"c6", "C6", "C6"},
	}

	dimensions := []struct {
		code string
		id   string
		en   string
	}{
		{"faktual", "Faktual", "Factual"},
		{"konseptual", "Konseptual", "Conceptual"},
		{"prosedural", "Prosedural", "Procedural"},
		{"metakognitif", "Metakognitif", "Metacognitive"},
	}

	var questionTypes []questionEntity.QuestionType
	for _, level := range levels {
		for _, dim := range dimensions {
			questionTypes = append(questionTypes, questionEntity.QuestionType{
				IndonesianName: level.id + " " + dim.id,
				EnglishName:    level.en + " " + dim.en,
				Code:           level.code + "_" + dim.code,
			})
		}
	}

	return questionTypes, nil
}
