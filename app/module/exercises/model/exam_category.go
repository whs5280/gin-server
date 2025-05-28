package model

type ExamCategory struct {
	BaseModel
	Name        string `gorm:"type:varchar(50);not null" json:"name"`
	Level       string `gorm:"type:varchar(50);not null" json:"level"`
	Description string `gorm:"type:text;not null" json:"description"`
}

// GetCategoryList 获取分类列表
func GetCategoryList() (categories []ExamCategory, err error) {
	err = DB.Find(&categories).Error
	return categories, err
}

func (ExamCategory) TableName() string {
	return "exam_category"
}
