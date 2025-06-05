package controller

import (
	"gin-server/app/module/exercises/helper"
	"gin-server/app/module/exercises/model"
	"gin-server/app/module/exercises/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CategoryIndex(g *gin.Context) {
	examService := service.ExaminationService{G: g}

	list, err := examService.GetCategoryList()
	if err != nil {
		helper.ResponseJson(g, true, "获取列表失败", err, http.StatusInternalServerError)
		return
	}
	helper.ResponseJson(g, false, "获取列表成功", list)
}

func QuestionIndex(g *gin.Context) {
	examService := service.ExaminationService{G: g}

	var req model.ExamQuestionReq
	if err := g.ShouldBindQuery(&req); err != nil {
		helper.ResponseJson(g, true, "参数错误", err, http.StatusFailedDependency)
		return
	}

	list, err := examService.GetQuestionByCategoryId(req)
	if err != nil {
		helper.ResponseJson(g, true, "获取列表失败", err, http.StatusInternalServerError)
	}
	helper.ResponseJson(g, false, "获取列表成功", list)
}

func AddFav(g *gin.Context) {
	examService := service.ExaminationService{G: g}
	questionId, _ := strconv.Atoi(g.Query("question_id"))
	userId := helper.CommonGetUserId(g)

	isFav, _ := examService.IsFav(questionId, userId)
	if isFav {
		helper.ResponseJson(g, true, "已收藏", nil, http.StatusFailedDependency)
		return
	}

	err := examService.AddFav(questionId, userId)
	if err != nil {
		helper.ResponseJson(g, true, "添加失败", err, http.StatusInternalServerError)
	}
	helper.ResponseJson(g, false, "添加成功", nil)
}

func DelFav(g *gin.Context) {
	examService := service.ExaminationService{G: g}
	questionId, _ := strconv.Atoi(g.Query("question_id"))
	userId := helper.CommonGetUserId(g)

	isFav, _ := examService.IsFav(questionId, userId)
	if !isFav {
		helper.ResponseJson(g, true, "未收藏", nil, http.StatusFailedDependency)
		return
	}

	err := examService.DelFav(questionId, userId)
	if err != nil {
		helper.ResponseJson(g, true, "删除失败", err, http.StatusInternalServerError)
	}
	helper.ResponseJson(g, false, "删除成功", nil)
}

func FavList(g *gin.Context) {
	examService := service.ExaminationService{G: g}

	var req model.ExamQuestionFavReq
	if err := g.ShouldBindQuery(&req); err != nil {
		helper.ResponseJson(g, true, "参数错误", err, http.StatusFailedDependency)
		return
	}

	userId := helper.CommonGetUserId(g)
	req.UserId = userId

	list, err := examService.GetFavList(req)
	if err != nil {
		helper.ResponseJson(g, true, "获取列表失败", err, http.StatusInternalServerError)
	}
	helper.ResponseJson(g, false, "获取列表成功", list)
}
