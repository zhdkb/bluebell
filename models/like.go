package models

import "time"

const (
	LikeActionLike   = "like"
	LikeActionUnlike = "unlike"
)

type ParamLikeData struct {
	PostID string `json:"post_id" binding:"required"`
	Action string `json:"action" binding:"required,oneof=like unlike"`
}

type LikeResult struct {
	PostID    string `json:"post_id"`
	UserID    string `json:"user_id"`
	Liked     bool   `json:"liked"`
	Changed   bool   `json:"changed"`
	Delta     int64  `json:"delta"`
	LikeCount int64  `json:"like_count"`
}

type PostLikeEvent struct {
	EventID int64 `json:"event_id,string" gorm:"column:event_id"`
	PostID  int64 `json:"post_id,string" gorm:"column:post_id"`
	UserID  int64 `json:"user_id,string" gorm:"column:user_id"`
	Liked   int8  `json:"liked" gorm:"column:liked"`
	Delta   int64 `json:"delta" gorm:"column:delta"`
}

type PostLike struct {
	PostID     int64     `gorm:"column:post_id;primaryKey"`
	UserID     int64     `gorm:"column:user_id;primaryKey"`
	Liked      int8      `gorm:"column:liked"`
	UpdateTime time.Time `gorm:"column:update_time"`
}

func (PostLike) TableName() string {
	return "post_like"
}

type FailedLikeEvent struct {
	ID         int64     `gorm:"column:id;primaryKey"`
	EventID    int64     `gorm:"column:event_id"`
	PostID     int64     `gorm:"column:post_id"`
	UserID     int64     `gorm:"column:user_id"`
	Liked      int8      `gorm:"column:liked"`
	Delta      int64     `gorm:"column:delta"`
	RetryCount int       `gorm:"column:retry_count"`
	CreateTime time.Time `gorm:"column:create_time"`
	UpdateTime time.Time `gorm:"column:update_time"`
}

func (FailedLikeEvent) TableName() string {
	return "like_event_failed"
}
