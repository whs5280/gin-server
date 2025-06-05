package model

import (
	"gin-server/app/module/exercises/pagination"
)

// 1-单选 2-多选 3-判断 4-案例 5-论文
const (
	QuestionTypeSingle = 1
	QuestionTypeMulti  = 2
	QuestionTypeJudge  = 3
	QuestionTypeCase   = 4
	QuestionTypePaper  = 5
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
	CategoryId int `form:"category_id" binding:"required"`
	Type       int `form:"type"` // 1 上午题 2 下午题
	Page       int `form:"page,default=1"`
	PageSize   int `form:"page_size,default=5"`
}

type ExamQuestionResp struct {
	List       []ExamQuestion         `json:"list"`
	Pagination *pagination.Pagination `json:"pagination"`
}

// GetQuestionByCategoryId 获取列表
func GetQuestionByCategoryId(req ExamQuestionReq) (examQuestion []ExamQuestion, err error) {
	query := DB.Preload("Options").Where("category_id = ?", req.CategoryId)

	if req.Type != 0 {
		if req.Type == 1 { // 上午题
			query = query.Where("question_type IN (?)", []int8{QuestionTypeSingle, QuestionTypeMulti, QuestionTypeJudge})
		}
		if req.Type == 2 { // 下午题
			query = query.Where("question_type IN (?)", []int8{QuestionTypeCase, QuestionTypePaper})
		}
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
