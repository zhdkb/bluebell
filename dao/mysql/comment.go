package mysql

import (
	"bluebell/models"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

func GetTopComments(ctx context.Context, postID, page, size int64) ([]*models.ApiComment, error) {
	comments := make([]*models.ApiComment, 0)
	sqlStr := `select c.comment_id, c.post_id, c.author_id, c.author_name, c.content,
					r.parent_id, r.reply_id, c.create_time, c.update_time
				from comment c
				inner join comment_relation r on c.comment_id = r.comment_id
				where c.post_id = ?
					and c.status = 1
					and c.delete_time = 0
					and r.delete_time = 0
					and r.parent_id = 0
				order by c.create_time desc
				limit ?, ?`
	err := db.WithContext(ctx).Raw(sqlStr, postID, (page-1)*size, size).Scan(&comments).Error
	return comments, err
}

func GetSubComments(ctx context.Context, postID, parentID, page, size int64) ([]*models.ApiComment, error) {
	comments := make([]*models.ApiComment, 0)
	sqlStr := `select c.comment_id, c.post_id, c.author_id, c.author_name, c.content,
					r.parent_id, r.reply_id, c.create_time, c.update_time
				from comment c
				inner join comment_relation r on c.comment_id = r.comment_id
				where c.post_id = ?
					and r.parent_id = ?
					and c.status = 1
					and c.delete_time = 0
					and r.delete_time = 0
				order by c.create_time desc
				limit ?, ?`
	err := db.WithContext(ctx).Raw(sqlStr, postID, parentID, (page-1)*size, size).Scan(&comments).Error
	return comments, err
}

func GetCommentByID(ctx context.Context, commentID int64) (*models.ApiComment, error) {
	comment := new(models.ApiComment)
	sqlStr := `select c.comment_id, c.post_id, c.author_id, c.author_name, c.content,
					r.parent_id, r.reply_id, c.create_time, c.update_time
				from comment c
				inner join comment_relation r on c.comment_id = r.comment_id
				where c.comment_id = ?
					and c.status = 1
					and c.delete_time = 0
					and r.delete_time = 0`
	if err := db.WithContext(ctx).Raw(sqlStr, commentID).Scan(comment).Error; err != nil {
		return nil, err
	}
	if comment.CommentID == 0 {
		return nil, ErrorInvalidID
	}
	return comment, nil
}

func GetCommentRelationByID(ctx context.Context, commentID int64) (*models.CommentRelation, error) {
	relation := new(models.CommentRelation)
	sqlStr := `select post_id, comment_id, parent_id, reply_id, create_time, update_time
				from comment_relation
				where comment_id = ? and delete_time = 0`
	if err := db.WithContext(ctx).Raw(sqlStr, commentID).Scan(relation).Error; err != nil {
		return nil, err
	}
	if relation.CommentID == 0 {
		return nil, ErrorInvalidID
	}
	return relation, nil
}

func CreateComment(ctx context.Context, comment *models.Comment, relation *models.CommentRelation) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		sqlStrCreateComment := `insert into comment
					(comment_id, post_id, author_id, author_name, content, status)
					values (?, ?, ?, ?, ?, ?)`
		if err := tx.Exec(
			sqlStrCreateComment,
			comment.CommentID,
			comment.PostID,
			comment.AuthorID,
			comment.AuthorName,
			comment.Content,
			comment.Status,
		).Error; err != nil {
			return err
		}

		sqlStrCreateRelation := `insert into comment_relation
					(post_id, comment_id, parent_id, reply_id)
					values (?, ?, ?, ?)`
		return tx.Exec(
			sqlStrCreateRelation,
			relation.PostID,
			relation.CommentID,
			relation.ParentID,
			relation.ReplyID,
		).Error
	})
}

func UpdateComment(ctx context.Context, commentID, authorID int64, content string) error {
	result := db.WithContext(ctx).Exec(
		`update comment
			set content = ?
			where comment_id = ?
				and author_id = ?
				and status = 1
				and delete_time = 0`,
		content,
		commentID,
		authorID,
	)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrorInvalidID
	}
	return nil
}

func DeleteComment(ctx context.Context, commentID, authorID int64) error {
	now := time.Now().Unix()
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&models.Comment{}).
			Where("comment_id = ? and author_id = ? and status = 1 and delete_time = 0", commentID, authorID).
			Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return ErrorInvalidID
		}

		var deletedIDs []int64
		if err := tx.Model(&models.CommentRelation{}).
			Where("(comment_id = ? or parent_id = ? or reply_id = ?) and delete_time = 0", commentID, commentID, commentID).
			Pluck("comment_id", &deletedIDs).Error; err != nil {
			return err
		}
		if len(deletedIDs) == 0 {
			return ErrorInvalidID
		}

		if err := tx.Model(&models.Comment{}).
			Where("comment_id in ?", deletedIDs).
			Updates(map[string]interface{}{
				"status":      0,
				"delete_time": now,
			}).Error; err != nil {
			return err
		}

		return tx.Model(&models.CommentRelation{}).
			Where("comment_id in ?", deletedIDs).
			Update("delete_time", now).Error
	})
}

func GetCommentCount(ctx context.Context, postID int64) (int64, error) {
	var count int64
	sqlStr := `select count(*)
				from comment c
				inner join comment_relation r on c.comment_id = r.comment_id
				where c.post_id = ?
					and c.status = 1
					and c.delete_time = 0
					and r.delete_time = 0`
	err := db.WithContext(ctx).Raw(sqlStr, postID).Scan(&count).Error
	return count, err
}

func GetTopCommentCount(ctx context.Context, postID int64) (int64, error) {
	var count int64
	sqlStr := `select count(*)
				from comment c
				inner join comment_relation r on c.comment_id = r.comment_id
				where c.post_id = ?
					and c.status = 1
					and c.delete_time = 0
					and r.delete_time = 0
					and r.parent_id = 0`
	err := db.WithContext(ctx).Raw(sqlStr, postID).Scan(&count).Error
	return count, err
}

func GetSubCommentCount(ctx context.Context, parentID int64) (int64, error) {
	var count int64
	sqlStr := `select count(*)
				from comment c
				inner join comment_relation r on c.comment_id = r.comment_id
				where r.parent_id = ?
					and c.status = 1
					and c.delete_time = 0
					and r.delete_time = 0`
	err := db.WithContext(ctx).Raw(sqlStr, parentID).Scan(&count).Error
	return count, err
}

func GetCommentCountByUserID(ctx context.Context, userID int64) (int64, error) {
	var count int64
	sqlStr := `select count(*)
				from comment c
				inner join comment_relation r on c.comment_id = r.comment_id
				where c.author_id = ?
					and c.status = 1
					and c.delete_time = 0
					and r.delete_time = 0`
	err := db.WithContext(ctx).Raw(sqlStr, userID).Scan(&count).Error
	return count, err
}

func IsCommentNotFound(err error) bool {
	return errors.Is(err, ErrorInvalidID)
}
