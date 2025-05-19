package models

import "time"

type Statistics struct {
	ID              int64     `db:"id"`
	Progress        string    `db:"progress"`
	PercentageView  int64     `db:"percentage_view"`
	EmployeeName    string    `db:"-"`
	EmployeeSurname string    `db:"-"`
	MediaTitle      string    `db:"department_id"`
	CreatedAt       time.Time `db:"created_at"`
}
