package models

import "gorm.io/gorm"

type User struct {
	ID           uint           `gorm:"primarykey"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Organization string
	Department   string
	Position     string
	UserId       uint
}

//Нет методов создания и обновления таблицы Users
