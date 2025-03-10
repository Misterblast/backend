package repo

import quizEntity "github.com/ghulammuzz/misterblast/internal/quiz/entity"

func (r *quizRepository) GetLast(userID int) (quizEntity.QuizExp, error) {
	var quizExp quizEntity.QuizExp
	// TODO: implement retrieval logic
	return quizExp, nil
}
