package models

import "time"

type Media struct {
	ID            int64     `db:"id"`
	Title         string    `db:"title"`
	PressignedUrl string    `db:"-"`
	AdminID       int64     `db:"admin_id"`
	DepartmentID  int64     `db:"department_id"`
	CreatedAt     time.Time `db:"created_at"`
}
