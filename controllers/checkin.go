package controllers

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CheckInHandler 处理用户签到。
func CheckInHandler(c *gin.Context) {
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	data, err := logic.CheckIn(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, mysql.ErrorAlreadyCheckIn) {
			ResponseError(c, CodeAlreadyCheckIn)
			return
		}
		zap.L().Error("logic.CheckIn failed", zap.Int64("userID", userID), zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, data)
}

// CheckInMonthHandler 查询用户某个月的签到记录。
func CheckInMonthHandler(c *gin.Context) {
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	month := c.Query("month")
	data, err := logic.GetMonthlyCheckIn(c.Request.Context(), userID, month)
	if err != nil {
		zap.L().Error("logic.GetMonthlyCheckIn failed", zap.Int64("userID", userID), zap.String("month", month), zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	ResponseSuccess(c, data)
}
