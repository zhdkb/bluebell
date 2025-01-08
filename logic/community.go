package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"context"
)

func GetCommunityList(ctx context.Context) ([]*models.Community, error) {
	// 查找到所有的community并返回
	return mysql.GetCommunityList(ctx)

}

func GetCommunityDetail(ctx context.Context, id int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetailByID(ctx, id)
}
