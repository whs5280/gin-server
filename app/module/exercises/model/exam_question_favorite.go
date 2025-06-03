package model

type ExamQuestionFav struct {
	BaseModel
	QuestionId int `gorm:"type:int(11);not null" json:"question_id"`
	UserId     int `gorm:"type:int(11);not null" json:"user_id"`
}

func AddFav(questionId int, userId int) (err error) {
	err = DB.Create(&ExamQuestionFav{
		QuestionId: questionId,
		UserId:     userId,
	}).Error
	return err
}

func IsFav(questionId int, userId int) (bool, error) {
	var count int64
	err := DB.Model(&ExamQuestionFav{}).Where("question_id = ? and user_id = ?", questionId, userId).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (ExamQuestionFav) TableName() string {
	return "exam_question_favorite"
}
