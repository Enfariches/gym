package models

type Employee struct{
	ID string
	Name string
	SecondName string
	Surname string
	Age uint
	Sex bool
	Number string
	Departament string
	Post string // Должность
	Email string
	PassHash []byte
}