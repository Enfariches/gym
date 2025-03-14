package models

import (
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

type Token struct {
	Email     string
	SessionId string
	jwt.StandardClaims
}

type Session struct {
	ID        string `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Token     string
	Email     string
}

func AddSession(id string, token string, email string) (*Session, error) {
	session := Session{
		ID:    id,
		Token: token,
		Email: email,
	}
	Db.Create(&session)
	return &session, nil
}

func GetSession(token string) *Session {
	session := &Session{}
	err := Db.Table("sessions").Where("token = ?", token).First(&session).Error
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return session
}

func DeleteSessionByToken(token string) { //Не использутеся
	Db.Table("sessions").Where("token = ?", token).Unscoped().Delete(&Session{})
}

func DeleteSessionById(id string) {
	Db.Table("sessions").Where("id = ?", id).Unscoped().Delete(&Session{})
}
