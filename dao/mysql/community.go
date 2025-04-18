package mysql

import (
	"bluebell/models"
	"context"
	"database/sql"
	"fmt"

	"go.uber.org/zap"
)

func GetCommunityList(ctx context.Context) (data []*models.Community, err error) {
	sqlStr := `select community_id, community_name from community`
	// db.Select(&data, sqlStr)
	if err = db.WithContext(ctx).Raw(sqlStr).Scan(&data).Error; err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in db")
			err = nil
		}
	}
	for _, i := range data {
		fmt.Println(i)
	}
	return
}

// GetCommunityDetailByID 根据ID查询社区详情
func GetCommunityDetailByID(ctx context.Context, id int64) (community *models.CommunityDetail, err error) {
	community = new(models.CommunityDetail)
	sqlStr := `select community_id, community_name, introduction, create_time
				from community
				where community_id = ?`
	// db.Get(community, sqlStr, id)
	if err = db.WithContext(ctx).Raw(sqlStr, id).Scan(&community).Error; err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidID
		}
	}
	fmt.Println(community)
	return
}
