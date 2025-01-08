package logic

import (
	"bluebell/dao/redis"
	"bluebell/models"
	"context"
	"strconv"
)

// 投票功能
// 用户投票的数据

// VoteForPost 为帖子投票
func VoteForPost(ctx context.Context, userID int64, p *models.ParamVoteData) error {
	return redis.VoteForPost(ctx, strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
}
