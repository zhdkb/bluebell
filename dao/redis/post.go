package redis

import (
	"bluebell/models"
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

func getIDsFromKey(ctx context.Context, key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1
	// ZRevRange 按分数从大到小查询指定数量的元素查询
	return rdb.ZRevRange(ctx, key, start, end).Result()
}

func GetPostIDsInorder(ctx context.Context, p *models.ParamPostList) ([]string, error) {
	// 从redis获取id
	// 根据用户请求中携带的order参数查询id
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}

	// 确定查询的索引起始点
	return getIDsFromKey(ctx, key, p.Page, p.Size)
}

// GetCommunityPostIDsInorder 按社区查询ids
func GetCommunityPostIDsInorder(ctx context.Context, p *models.ParamPostList) ([]string, error) {
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}

	// 使用zinterstore 把分区的帖子set与帖子分数的zset 生成一个新的zset
	// 针对新的zset 按之前的逻辑取数据

	// 社区的key
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))

	// 利用缓存key减少zinterstore执行的次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if rdb.Exists(ctx, key).Val() < 1 {
		// 不存在，需要计算
		pipeline := rdb.Pipeline()
		pipeline.ZInterStore(ctx, key, &redis.ZStore{
			Aggregate: "MAX",
			Keys:      []string{cKey, orderKey},
		}) // ZInterStore 计算
		pipeline.Expire(ctx, key, 60*time.Second)
		_, err := pipeline.Exec(ctx)
		if err != nil {
			return nil, err
		}
	}
	// 存在的话就直接根据key查询ids
	return getIDsFromKey(ctx, key, p.Page, p.Size)
}
