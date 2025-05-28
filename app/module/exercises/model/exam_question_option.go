package model

type ExamQuestionOption struct {
	BaseModel
	QuestionId int    `gorm:"type:int(11);not null" json:"question_id"`
	OptionKey  string `gorm:"type:varchar(255);not null" json:"option_key"`
	Content    string `gorm:"type:text;not null" json:"content"`
	IsCorrect  int    `gorm:"type:tinyint(1);not null" json:"is_correct"`
}

func (ExamQuestionOption) TableName() string {
	return "exam_question_option"
}
