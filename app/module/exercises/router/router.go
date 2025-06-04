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

	userGroup := router.Group("/user")
	{
		userGroup.GET("/register", controller.Register)
		userGroup.GET("/login", controller.Login)
		userGroup.GET("/logout", controller.Logout)
	}

	apiGroup := router.Group("/examination").
		Use(middleware.CheckToken)
	{
		apiGroup.GET("/category", controller.CategoryIndex)
		apiGroup.GET("/question", controller.QuestionIndex)
		apiGroup.GET("/addFav", controller.AddFav)
		apiGroup.GET("/delFav", controller.DelFav)
		apiGroup.GET("/favList", controller.FavList)
	}

	return router
}

// loadDefaultMiddleware
func loadDefaultMiddleware(router *gin.Engine) {
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.NoCache)
	router.Use(middleware.Options)
}
