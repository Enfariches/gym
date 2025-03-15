package models

import (
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	Name     string
	Archived bool
}

func AddVideo(filename string) (*Video, error) {
	v := Video{
		Name: filename,
	}
	Db.Create(&v)
	return &v, nil
}
func ToggleVid(id uint) error {
	v := Video{}
	v.ID = id
	Db.First(&v)
	v.Archived = !v.Archived
	Db.Save(&v)
	return nil
}
func UpdVid(id uint, name string) (*Video, error) {
	v := Video{}
	v.ID = id
	Db.Model(&v).Update("name", name)
	return &v, nil
}

func GetVideosList() []Video {
	var videos []Video
	Db.Find(&videos)
	return videos
}

func DeleteVideo(id uint) {
	Db.Where("id = ?", id).Delete(&Video{})
}
