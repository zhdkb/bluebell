package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (error) {
	// 生成post id
	p.ID = int64(snowflake.GenID())

	// 保存到数据库
	err := mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.ID)
	return err

}

// GetPostById 根据帖子id查询帖子详情数据
func GetPostById(pid int64) (data *models.ApiPostDetail, err error) {
	// 查询并组合我们接口想用的数据
	post, err := mysql.GetPostById(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostById(pid) failed", zap.Error(err))
		return
	}
	// 根据作者id查询作者信息
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserById(pid) failed", zap.Error(err))
		return
	}
	// 根据社区id查询社区详细信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
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
func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range(posts) {
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(pid) failed", zap.Error(err))
			continue
		}
		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
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

func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 去redis查询ID列表
	ids, err := redis.GetPostIDsInorder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInorder(p) return 0 data")
		return
	}
    // 根据id去MySQL数据库查询帖子详细信息
	posts, err := mysql.GetPostListByIDs((ids))
	if err != nil {
		return
	}

	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range(posts) {
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(pid) failed", zap.Error(err))
			continue
		}
		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
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
