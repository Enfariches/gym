package models

import (
	"gorm.io/gorm"
)

type Shedule struct {
	gorm.Model
	DayOfWeek uint8
	Time      string
	VideoID   uint
}

func AddSchedule(day uint8) *Shedule {
	shed := Shedule{
		DayOfWeek: day,
	}
	Db.Create(&shed)
	return &shed
}

func ListSchedule(day uint8) []Shedule {
	var shedules []Shedule
	Db.Where("day_of_week = ?", day).Find(&shedules)
	return shedules
}

func ListAllSchedule() []Shedule {
	var shedules []Shedule
	Db.Find(&shedules)
	return shedules
}

func (s *Shedule) Update() {
	Db.Save(&s)
}

func DeleteSchedule(id uint) {
	Db.Where("id = ?", id).Delete(&Shedule{})
}
