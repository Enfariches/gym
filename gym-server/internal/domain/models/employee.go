package models

type Employee struct {
	ID         int64  `db:"id"`
	Name       string `db:"name"`
	SecondName string `db:"second_name"`
	Surname    string `db:"surname"`
	Age        int64  `db:"age"`
	Sex        bool   `db:"sex"`
	Phone      string `db:"number"`
	Department string `db:"department"`
	Post       string `db:"post"` // Должность
	Email      string `db:"email"`
}
