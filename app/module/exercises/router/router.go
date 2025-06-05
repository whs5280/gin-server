package router

import (
	"gin-server/app/module/exercises/controller"
	"gin-server/app/module/exercises/router/middleware"
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

	// API 路由组
	apiRouter := router.Group("/api")
	{
		userGroup := apiRouter.Group("/user")
		{
			userGroup.GET("/register", controller.Register)
			userGroup.GET("/login", controller.Login)
			userGroup.GET("/logout", controller.Logout)
		}
		examGroup := apiRouter.Group("/examination").
			Use(middleware.CheckToken)
		{
			examGroup.GET("/category", controller.CategoryIndex)
			examGroup.GET("/question", controller.QuestionIndex)
			examGroup.GET("/addFav", controller.AddFav)
			examGroup.GET("/delFav", controller.DelFav)
			examGroup.GET("/favList", controller.FavList)
		}
	}

	return router
}

// loadDefaultMiddleware
func loadDefaultMiddleware(router *gin.Engine) {
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.NoCache)
	router.Use(middleware.Options)
	router.Use(middleware.ErrorHandler())
}
