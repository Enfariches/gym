package models

type AuthUser struct {
	ID           int64  `db:"id"`
	Email        string `db:"email"`
	DepartmentID int64  `db:"department_id"`
	PassHash     []byte `db:"passhash"`
	Source       string `db:"-"`
}
