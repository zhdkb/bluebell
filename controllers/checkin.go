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
