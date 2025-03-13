package models

import "strings"

type Username struct {
	ID       uint   `gorm:"primary_key" ;json:"id"`
	Username string `gorm:"index,unique" ;json:"username"`
}

func GetOrCreateUsernameId(username string) uint {
	name := &Username{}
	Db.Table("usernames").
		Where("lower(username) = lower(?)", strings.ToLower(username)).
		Limit(1).
		Find(&name)
	if name.Username == "" {
		name.Username = strings.ToLower(username)
		Db.Create(&name)
	}
	return name.ID
}
