package model

import (
	"gin-server/app/module/exercises/pagination"
)

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
	CategoryId   int `form:"category_id" binding:"required"`
	QuestionType int `form:"question_type"`
	Page         int `form:"page,default=1"`
	PageSize     int `form:"page_size,default=5"`
}

type ExamQuestionResp struct {
	List       []ExamQuestion         `json:"list"`
	Pagination *pagination.Pagination `json:"pagination"`
}

// GetQuestionByCategoryId 获取列表
func GetQuestionByCategoryId(req ExamQuestionReq) (examQuestion []ExamQuestion, err error) {
	query := DB.Preload("Options").Where("category_id = ?", req.CategoryId)

	if req.QuestionType != 0 {
		query = query.Where("question_type = ?", req.QuestionType)
	}

	err = query.Order("RAND()").
		Offset(req.PageSize * (req.Page - 1)).
		Limit(req.PageSize).
		Find(&examQuestion).
		Error
	return examQuestion, err
}

func (ExamQuestion) TableName() string {
	return "exam_question"
}
