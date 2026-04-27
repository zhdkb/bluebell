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
