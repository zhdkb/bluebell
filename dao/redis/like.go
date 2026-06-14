package redis

import (
	"bluebell/models"
	"context"
	"errors"
	"strconv"
	"time"

	goredis "github.com/go-redis/redis/v8"
)

var likePostScript = goredis.NewScript(`
local likedKey = KEYS[1]
local countKey = KEYS[2]
local scoreKey = KEYS[3]

local userID = ARGV[1]
local postID = ARGV[2]
local liked = tonumber(ARGV[3])

local currentCount = redis.call("GET", countKey)
if currentCount then
	currentCount = tonumber(currentCount)
else
	currentCount = redis.call("SCARD", likedKey)
end

local exists = redis.call("SISMEMBER", likedKey, userID)
if liked == 1 then
	if exists == 1 then
		return {0, 1, 0, currentCount}
	end
	redis.call("SADD", likedKey, userID)
	currentCount = currentCount + 1
	redis.call("SET", countKey, currentCount)
	redis.call("ZINCRBY", scoreKey, 1, postID)
	return {1, 1, 1, currentCount}
end

if exists == 0 then
	return {0, 0, 0, currentCount}
end
redis.call("SREM", likedKey, userID)
currentCount = currentCount - 1
if currentCount < 0 then
	currentCount = 0
end
redis.call("SET", countKey, currentCount)
redis.call("ZINCRBY", scoreKey, -1, postID)
return {1, 0, -1, currentCount}
`)

func CreatePost(ctx context.Context, postID, communityID int64) error {

	pipeline := rdb.TxPipeline()

	pipeline.ZAdd(ctx, getRedisKey(KeyPostTimeZSet), &goredis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	pipeline.ZAdd(ctx, getRedisKey(KeyPostScoreZSet), &goredis.Z{
		Score:  0,
		Member: postID,
	})

	// 把帖子id加到社区的set
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	pipeline.SAdd(ctx, cKey, postID)
	_, err := pipeline.Exec(ctx)

	return err
}

func LikePost(ctx context.Context, userID, postID string, liked int8) (*models.LikeResult, error) {
	keys := []string{
		getRedisKey(KeyPostLikedSetPF + postID),
		getRedisKey(KeyPostLikeCountPF + postID),
		getRedisKey(KeyPostScoreZSet),
	}
	values, err := likePostScript.Run(ctx, rdb, keys, userID, postID, liked).Slice()
	if err != nil {
		return nil, err
	}
	if len(values) != 4 {
		return nil, errors.New("redis like script return invalid data")
	}

	changed, err := redisInt64(values[0])
	if err != nil {
		return nil, err
	}
	currentLiked, err := redisInt64(values[1])
	if err != nil {
		return nil, err
	}
	delta, err := redisInt64(values[2])
	if err != nil {
		return nil, err
	}
	count, err := redisInt64(values[3])
	if err != nil {
		return nil, err
	}

	return &models.LikeResult{
		PostID:    postID,
		UserID:    userID,
		Liked:     currentLiked == 1,
		Changed:   changed == 1,
		Delta:     delta,
		LikeCount: count,
	}, nil
}

func GetPostLikeCount(ctx context.Context, postID string) (int64, bool, error) {
	key := getRedisKey(KeyPostLikeCountPF + postID)
	count, err := rdb.Get(ctx, key).Int64()
	if errors.Is(err, goredis.Nil) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return count, true, nil
}

func CachePostLikeCount(ctx context.Context, postID string, count int64) error {
	key := getRedisKey(KeyPostLikeCountPF + postID)
	return rdb.Set(ctx, key, count, 0).Err()
}

func GetPostLikeData(ctx context.Context, ids []string) ([]int64, error) {
	data := make([]int64, 0, len(ids))
	for _, id := range ids {
		count, _, err := GetPostLikeCount(ctx, id)
		if err != nil {
			return nil, err
		}
		data = append(data, count)
	}
	return data, nil
}

func redisInt64(v interface{}) (int64, error) {
	switch value := v.(type) {
	case int64:
		return value, nil
	case int:
		return int64(value), nil
	case string:
		return strconv.ParseInt(value, 10, 64)
	case []byte:
		return strconv.ParseInt(string(value), 10, 64)
	default:
		return 0, errors.New("unsupported redis integer type")
	}
}
