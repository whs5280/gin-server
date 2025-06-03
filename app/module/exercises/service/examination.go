package service

import (
	"fmt"
	"gin-server/app/module/exercises/model"
	"gin-server/app/module/exercises/pagination"
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

func (s *ExaminationService) GetQuestionByCategoryId(req model.ExamQuestionReq) ([]model.ExamQuestionResp, error) {
	questions, err := model.GetQuestionByCategoryId(req)
	if err != nil {
		s.logger.Printf("获取题库列表失败: %v", err)
		return nil, fmt.Errorf("获取题库列表失败: %w", err)
	}

	if len(questions) == 0 {
		return []model.ExamQuestionResp{}, nil
	}

	questionResp := new(model.ExamQuestionResp)
	questionResp.List = questions
	questionResp.Pagination = pagination.MakePagination(questions, req.Page, req.PageSize)

	return []model.ExamQuestionResp{*questionResp}, nil
}

func (s *ExaminationService) IsFav(questionId int, userId int) (bool, error) {
	return model.IsFav(questionId, userId)
}

func (s *ExaminationService) AddFav(questionId int, userId int) (err error) {
	return model.AddFav(questionId, userId)
}

func (s *ExaminationService) DelFav(questionId int, userId int) (err error) {
	return model.DelFav(questionId, userId)
}

func (s *ExaminationService) GetFavList(req model.ExamQuestionFavReq) ([]model.ExamQuestionFav, error) {
	return model.GetFavList(req)
}
