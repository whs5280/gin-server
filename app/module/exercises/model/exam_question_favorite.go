package model

type ExamQuestionFav struct {
	BaseModel
	QuestionId int          `gorm:"type:int(11);not null" json:"question_id"`
	UserId     int          `gorm:"type:int(11);not null" json:"user_id"`
	Question   ExamQuestion `gorm:"foreignKey:QuestionId;references:ID" json:"question"`
}

type ExamQuestionFavReq struct {
	UserId   int `form:"user_id"`
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=5"`
}

func AddFav(questionId int, userId int) (err error) {
	err = DB.Create(&ExamQuestionFav{
		QuestionId: questionId,
		UserId:     userId,
	}).Error
	return err
}

func DelFav(questionId int, userId int) (err error) {
	err = DB.Where("question_id = ? and user_id = ?", questionId, userId).Delete(&ExamQuestionFav{}).Error
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

func GetFavList(req ExamQuestionFavReq) (favList []ExamQuestionFav, err error) {
	err = DB.Preload("Question").Preload("Question.Options").
		Where("user_id = ?", req.UserId).
		Offset(req.PageSize * (req.Page - 1)).
		Limit(req.PageSize).
		Find(&favList).
		Error
	return favList, err
}

func (ExamQuestionFav) TableName() string {
	return "exam_question_favorite"
}
