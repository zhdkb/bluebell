package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/kafka"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"context"
	"encoding/json"
	"strconv"

	"go.uber.org/zap"
)

func LikePost(ctx context.Context, userID int64, p *models.ParamLikeData) (*models.LikeResult, error) {
	liked := int8(0)
	if p.Action == models.LikeActionLike {
		liked = 1
	}

	result, err := redis.LikePost(ctx, strconv.FormatInt(userID, 10), p.PostID, liked)
	if err != nil {
		return nil, err
	}
	if !result.Changed {
		return result, nil
	}

	postID, err := strconv.ParseInt(p.PostID, 10, 64)
	if err != nil {
		return nil, err
	}
	event := &models.PostLikeEvent{
		EventID: snowflake.GenID(),
		PostID:  postID,
		UserID:  userID,
		Liked:   liked,
		Delta:   result.Delta,
	}

	if err := publishLikeEvent(ctx, event); err != nil {
		zap.L().Warn("publish like event failed, save failed event",
			zap.Int64("postID", event.PostID),
			zap.Int64("userID", event.UserID),
			zap.Int64("eventID", event.EventID),
			zap.Error(err),
		)
		if saveErr := mysql.SaveFailedLikeEvent(ctx, event); saveErr != nil {
			zap.L().Error("mysql.SaveFailedLikeEvent failed",
				zap.Int64("postID", event.PostID),
				zap.Int64("userID", event.UserID),
				zap.Int64("eventID", event.EventID),
				zap.Error(saveErr),
			)
			return nil, saveErr
		}
	}

	return result, nil
}

func publishLikeEvent(ctx context.Context, event *models.PostLikeEvent) error {
	manager := kafka.GetManager()
	if manager == nil {
		return kafka.ErrManagerNotReady
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	key := []byte(strconv.FormatInt(event.PostID, 10) + ":" + strconv.FormatInt(event.UserID, 10))
	return manager.Publish(ctx, kafka.TopicLike, key, payload)
}
