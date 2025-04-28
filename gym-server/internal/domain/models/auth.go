package models

type AuthUser struct {
	Id int
	Email string
	PassHash []byte
	Source string
}