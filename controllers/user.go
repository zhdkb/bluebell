package controllers

import (
	"bluebell/logic"
	"bluebell/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SignupHandler 处理注册请求的函数
func SignupHandler(c *gin.Context) {
	// 获取参数和参数校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误，返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"msg": "请求参数有误",
		})
		return
	}
	// 手动对请求参数进行详细的业务规则校验
	if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.RePassword != p.Password {
		zap.L().Error("SignUp with invalid param")
		c.JSON(http.StatusOK, gin.H{
			"msg": "请求参数有误",
		})
		return
	}
	fmt.Println(p)
	
	// 业务处理
	logic.SignUp(p)
	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}
