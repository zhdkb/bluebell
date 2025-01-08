package redis

import (
	"context"
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// 本项目使用简化版的投票分数
// 用户投一票就加432分   86400/200 -> 需要200张赞成票可以给你的帖子续一天

/* 投票的几种情况
direction=1时，有两种情况：
	1.之前没有投过票，现在投赞成票 +432
	2.之前投反对票，现在改投赞成票 +432*2
direction=0时，有两种情况
	1.之前投过赞成票，现在要取消投票 -432
	2.之前投过反对票，现在要取消投票 +432
direction=-1时，有两种情况：
	1.之前没有投过票，现在投反对票 -432
	2.之前投赞成票，现在改投反对票 -432*2

投票的限制：
每个帖子自发表之日起，一个星期之内运行用户投票，超过一个星期就不运行在投票了
	1.到期之后将redis中保存的赞成票数和反对票数存储到MySQL表中
	2.到期之后删除那个 KeyPostVotedZsetPre
*/

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePreVote = 432 // 每一票值多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepested = errors.New("不允许重复投票")
)

func CreatePost(ctx context.Context, postID, communityID int64) error {

	pipeline := rdb.TxPipeline()

	pipeline.ZAdd(ctx, getRedisKey(KeyPostTimeZSet), &redis.Z{
		Score: float64(time.Now().Unix()),
		Member: postID,
	})

	pipeline.ZAdd(ctx, getRedisKey(KeyPostScoreZSet), &redis.Z{
		Score: float64(time.Now().Unix()),
		Member: postID,
	})

	// 把帖子id加到社区的set
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	pipeline.SAdd(ctx, cKey, postID)
	_, err := pipeline.Exec(ctx)

	return err
}

func VoteForPost(ctx context.Context, userID, postID string, value float64) error {
	// 1.判断投票的限制
	postTime := rdb.ZScore(ctx, getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix()) - postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}

	// 2和3需要放到一个pipeline事务中操作

	// 2.更新帖子的分数
	// 先查当前用户给当前帖子的投票记录
	ov := rdb.ZScore(ctx, getRedisKey(KeyPostVotedZsetPre + postID), userID).Val()
	if value == ov {
		return ErrVoteRepested
	}
	var dir float64
	if value > ov {
		dir = 1
	} else {
		dir = -1
	}
	diff := math.Abs(ov - value) // 计算两次投票的差值

	pipeline := rdb.TxPipeline()
	pipeline.ZIncrBy(ctx, getRedisKey(KeyPostScoreZSet), dir * diff * float64(scorePreVote), postID)

	// 3.记录用户为该帖子投过票
	if value == 0 {
		pipeline.ZRem(ctx, getRedisKey(KeyPostVotedZsetPre + postID), userID)
	} else {
		pipeline.ZAdd(ctx, getRedisKey(KeyPostVotedZsetPre + postID), &redis.Z{
			Score: value, // 赞成票还是反对票
			Member: userID,
		})
	}

	_, err := pipeline.Exec(ctx)
	return err

}

