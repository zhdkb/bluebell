package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"context"
	"strings"
)

func GetTopComments(ctx context.Context, postID, page, size int64) ([]*models.ApiComment, error) {
	return mysql.GetTopComments(ctx, postID, page, size)
}

func GetSubComments(ctx context.Context, postID, parentID, page, size int64) ([]*models.ApiComment, error) {
	return mysql.GetSubComments(ctx, postID, parentID, page, size)
}

func GetCommentByID(ctx context.Context, commentID int64) (*models.ApiComment, error) {
	return mysql.GetCommentByID(ctx, commentID)
}

func CreateComment(ctx context.Context, userID int64, p *models.ParamCreateComment) error {
	p.Content = strings.TrimSpace(p.Content)
	if p.Content == "" {
		return mysql.ErrorInvalidID
	}

	post, err := mysql.GetPostById(ctx, p.PostID)
	if err != nil {
		return err
	}
	if post.ID == 0 {
		return mysql.ErrorInvalidID
	}

	user, err := mysql.GetUserById(ctx, userID)
	if err != nil {
		return err
	}
	if user.UserID == 0 {
		return mysql.ErrorInvalidID
	}

	parentID, replyID, err := normalizeCommentRelation(ctx, p.PostID, p.ParentID, p.ReplyID)
	if err != nil { 
		return err
	}

	commentID := snowflake.GenID()
	comment := &models.Comment{
		CommentID:  commentID,
		PostID:     p.PostID,
		AuthorID:   userID,
		AuthorName: user.Username,
		Content:    p.Content,
		Status:     1,
	}
	relation := &models.CommentRelation{
		PostID:    p.PostID,
		CommentID: commentID,
		ParentID:  parentID,
		ReplyID:   replyID,
	}
	return mysql.CreateComment(ctx, comment, relation)
}

func UpdateComment(ctx context.Context, userID int64, p *models.ParamUpdateComment) error {
	p.Content = strings.TrimSpace(p.Content)
	if p.Content == "" {
		return mysql.ErrorInvalidID
	}
	return mysql.UpdateComment(ctx, p.CommentID, userID, p.Content)
}

func DeleteComment(ctx context.Context, userID, commentID int64) error {
	return mysql.DeleteComment(ctx, commentID, userID)
}

func GetCommentCount(ctx context.Context, postID int64) (int64, error) {
	return mysql.GetCommentCount(ctx, postID)
}

func GetTopCommentCount(ctx context.Context, postID int64) (int64, error) {
	return mysql.GetTopCommentCount(ctx, postID)
}

func GetSubCommentCount(ctx context.Context, parentID int64) (int64, error) {
	return mysql.GetSubCommentCount(ctx, parentID)
}

func GetCommentCountByUserID(ctx context.Context, userID int64) (int64, error) {
	return mysql.GetCommentCountByUserID(ctx, userID)
}

func normalizeCommentRelation(ctx context.Context, postID, parentID, replyID int64) (int64, int64, error) {
	if parentID == 0 {
		return 0, 0, nil
	}

	parentRelation, err := mysql.GetCommentRelationByID(ctx, parentID)
	if err != nil {
		return 0, 0, err
	}
	if parentRelation.PostID != postID || parentRelation.ParentID != 0 {
		return 0, 0, mysql.ErrorInvalidID
	}

	if replyID == 0 {
		replyID = parentID
	}

	replyComment, err := mysql.GetCommentByID(ctx, replyID)
	if err != nil {
		return 0, 0, err
	}
	if replyComment.PostID != postID {
		return 0, 0, mysql.ErrorInvalidID
	}
	if replyComment.CommentID != parentID && replyComment.ParentID != parentID {
		return 0, 0, mysql.ErrorInvalidID
	}

	return parentID, replyID, nil
}
