package redis

// redis key

// redis key尽量使用命名空间的方式，方便查询和拆分

const (
	KeyPrefix        = "bluebell:"
	KeyPostTimeZSet  = "post:time"  // zset;帖子及发帖时间
	KeyPostScoreZSet = "post:score" // zset;帖子及点赞数

	KeyCommunitySetPF  = "community:"       // set;保存每个分区下帖子的id
	KeyPostLikedSetPF  = "post:liked:"      // set;参数是post_id，保存点赞用户id
	KeyPostLikeCountPF = "post:like_count:" // string;参数是post_id，保存点赞数
	KeyCheckInDailyPF  = "checkin:daily:"   // string;参数是user_id:yyyy-mm-dd，兼容每日签到标记
	KeyCheckInBitmapPF = "checkin:bitmap:"  // bitmap;参数是user_id:yyyy-mm，offset为day-1
	KeyCheckInCountPF  = "checkin:count:"   // hash;参数是user_id，缓存签到统计

	KeyRefreshTokenBlacklistPF = "token:refresh:blacklist:" // string;参数是refresh token的jti
)

func getRedisKey(key string) string {
	return KeyPrefix + key
}
