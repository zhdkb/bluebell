package redis

import "bluebell/models"

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
