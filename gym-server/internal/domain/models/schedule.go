package models

import "time"

type Schedule struct {
	ID             int64     `db:"id"`
	CronExpression string    `db:"cron_expression"`
	IsActive       bool      `db:"is_active"`
	VideoID        int64     `db:"video_id"`
	AdminID        int64     `db:"admin_id"`
	CreatedAt      time.Time `db:"created_at"`
}
