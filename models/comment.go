package models

import "time"

type Comment struct {
	ID         int64     `json:"-" gorm:"column:id"`
	CommentID  int64     `json:"comment_id,string" gorm:"column:comment_id"`
	PostID     int64     `json:"post_id,string" gorm:"column:post_id"`
	AuthorID   int64     `json:"author_id,string" gorm:"column:author_id"`
	AuthorName string    `json:"author_name" gorm:"column:author_name"`
	Content    string    `json:"content" gorm:"column:content"`
	Status     int32     `json:"status" gorm:"column:status"`
	CreateTime time.Time `json:"create_time" gorm:"column:create_time"`
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time"`
	DeleteTime int64     `json:"-" gorm:"column:delete_time"`
}

func (Comment) TableName() string {
	return "comment"
}

type CommentRelation struct {
	ID         int64     `json:"-" gorm:"column:id"`
	PostID     int64     `json:"post_id,string" gorm:"column:post_id"`
	CommentID  int64     `json:"comment_id,string" gorm:"column:comment_id"`
	ParentID   int64     `json:"parent_id,string" gorm:"column:parent_id"`
	ReplyID    int64     `json:"reply_id,string" gorm:"column:reply_id"`
	CreateTime time.Time `json:"create_time" gorm:"column:create_time"`
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time"`
	DeleteTime int64     `json:"-" gorm:"column:delete_time"`
}

func (CommentRelation) TableName() string {
	return "comment_relation"
}

type ApiComment struct {
	CommentID  int64     `json:"comment_id,string" gorm:"column:comment_id"`
	PostID     int64     `json:"post_id,string" gorm:"column:post_id"`
	AuthorID   int64     `json:"author_id,string" gorm:"column:author_id"`
	AuthorName string    `json:"author_name" gorm:"column:author_name"`
	Content    string    `json:"content" gorm:"column:content"`
	ParentID   int64     `json:"parent_id,string" gorm:"column:parent_id"`
	ReplyID    int64     `json:"reply_id,string" gorm:"column:reply_id"`
	CreateTime time.Time `json:"create_time" gorm:"column:create_time"`
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time"`
}

type ParamCreateComment struct {
	PostID   int64  `json:"post_id,string" binding:"required"`
	Content  string `json:"content" binding:"required"`
	ParentID int64  `json:"parent_id,string"`
	ReplyID  int64  `json:"reply_id,string"`
}

type ParamUpdateComment struct {
	CommentID int64  `json:"comment_id,string" binding:"required"`
	Content   string `json:"content" binding:"required"`
}
