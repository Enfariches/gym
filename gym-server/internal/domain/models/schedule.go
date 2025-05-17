package models

import "time"

type Schedule struct {
	ID             int64     `db:"id"`
	CronExpression string    `db:"cron_expression"`
	IsActive       bool      `db:"is_active"`
	MediaID        int64     `db:"media_id"`
	AdminID        int64     `db:"admin_id"`
	CreatedAt      time.Time `db:"created_at"`
}
