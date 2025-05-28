package controller

import (
	"gin-server/app/module/exercises/helper"
	"gin-server/app/module/exercises/model"
	"gin-server/app/module/exercises/service"
	"github.com/gin-gonic/gin"
)

func CategoryIndex(g *gin.Context) {
	examService := service.ExaminationService{G: g}

	list, err := examService.GetCategoryList()
	if err != nil {
		helper.ResponseJson(g, true, "获取列表失败", err, 500)
	}
	helper.ResponseJson(g, false, "获取列表成功", list)
}

func QuestionIndex(g *gin.Context) {
	examService := service.ExaminationService{G: g}

	var req model.ExamQuestionReq
	if err := g.ShouldBindQuery(&req); err != nil {
		helper.ResponseJson(g, true, "参数错误", err, 422)
		return
	}

	list, err := examService.GetQuestionByCategoryId(req)
	if err != nil {
		helper.ResponseJson(g, true, "获取列表失败", err, 500)
	}
	helper.ResponseJson(g, false, "获取列表成功", list)
}
