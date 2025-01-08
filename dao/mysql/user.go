package mysql

import (
	"bluebell/models"
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
)

// 把每一步数据库操作封装成函数
// 待logic层根据业务需求调用

const secret = "sunliulei"

// CheckUserExist 检查指定用户名是否存在
func CheckUserExist(username string) (err error) {
	if db == nil {
        return errors.New("database connection is nil")
    }
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	// db.Get(&count, sqlStr, username);
	if err = db.Raw(sqlStr, username).Scan(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户已存在")
	}
	return
}

// InsertUser 像数据库中插入一条新的用户记录
func InsertUser(ctx context.Context, user *models.User) (err error) {
	// 对密码进行加密
	user.Password = encryptPassword(user.Password)
	// 执行SQL语句入库
	sqlStr := `insert into user(user_id, username, password) values(?,?,?)`
	err = db.WithContext(ctx).Exec(sqlStr, user.UserID, user.Username, user.Password).Error
	return
}

func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

func Login(ctx context.Context, user *models.User) (err error) {
	oPassword := user.Password //用户登录的密码
	sqlStr := `select user_id, username, password from user where username = ?`
	// err = db.Get(user, sqlStr, user.Username)
	err = db.WithContext(ctx).Raw(sqlStr, user.Username).Scan(user).Error
	if err == sql.ErrNoRows {
		return errors.New("用户不存在")
	}
	if err != nil {
		// 查询数据库失败
		return err
	}
	// 判断密码是否正确
	password := encryptPassword(oPassword)
	if password != user.Password {
		return errors.New("密码错误")
	}
	return

}

// GetUserById 根据id获取用户信息
func GetUserById(ctx context.Context, uid int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id, username from user where user_id = ?`
	// err = db.Get(user, sqlStr, uid)
	err = db.WithContext(ctx).Raw(sqlStr, uid).Scan(user).Error
	return
}
