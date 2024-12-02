package routes

import (
	"bluebell/controllers"
	"bluebell/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.New()

	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 注册业务路由
	r.POST("/signup", controllers.SignupHandler)

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	return r
}
