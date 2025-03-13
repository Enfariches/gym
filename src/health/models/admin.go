package models

import "log"

type Admin struct {
	Email string `gorm:"primary_key"`
}

func CreateAdmin(email string) *Admin {
	admin := &Admin{Email: email}
	Db.Create(admin)
	return admin
}

func HasAdmin(email string) bool {
	var count int64
	err := Db.Table("admins").
		Select("*").
		Where("email = ?", email).
		Count(&count).
		Error
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return count > 0
}
