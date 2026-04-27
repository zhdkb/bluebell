package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/kafka"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"context"
	"encoding/json"
	"strconv"
	"time"

	"go.uber.org/zap"
)

// CheckIn 完成用户签到：MySQL 是最终事实源，Redis 只在成功后做缓存，Kafka 发送签到事件。
func CheckIn(ctx context.Context, userID int64) (*models.CheckInResult, error) {
	now := time.Now()

	checked, err := redis.HasCheckInToday(ctx, userID, now)
	if err != nil {
		zap.L().Warn("redis.HasCheckInToday failed", zap.Int64("userID", userID), zap.Error(err))
	} else if checked {
		return nil, mysql.ErrorAlreadyCheckIn
	}

	detail := &models.CheckInDetail{
		CheckInID: snowflake.GenID(),
		UserID:    userID,
		SignTime:  now,
	}

	result, err := mysql.CreateCheckIn(ctx, detail)
	if err != nil {
		return nil, err
	}

	if ok, _, err := redis.MarkCheckInToday(ctx, userID, now); err != nil {
		zap.L().Warn("redis.MarkCheckInToday failed", zap.Int64("userID", userID), zap.Error(err))
	} else if !ok {
		zap.L().Warn("redis checkin mark already exists after mysql success", zap.Int64("userID", userID))
	}

	if err := redis.CacheCheckInResult(ctx, result); err != nil {
		zap.L().Warn("redis.CacheCheckInResult failed", zap.Int64("userID", userID), zap.Error(err))
	}

	publishCheckInEvent(ctx, result)
	return result, nil
}

func GetMonthlyCheckIn(ctx context.Context, userID int64, month string) (*models.CheckInMonthResult, error) {
	targetMonth, err := parseCheckInMonth(month)
	if err != nil {
		return nil, err
	}

	// 月度打卡记录优先读 Redis bitmap；没有缓存时再查 MySQL。
	days, ok, err := redis.GetMonthlyCheckInDays(ctx, userID, targetMonth)
	if err != nil {
		zap.L().Warn("redis.GetMonthlyCheckInDays failed", zap.Int64("userID", userID), zap.String("month", targetMonth.Format("2006-01")), zap.Error(err))
	}
	if ok && err == nil {
		return &models.CheckInMonthResult{
			UserID:      userID,
			Month:       targetMonth.Format("2006-01"),
			CheckInDays: days,
			TotalDays:   len(days),
		}, nil
	}

	days, err = mysql.GetMonthlyCheckInDays(ctx, userID, targetMonth)
	if err != nil {
		return nil, err
	}
	// MySQL 是最终事实源，回源成功后把结果写回 bitmap，加速下次查询。
	if err := redis.CacheMonthlyCheckInDays(ctx, userID, targetMonth, days); err != nil {
		zap.L().Warn("redis.CacheMonthlyCheckInDays failed", zap.Int64("userID", userID), zap.String("month", targetMonth.Format("2006-01")), zap.Error(err))
	}

	return &models.CheckInMonthResult{
		UserID:      userID,
		Month:       targetMonth.Format("2006-01"),
		CheckInDays: days,
		TotalDays:   len(days),
	}, nil
}

func parseCheckInMonth(month string) (time.Time, error) {
	if month == "" {
		now := time.Now()
		return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()), nil
	}
	return time.ParseInLocation("2006-01", month, time.Local)
}

func publishCheckInEvent(ctx context.Context, data *models.CheckInResult) {
	manager := kafka.GetManager()
	if manager == nil {
		return
	}

	payload, err := json.Marshal(data)
	if err != nil {
		zap.L().Warn("marshal checkin event failed", zap.Error(err))
		return
	}

	key := []byte(strconv.FormatInt(data.UserID, 10))
	if err := manager.Publish(ctx, kafka.TopicCheckIn, key, payload); err != nil {
		zap.L().Warn("publish checkin event failed", zap.Int64("userID", data.UserID), zap.Error(err))
	}
}
