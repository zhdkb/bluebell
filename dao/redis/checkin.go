package redis

import (
	"bluebell/models"
	"context"
	"strconv"
	"time"
)

// MarkCheckInToday 使用 SETNX 标记用户今天已经签到。
func MarkCheckInToday(ctx context.Context, userID int64, now time.Time) (bool, string, error) {
	key := getCheckInDailyKey(userID, now)
	expire := time.Until(nextDay(now).Add(time.Hour))
	ok, err := rdb.SetNX(ctx, key, "1", expire).Result()
	return ok, key, err
}

func HasCheckInToday(ctx context.Context, userID int64, now time.Time) (bool, error) {
	key := getCheckInDailyKey(userID, now)
	n, err := rdb.Exists(ctx, key).Result()
	return n > 0, err
}

func RollbackCheckInMark(ctx context.Context, key string) {
	if key == "" {
		return
	}
	_ = rdb.Del(ctx, key).Err()
}

func CacheCheckInResult(ctx context.Context, data *models.CheckInResult) error {
	key := getRedisKey(KeyCheckInCountPF + strconv.FormatInt(data.UserID, 10))
	return rdb.HSet(ctx, key, map[string]interface{}{
		"total_count":      data.TotalCount,
		"continuous_count": data.ContinuousCount,
		"last_sign_date":   data.SignDate,
	}).Err()
}

func nextDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d+1, 0, 0, 0, 0, t.Location())
}

func getCheckInDailyKey(userID int64, now time.Time) string {
	date := now.Format("2006-01-02")
	return getRedisKey(KeyCheckInDailyPF + strconv.FormatInt(userID, 10) + ":" + date)
}
