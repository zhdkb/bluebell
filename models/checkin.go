package models

import "time"

// CheckInDetail 记录用户每一次签到的明细。
type CheckInDetail struct {
	ID        int64     `json:"-" gorm:"column:id"`
	CheckInID int64     `json:"checkin_id,string" gorm:"column:checkin_id"`
	UserID    int64     `json:"user_id,string" gorm:"column:user_id"`
	SignDate  time.Time `json:"sign_date" gorm:"column:sign_date"`
	SignTime  time.Time `json:"sign_time" gorm:"column:sign_time"`
}

func (CheckInDetail) TableName() string {
	return "user_checkin_detail"
}

// CheckInCount 统计用户累计签到和连续签到信息。
type CheckInCount struct {
	ID              int64     `json:"-" gorm:"column:id"`
	UserID          int64     `json:"user_id,string" gorm:"column:user_id"`
	TotalCount      int       `json:"total_count" gorm:"column:total_count"`
	ContinuousCount int       `json:"continuous_count" gorm:"column:continuous_count"`
	LastSignDate    time.Time `json:"last_sign_date" gorm:"column:last_sign_date"`
}

func (CheckInCount) TableName() string {
	return "user_checkin_count"
}

type CheckInResult struct {
	CheckInID       int64  `json:"checkin_id,string"`
	UserID          int64  `json:"user_id,string"`
	SignDate        string `json:"sign_date"`
	SignTime        string `json:"sign_time"`
	TotalCount      int    `json:"total_count"`
	ContinuousCount int    `json:"continuous_count"`
}

type CheckInMonthResult struct {
	UserID      int64    `json:"user_id,string"`
	Month       string   `json:"month"`
	CheckInDays []string `json:"checkin_days"`
	TotalDays   int      `json:"total_days"`
}
