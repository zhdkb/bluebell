package logic

import (
	"bluebell/dao/redis"
	"bluebell/models"
	"strconv"
)

// 投票功能
// 用户投票的数据

// VoteForPost 为帖子投票
func VoteForPost(userID int64, p *models.ParamVoteData) error {
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
}
