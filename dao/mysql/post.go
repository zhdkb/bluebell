package mysql

import (
	"bluebell/models"
	"context"
	"strings"
)

func CreatePost(ctx context.Context, p *models.Post) (err error) {
	sqlStr := `insert into post
				(post_id, title, content, author_id, community_id)
				values (?, ?, ?, ?, ?)`
	err = db.WithContext(ctx).Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID).Error
	return
}

// GetPostById 根据id查询单个帖子数据
func GetPostById(ctx context.Context, pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select post_id, title, content, author_id, community_id, create_time 
				from post 
				where post_id = ?`
	// result := db.First(post, pid)
	err = db.WithContext(ctx).Raw(sqlStr, pid).Scan(&post).Error
	return
}

// GetPostList 查询帖子列表函数
func GetPostList(ctx context.Context, page, size int64) (posts []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time 
				from post 
				order by create_time
				DESC
				limit ?, ?`
	
	posts = make([]*models.Post, 0, 2)
	err = db.WithContext(ctx).Raw(sqlStr, (page - 1) * size, size).Scan(&posts).Error
	return
}

// GetPostListByIDs 根据给定的id列表查询帖子数据
func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	idString := strings.Join(ids, ",")
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
				from post
				where post_id in (?)
				order by FIND_IN_SET(post_id, ?)`
	err = db.Raw(sqlStr, idString, idString).Scan(&postList).Error
	// query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	// if err != nil {
	// 	return nil, err
	// }
	// query = db.Rebind(query)
	// err = db.Select(&postList, query, args...)
	return
}
