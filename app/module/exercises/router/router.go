package router

import (
	"gin-server/app/module/exercises/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() *gin.Engine {
	router := gin.New()

	loadDefaultMiddleware(router)

	// 404 Handler.
	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	router.GET("/", controller.Index)

	apiGroup := router.Group("/examination")
	{
		apiGroup.GET("/category", controller.CategoryIndex)
		apiGroup.GET("/question", controller.QuestionIndex)
	}

	return router
}

// loadDefaultMiddleware
func loadDefaultMiddleware(router *gin.Engine) {
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
}
