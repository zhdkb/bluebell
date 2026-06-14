package controllers

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetTopComments(c *gin.Context) {
	postID, err := parseQueryInt64(c, "post_id")
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	page, size := getPageInfo(c)

	data, err := logic.GetTopComments(c.Request.Context(), postID, page, size)
	if err != nil {
		zap.L().Error("logic.GetTopComments failed", zap.Int64("postID", postID), zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

func GetSubComments(c *gin.Context) {
	postID, err := parseQueryInt64(c, "post_id")
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	parentID, err := parseQueryInt64(c, "parent_id")
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	page, size := getPageInfo(c)

	data, err := logic.GetSubComments(c.Request.Context(), postID, parentID, page, size)
	if err != nil {
		zap.L().Error("logic.GetSubComments failed", zap.Int64("postID", postID), zap.Int64("parentID", parentID), zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

func GetCommentByCommentID(c *gin.Context) {
	commentID, err := parseQueryInt64(c, "comment_id")
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	data, err := logic.GetCommentByID(c.Request.Context(), commentID)
	if err != nil {
		handleCommentError(c, "logic.GetCommentByID failed", err)
		return
	}
	ResponseSuccess(c, data)
}

func CreateComment(c *gin.Context) {
	p := new(models.ParamCreateComment)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("CreateComment with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	if err := logic.CreateComment(c.Request.Context(), userID, p); err != nil {
		handleCommentError(c, "logic.CreateComment failed", err)
		return
	}
	ResponseSuccess(c, nil)
}

func UpdateComment(c *gin.Context) {
	p := new(models.ParamUpdateComment)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("UpdateComment with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	if err := logic.UpdateComment(c.Request.Context(), userID, p); err != nil {
		handleCommentError(c, "logic.UpdateComment failed", err)
		return
	}
	ResponseSuccess(c, nil)
}

func DeleteComment(c *gin.Context) {
	commentID, err := parseQueryInt64(c, "comment_id")
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	if err := logic.DeleteComment(c.Request.Context(), userID, commentID); err != nil {
		handleCommentError(c, "logic.DeleteComment failed", err)
		return
	}
	ResponseSuccess(c, nil)
}

func GetCommentCount(c *gin.Context) {
	postID, err := parseQueryInt64(c, "post_id")
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	count, err := logic.GetCommentCount(c.Request.Context(), postID)
	if err != nil {
		zap.L().Error("logic.GetCommentCount failed", zap.Int64("postID", postID), zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, count)
}

func GetTopCommentCount(c *gin.Context) {
	postID, err := parseQueryInt64(c, "post_id")
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	count, err := logic.GetTopCommentCount(c.Request.Context(), postID)
	if err != nil {
		zap.L().Error("logic.GetTopCommentCount failed", zap.Int64("postID", postID), zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, count)
}

func GetSubCommentCount(c *gin.Context) {
	parentID, err := parseQueryInt64(c, "parent_id")
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	count, err := logic.GetSubCommentCount(c.Request.Context(), parentID)
	if err != nil {
		zap.L().Error("logic.GetSubCommentCount failed", zap.Int64("parentID", parentID), zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, count)
}

func GetCommentCountByUserID(c *gin.Context) {
	userID, err := parseQueryInt64(c, "user_id")
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	count, err := logic.GetCommentCountByUserID(c.Request.Context(), userID)
	if err != nil {
		zap.L().Error("logic.GetCommentCountByUserID failed", zap.Int64("userID", userID), zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, count)
}

func parseQueryInt64(c *gin.Context, key string) (int64, error) {
	value := c.Query(key)
	if value == "" {
		return 0, strconv.ErrSyntax
	}
	return strconv.ParseInt(value, 10, 64)
}

func handleCommentError(c *gin.Context, msg string, err error) {
	zap.L().Error(msg, zap.Error(err))
	if errors.Is(err, mysql.ErrorInvalidID) {
		ResponseError(c, CodeInvalidParam)
		return
	}
	ResponseError(c, CodeServerBusy)
}
