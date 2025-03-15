package models

import (
	"github.com/google/uuid"
	"log"
)

type EmailConfirm struct {
	ID    string `gorm:"primary_key"`
	Email string `sql:"index,unique" ;gorm:"index,unique"`
	Code  uint
}

func GetEmailConfirm(id string) *EmailConfirm {
	emailConfirm := &EmailConfirm{}
	err := Db.Table("email_confirms").Where("id = ?", id).First(&emailConfirm).Error
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	if emailConfirm.Email == "" {
		return nil
	}
	return emailConfirm
}

func CreateEmailConfirm(email string) *EmailConfirm {
	exEmail := &EmailConfirm{}
	Db.Table("email_confirms").Where("email = ?", email).First(&exEmail)
	if exEmail != nil && exEmail.Email != "" {
		UpdateEmailConfirm(exEmail, 0)
		return exEmail
	}
	emailConfirm := &EmailConfirm{
		ID:    uuid.New().String(),
		Email: email,
	}
	Db.Create(&emailConfirm)
	return emailConfirm
}

func UpdateEmailConfirm(emailConfirm *EmailConfirm, code uint) *EmailConfirm {
	Db.Model(&emailConfirm).Update("code", code)
	return emailConfirm
}

func DeleteEmailConfirm(email string) {
	Db.Table("email_confirms").Where("email = ?", email).Unscoped().Delete(&Session{})
}
