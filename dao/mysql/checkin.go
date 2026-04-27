package mysql

import (
	"bluebell/models"
	"context"
	"errors"
	"time"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrorAlreadyCheckIn = errors.New("今日已签到")

// CreateCheckIn 在一个事务里写入签到明细，并更新用户签到统计。
func CreateCheckIn(ctx context.Context, detail *models.CheckInDetail) (*models.CheckInResult, error) {
	signDate := dateOnly(detail.SignTime)
	detail.SignDate = signDate

	var result *models.CheckInResult
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		counter := new(models.CheckInCount)
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", detail.UserID).
			First(counter).Error

		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			counter = &models.CheckInCount{
				UserID:          detail.UserID,
				TotalCount:      1,
				ContinuousCount: 1,
				LastSignDate:    signDate,
			}
			if err := tx.Create(detail).Error; err != nil {
				return normalizeCheckInError(err)
			}
			if err := tx.Create(counter).Error; err != nil {
				return err
			}
		case err != nil:
			return err
		default:
			if sameDate(counter.LastSignDate, signDate) {
				return ErrorAlreadyCheckIn
			}

			if err := tx.Create(detail).Error; err != nil {
				return normalizeCheckInError(err)
			}

			counter.TotalCount++
			if sameDate(counter.LastSignDate, signDate.AddDate(0, 0, -1)) {
				counter.ContinuousCount++
			} else {
				counter.ContinuousCount = 1
			}
			counter.LastSignDate = signDate

			if err := tx.Model(counter).
				Select("total_count", "continuous_count", "last_sign_date").
				Updates(counter).Error; err != nil {
				return err
			}
		}

		result = &models.CheckInResult{
			CheckInID:       detail.CheckInID,
			UserID:          detail.UserID,
			SignDate:        signDate.Format("2006-01-02"),
			SignTime:        detail.SignTime.Format(time.RFC3339),
			TotalCount:      counter.TotalCount,
			ContinuousCount: counter.ContinuousCount,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func normalizeCheckInError(err error) error {
	var mysqlErr *mysqlDriver.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return ErrorAlreadyCheckIn
	}
	return err
}

func dateOnly(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func sameDate(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}
