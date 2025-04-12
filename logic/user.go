package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
	"context"
	"errors"
	"time"

)

// 存放业务逻辑的代码

func SignUp(ctx context.Context, p *models.ParamSignUp) (err error) {
	// 判断用户存不存在
	if err = mysql.CheckUserExist(ctx, p.Username); err != nil {
		// 数据库查询出错
		return err
	}

	// 生成UID
	userID := snowflake.GenID()
	// 构造一个User实例
	user := &models.User{
		UserID:		userID,
		Username: 	p.Username,
		Password: 	p.Password,
	}
	// 保存进数据库
	return mysql.InsertUser(ctx, user)
	// redis.xxx
}

func Login(ctx context.Context, p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}

	// 传递的是指针，就能拿到userID
	if err := mysql.Login(ctx, user); err != nil {
		return nil, err
	}

	// 生成JWT
	accesstoken, refreshtoken, err := jwt.GenDoubleToken(user.UserID, user.Username)
	if err != nil {
		return
	}
	user.AccessToken = accesstoken
	user.RefreshToken = refreshtoken
	return
}

func Refresh(ctx context.Context, refreshtoken string) (string, string, error) {
	mc, err := jwt.ParseToken(refreshtoken)
	if err != nil {
		return "", "", err
	}

	if mc.Tokentype != "refresh" {
		return "", "", errors.New("token类型错误")
	}

	accesstoken, err := jwt.GenToken(mc.UserID, mc.Username, "access", time.Hour * 24 *7)
	if err != nil {
		return "", "", errors.New("access token 生成失败")
	}

	return accesstoken, refreshtoken, nil
}
