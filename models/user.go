package models

type User struct {
	UserID		int64	`gorm:"column:user_id"`
	Username	string	`gorm:"column:username"`
	Password	string	`gorm:"column:password"`
	AccessToken		string
	RefreshToken	string
}
