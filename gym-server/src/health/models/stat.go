package models

import (
	"gorm.io/gorm"
	"time"
)

type Stat struct {
	gorm.Model
	VideoID  uint
	Username string
	Type     string
}

type PublicStat struct {
	ID        uint `gorm:"primarykey"`
	Time      time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	VideoID   uint
	UserId    uint
	Type      string
}

func AddStat(VID uint, username string, viewType string) Stat {
	st := Stat{
		VideoID:  VID,
		Username: username,
		Type:     viewType,
	}
	Db.Create(&st)
	return st
}

func AddPublicStat(VID uint, userId uint, viewType string) PublicStat {
	st := PublicStat{
		VideoID: VID,
		UserId:  userId,
		Type:    viewType,
		Time:    time.Now().UTC(),
	}
	Db.Create(&st)
	return st
}

func ListStats() []Stat {
	var stats []Stat
	Db.Find(&stats)
	return stats
}
