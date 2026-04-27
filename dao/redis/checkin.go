package redis

import (
	"bluebell/models"
	"context"
	"strconv"
	"time"

	goredis "github.com/go-redis/redis/v8"
)

const checkInBitmapTTL = 400 * 24 * time.Hour

// MarkCheckInToday 使用 SETNX 兼容每日标记，并用 bitmap 标记用户今天已经签到。
func MarkCheckInToday(ctx context.Context, userID int64, now time.Time) (bool, string, error) {
	key := getCheckInDailyKey(userID, now)
	expire := time.Until(nextDay(now).Add(time.Hour))

	bitmapKey := getCheckInBitmapKey(userID, now)
	// bitmap 的 offset 从 0 开始，所以 1 号对应 0，31 号对应 30。
	offset := int64(now.Day() - 1)

	pipeline := rdb.Pipeline()
	setNXCmd := pipeline.SetNX(ctx, key, "1", expire)
	pipeline.SetBit(ctx, bitmapKey, offset, 1)
	pipeline.Expire(ctx, bitmapKey, checkInBitmapTTL)
	_, err := pipeline.Exec(ctx)
	return setNXCmd.Val(), key, err
}

func HasCheckInToday(ctx context.Context, userID int64, now time.Time) (bool, error) {
	bitmapKey := getCheckInBitmapKey(userID, now)
	// 优先读 bitmap，命中后可以直接拦截重复签到。
	bit, err := rdb.GetBit(ctx, bitmapKey, int64(now.Day()-1)).Result()
	if err != nil {
		return false, err
	}
	if bit == 1 {
		return true, nil
	}

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

func GetMonthlyCheckInDays(ctx context.Context, userID int64, month time.Time) ([]string, bool, error) {
	key := getCheckInBitmapKey(userID, month)
	// 不存在 bitmap key 时交给上层回源 MySQL，再把结果回填到 Redis。
	exists, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		return nil, false, err
	}
	if exists == 0 {
		return nil, false, nil
	}

	monthStart := monthStart(month)
	daysInMonth := monthStart.AddDate(0, 1, -1).Day()

	pipeline := rdb.Pipeline()
	cmds := make([]*goredis.IntCmd, 0, daysInMonth)
	// 一次 pipeline 读完当月所有 bit，减少 Redis 往返次数。
	for day := 1; day <= daysInMonth; day++ {
		cmds = append(cmds, pipeline.GetBit(ctx, key, int64(day-1)))
	}
	if _, err := pipeline.Exec(ctx); err != nil {
		return nil, true, err
	}

	days := make([]string, 0)
	for i, cmd := range cmds {
		if cmd.Val() == 1 {
			days = append(days, monthStart.AddDate(0, 0, i).Format("2006-01-02"))
		}
	}
	return days, true, nil
}

func CacheMonthlyCheckInDays(ctx context.Context, userID int64, month time.Time, days []string) error {
	key := getCheckInBitmapKey(userID, month)
	pipeline := rdb.Pipeline()
	// MySQL 回源得到的日期列表会被回填成 bitmap，后续月度查询直接读 Redis。
	for _, day := range days {
		t, err := time.ParseInLocation("2006-01-02", day, time.Local)
		if err != nil {
			return err
		}
		pipeline.SetBit(ctx, key, int64(t.Day()-1), 1)
	}
	pipeline.Expire(ctx, key, checkInBitmapTTL)
	_, err := pipeline.Exec(ctx)
	return err
}

func nextDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d+1, 0, 0, 0, 0, t.Location())
}

func getCheckInDailyKey(userID int64, now time.Time) string {
	date := now.Format("2006-01-02")
	return getRedisKey(KeyCheckInDailyPF + strconv.FormatInt(userID, 10) + ":" + date)
}

func getCheckInBitmapKey(userID int64, month time.Time) string {
	return getRedisKey(KeyCheckInBitmapPF + strconv.FormatInt(userID, 10) + ":" + month.Format("2006-01"))
}

func monthStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}
