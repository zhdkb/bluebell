package controllers

import (
	"bluebell/logic"
	"bluebell/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

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
		
	// 业务处理
	if err := logic.SignUp(p); err != nil {
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
