package controllers

import (
	"bluebell/logic"
	"bluebell/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// gin.Context 是一个非常重要的结构体，
// 它承载着 HTTP 请求的上下文信息，
// 包含了请求、响应、路由、以及与请求生命周期相关的一些方法。

// SignupHandler 处理注册请求的函数
func SignupHandler(c *gin.Context) {
	// 获取参数和参数校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误，返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		// 判断err是不是validator类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}

		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	// 业务处理
	if err := logic.SignUp(c.Request.Context(), p); err != nil {
		zap.L().Error("logic.Signup failed", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
		return
	}
	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}

func LoginHandler(c *gin.Context) {
	// 获取请求参数及参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误，返回响应
		zap.L().Error("Login with invalid param", zap.Error(err))
		// 判断err是不是validator类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg": removeTopStruct(errs.Translate(trans)), // 翻译错误
		})
		return
	}

	// 业务逻辑处理
	user, err := logic.Login(c.Request.Context(), p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"msg": "用户名或密码错误",
		})
		return
	}

	// 返回响应
	ResponseSuccess(c, gin.H{
		"user_id": fmt.Sprintf("%d", user.UserID),   // id值大于1<<53-1 int64类型的最大值是1<<63-1
		"user_name": user.Username,
		"token": user.Token,
	})

}
