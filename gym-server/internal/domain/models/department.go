package models

type Department struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}
