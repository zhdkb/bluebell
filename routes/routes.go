package routes

import (
	"bluebell/controllers"
	"bluebell/logger"
	"bluebell/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

// gin.Context 是一个非常重要的结构体，
// 它承载着 HTTP 请求的上下文信息，
// 包含了请求、响应、路由、以及与请求生命周期相关的一些方法。
func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // gin设置成发布模式
	}

	r := gin.New()

	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 注册业务路由
	r.POST("/signup", controllers.SignupHandler)
	r.POST("/login", controllers.LoginHandler)

	r.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
		// 如果是登录的用户，判断请求头中是否有 有效的JWT
		c.String(http.StatusOK, "pong")
	})
	return r
}

