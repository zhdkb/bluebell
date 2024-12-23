package redis

import (
	"bluebell/models"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func getIDsFromKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1
	// ZRevRange 按分数从大到小查询指定数量的元素查询
	return rdb.ZRevRange(key, start, end).Result()
}

func GetPostIDsInorder(p *models.ParamPostList) ([]string, error) {
	// 从redis获取id
	// 根据用户请求中携带的order参数查询id
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}

	// 确定查询的索引起始点
	return getIDsFromKey(key, p.Page, p.Size)
}

// GetPostVoteData 根据ids查询每篇帖子的投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	data = make([]int64, 0)
	for _, id := range(ids) {
		key := getRedisKey(KeyPostVotedZsetPre + id)
		// 查找key中分数是1的元素的数量 -> 统计每篇帖子的赞成票的数量
		v1 := rdb.ZCount(key, "1", "1").Val()
		data = append(data, v1)
	}

	// // 使用pipeline一次发送多条命令，减少RTT
	// pipeline := rdb.Pipeline()
	// for _, id := range ids {
	// 	key := getRedisKey(KeyPostVotedZsetPre + id)
	// 	pipeline.ZCount(key, "1", "1")
	// }
	// cmders, err := pipeline.Exec()
	// if err != nil {
	// 	return nil, err
	// }
	// data = make([]int64, 0)
	// for _, cmder := range cmders {
	// 	v := cmder.(*redis.IntCmd).Val()
	// 	data = append(data, v)
	// }
	return
}

// GetCommunityPostIDsInorder 按社区查询ids
func GetCommunityPostIDsInorder(p *models.ParamPostList) ([]string, error) {
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
	if rdb.Exists(key).Val() < 1 {
		// 不存在，需要计算
		pipeline := rdb.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, cKey, orderKey) // ZInterStore 计算
		pipeline.Expire(key, 60 * time.Second)
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	// 存在的话就直接根据key查询ids
	return getIDsFromKey(key, p.Page, p.Size)
}
