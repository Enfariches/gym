package models

type Admin struct {
	ID          int64
	Name        *string
	Surname     *string
	Email       string
	Departament *string
}
