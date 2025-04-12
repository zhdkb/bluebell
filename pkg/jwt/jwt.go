package jwt

import (
	"errors"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const TokenExpireDuration = time.Hour * 24 * 365

var mySecret = []byte("夏天夏天悄悄过去")

type MyClaims struct {
	UserID		int64	`json:"user_id"`
	Username	string	`json:"username"`
	Tokentype	string	`json:"tokentype"`
	jwt.RegisteredClaims
}

// GenToken 生成JWT
func GenToken(userID int64, username string, tokentype string, validTime time.Duration) (string, error) {
	// 创建一个我们自己的声明的数据
	c := MyClaims{
		UserID: userID,
		Username: username,
		Tokentype: tokentype,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(validTime)),
			Issuer:    "bluebell",
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

// 生成双token
func GenDoubleToken(userID int64, username string) (string, string, error) {
	var accessToken, refreshToken string
	errapichan := make(chan error)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		accessTokenstr, err := GenToken(userID, username, "access", time.Hour * 24 * 7)
		if err != nil {
			errapichan <- err
			return
		}
		accessToken = accessTokenstr
	} ()

	go func ()  {
		defer wg.Done()
		refreshTokenstr, err := GenToken(userID, username, "refresh", TokenExpireDuration)
		if err != nil {
			errapichan <- err
			return
		}
		refreshToken = refreshTokenstr
	} ()

	go func ()  {
		wg.Wait()
		close(errapichan)
	} ()

	for err := range errapichan {
		if err != nil {
			return "", "", err
		}
	}

	return accessToken, refreshToken, nil
	
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid { // 校验token
		return mc, nil
	}
	return nil, errors.New("invalid token")
}
