package redis

import (
	"bluebell/models"

)

func GetPostIDsInorder(p *models.ParamPostList) ([]string, error) {
	// 从redis获取id
	// 根据用户请求中携带的order参数查询id
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}

	// 确定查询的索引起始点
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	// ZRevRange 按分数从大到小查询指定数量的元素查询
	return rdb.ZRevRange(key, start, end).Result()
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
