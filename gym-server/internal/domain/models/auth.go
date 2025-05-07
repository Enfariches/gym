package models

type AuthUser struct {
	ID int
	Email string
	PassHash []byte
	Source string
}