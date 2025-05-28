package service

import (
	"fmt"
	"gin-server/app/module/exercises/model"
	"github.com/gin-gonic/gin"
	"log"
)

type ExaminationService struct {
	G      *gin.Context
	logger *log.Logger
}

func (s *ExaminationService) GetCategoryList() ([]model.ExamCategory, error) {
	categories, err := model.GetCategoryList()
	if err != nil {
		s.logger.Printf("获取分类列表失败: %v", err)
		return nil, fmt.Errorf("获取分类列表失败: %w", err)
	}

	if len(categories) == 0 {
		return []model.ExamCategory{}, nil // 返回空切片而非nil
	}

	return categories, nil
}

func (s *ExaminationService) GetQuestionByCategoryId(req model.ExamQuestionReq) ([]model.ExamQuestion, error) {
	questions, err := model.GetQuestionByCategoryId(req)
	if err != nil {
		s.logger.Printf("获取题库列表失败: %v", err)
		return nil, fmt.Errorf("获取题库列表失败: %w", err)
	}

	if len(questions) == 0 {
		return []model.ExamQuestion{}, nil // 返回空切片而非nil
	}

	return questions, nil
}
