package models

type Admin struct{
	ID string
	Email string
	PassHash []byte
}