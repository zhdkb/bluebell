package mysql

import (
	"bluebell/models"
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ApplyPostLikeEvent applies a like event idempotently by comparing the stored
// user-post state with the event state before changing post.like_count.
func ApplyPostLikeEvent(ctx context.Context, event *models.PostLikeEvent) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		like := new(models.PostLike)
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("post_id = ? and user_id = ?", event.PostID, event.UserID).
			First(like).Error

		oldLiked := int8(0)
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			like = &models.PostLike{
				PostID: event.PostID,
				UserID: event.UserID,
				Liked:  event.Liked,
			}
			if err := tx.Create(like).Error; err != nil {
				return err
			}
		case err != nil:
			return err
		default:
			oldLiked = like.Liked
			if oldLiked != event.Liked {
				if err := tx.Model(like).Update("liked", event.Liked).Error; err != nil {
					return err
				}
			}
		}

		delta := int64(event.Liked - oldLiked)
		if delta == 0 {
			return nil
		}

		if delta > 0 {
			return tx.Exec(
				`update post set like_count = like_count + ? where post_id = ?`,
				delta,
				event.PostID,
			).Error
		}

		return tx.Exec(
			`update post set like_count = greatest(like_count + ?, 0) where post_id = ?`,
			delta,
			event.PostID,
		).Error
	})
}

func SaveFailedLikeEvent(ctx context.Context, event *models.PostLikeEvent) error {
	sqlStr := `insert into like_event_failed(event_id, post_id, user_id, liked, delta)
				values (?, ?, ?, ?, ?)
				on duplicate key update
					retry_count = retry_count + 1,
					update_time = current_timestamp`
	return db.WithContext(ctx).Exec(
		sqlStr,
		event.EventID,
		event.PostID,
		event.UserID,
		event.Liked,
		event.Delta,
	).Error
}

func ListFailedLikeEvents(ctx context.Context, limit int) ([]*models.PostLikeEvent, error) {
	var rows []*models.FailedLikeEvent
	if err := db.WithContext(ctx).
		Order("id asc").
		Limit(limit).
		Find(&rows).Error; err != nil {
		return nil, err
	}

	events := make([]*models.PostLikeEvent, 0, len(rows))
	for _, row := range rows {
		events = append(events, &models.PostLikeEvent{
			EventID: row.EventID,
			PostID:  row.PostID,
			UserID:  row.UserID,
			Liked:   row.Liked,
			Delta:   row.Delta,
		})
	}
	return events, nil
}

func DeleteFailedLikeEvent(ctx context.Context, eventID int64) error {
	return db.WithContext(ctx).
		Where("event_id = ?", eventID).
		Delete(&models.FailedLikeEvent{}).Error
}

func GetPostLikeCount(ctx context.Context, postID int64) (int64, error) {
	var count int64
	err := db.WithContext(ctx).
		Raw(`select like_count from post where post_id = ?`, postID).
		Scan(&count).Error
	return count, err
}
