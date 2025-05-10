package models

type AuthUser struct {
	ID       int64  `db:"id"`
	Email    string `db:"email"`
	PassHash []byte `db:"passhash"`
	Source   string `db:"-"`
}
