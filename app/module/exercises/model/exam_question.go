package model

type ExamQuestion struct {
	BaseModel
	CategoryId   int                  `gorm:"type:int(11);not null" json:"category_id"`
	QuestionType int                  `gorm:"type:tinyint(4);not null" json:"question_type"`
	Content      string               `gorm:"type:text;not null" json:"content"`
	Year         int                  `gorm:"type:smallint(6);default:null" json:"year"`
	CreatedAt    string               `gorm:"type:datetime;" json:"created_at"`
	Options      []ExamQuestionOption `gorm:"foreignkey:QuestionId" json:"options"`
}

type ExamQuestionReq struct {
	CategoryId string `form:"category_id" binding:"required"`
	Page       int32  `form:"page,default=1"`
	PageSize   int32  `form:"page_size,default=5"`
}

// GetQuestionByCategoryId 获取列表
func GetQuestionByCategoryId(req ExamQuestionReq) (examQuestion []ExamQuestion, err error) {
	query := DB.Where("category_id = ?", req.CategoryId).Preload("Options").Order("RAND()")

	if req.Page > 0 {
		query.Limit(int(req.PageSize)).Offset(int(req.PageSize * (req.Page - 1)))
	}

	err = query.Find(&examQuestion).Error
	return examQuestion, err
}

func (ExamQuestion) TableName() string {
	return "exam_question"
}
