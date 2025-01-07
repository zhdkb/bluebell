package mysql

import (
	"bluebell/settings"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init(cfg *settings.MySQLConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Dbname,
		)
	// 也可以使用MustConnect连接不成功就panic
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Error("connect DB failed", zap.Error(err))
		return
	}
	sqldb, err := db.DB()
	if err != nil {
		zap.L().Error("get sqldb failed", zap.Error(err))
		return
	}
	sqldb.SetMaxOpenConns(cfg.MaxOpenConns)
	sqldb.SetMaxIdleConns(cfg.MaxIdleConns)
	return
}


func Close() {
	sqldb, _ := db.DB()
	if err := sqldb.Close(); err != nil {
		return
	}
}

