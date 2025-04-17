package models

type Admin struct{
	ID string
	Name string
	Surname string
	Email string
	Departament string
	PassHash []byte
}