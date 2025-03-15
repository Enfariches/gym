package dto

import (
	"health/src/health/models"
	"time"
)

type AccountDto struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
}

type SessionDto struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Token     string     `json:"token"`
	Email     string     `json:"email"`
}

type EmailConfirmDto struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type StatDto struct {
	ID      uint      `json:"id"`
	VideoId uint      `json:"video_id"`
	UserId  uint      `json:"user_id"`
	Type    string    `json:"type"`
	Time    time.Time `json:"created_at"`
}

type StatDetailDto struct {
	UserId       uint   `json:"user_id"`
	Organization string `json:"organization"`
	Position     string `json:"position"`
	Department   string `json:"department"`
	Name         string `json:"name"`
	Count        uint   `json:"count"`
	Type         string `json:"type"`
}

func SessionToDto(session *models.Session) *SessionDto {
	return &SessionDto{
		ID:        session.ID,
		CreatedAt: session.CreatedAt,
		UpdatedAt: session.UpdatedAt,
		DeletedAt: session.DeletedAt,
		Token:     session.Token,
		Email:     session.Email,
	}
}

func EmailConfirmToDto(emailConfirm *models.EmailConfirm) *EmailConfirmDto {
	return &EmailConfirmDto{
		ID:    emailConfirm.ID,
		Email: emailConfirm.Email,
	}
}
