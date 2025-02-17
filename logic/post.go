package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"context"

	"go.uber.org/zap"
)

func CreatePost(ctx context.Context, p *models.Post) (error) {
	// 生成post id
	p.ID = int64(snowflake.GenID())

	// 保存到数据库
	err := mysql.CreatePost(ctx, p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(ctx, p.ID, p.CommunityID)
	return err

}

// GetPostById 根据帖子id查询帖子详情数据
func GetPostById(ctx context.Context, pid int64) (data *models.ApiPostDetail, err error) {
	// 查询并组合我们接口想用的数据
	post, err := mysql.GetPostById(ctx, pid)
	if err != nil {
		zap.L().Error("mysql.GetPostById(pid) failed", zap.Error(err))
		return
	}
	// 根据作者id查询作者信息
	user, err := mysql.GetUserById(ctx, post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserById(pid) failed", zap.Error(err))
		return
	}
	// 根据社区id查询社区详细信息
	community, err := mysql.GetCommunityDetailByID(ctx, post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID(pid) failed", zap.Error(err))
		return
	}
	data = &models.ApiPostDetail{
		AuthorName: user.Username,
		Post: post,
		CommunityDetail: community,
	}
	return
}

// GetPostList 获取帖子列表
func GetPostList(ctx context.Context, page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(ctx, page, size)
	if err != nil {
		return
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range(posts) {
		user, err := mysql.GetUserById(ctx, post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(pid) failed", zap.Error(err))
			continue
		}
		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(ctx, post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(pid) failed", zap.Error(err))
			continue
		}
		postdetail := &models.ApiPostDetail{
			AuthorName: user.Username,
			Post: post,
			CommunityDetail: community,
		}
		data = append(data, postdetail)
	}
	return
}

func GetPostList2(ctx context.Context, p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 去redis查询ID列表
	ids, err := redis.GetPostIDsInorder(ctx, p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInorder(p) return 0 data")
		return
	}
    // 根据id去MySQL数据库查询帖子详细信息
	posts, err := mysql.GetPostListByIDs(ctx, ids)
	if err != nil {
		return
	}

	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ctx, ids)
	if err != nil {
		return
	}

	data = make([]*models.ApiPostDetail, 0, len(posts))
	for i, post := range(posts) {
		user, err := mysql.GetUserById(ctx, post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(pid) failed", zap.Error(err))
			continue
		}
		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(ctx, post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(pid) failed", zap.Error(err))
			continue
		}
		postdetail := &models.ApiPostDetail{
			AuthorName: user.Username,
			VoteNum: voteData[i],
			Post: post,
			CommunityDetail: community,
		}
		data = append(data, postdetail)
	}
	return
}


func GetCommunityPostList(ctx context.Context, p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 去redis查询ID列表
	ids, err := redis.GetCommunityPostIDsInorder(ctx, p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInorder(p) return 0 data")
		return
	}
    // 根据id去MySQL数据库查询帖子详细信息
	posts, err := mysql.GetPostListByIDs(ctx, ids)
	if err != nil {
		return
	}

	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ctx, ids)
	if err != nil {
		return
	}

	data = make([]*models.ApiPostDetail, 0, len(posts))
	for i, post := range(posts) {
		user, err := mysql.GetUserById(ctx, post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(pid) failed", zap.Error(err))
			continue
		}
		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(ctx, post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(pid) failed", zap.Error(err))
			continue
		}
		postdetail := &models.ApiPostDetail{
			AuthorName: user.Username,
			VoteNum: voteData[i],
			Post: post,
			CommunityDetail: community,
		}
		data = append(data, postdetail)
	}
	return
}

// GetPostListNew 将两个查询逻辑合二为一
func GetPostListNew(ctx context.Context, p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	if p.CommunityID == 0 {
		data, err = GetPostList2(ctx, p)
	} else {
		data, err = GetCommunityPostList(ctx, p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return
	}
	return
}
