package controllers

import (
	"encoding/json"
	"health/src/health/dto"
	"health/src/health/models"
	u "health/src/health/utils"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var CreateAdmin = func(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		u.Respond(w, u.Message(u.InvalidBody, "Invalid request body"))
		return
	}
	admin := models.CreateAdmin(req.Email)
	u.Respond(w, admin)
}

var GetSession = func(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value("session").(*models.Session)

	u.Respond(w, dto.SessionToDto(session))
}

var Logout = func(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value("session").(*models.Session)
	models.DeleteSessionById(session.ID)
	u.Respond(w, u.Message(u.Ok, "Logged out"))
}

var SendConfirmationCode = func(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		u.Respond(w, u.Message(u.InvalidBody, "Invalid request body"))
		return
	}
	if !models.HasAdmin(req.Email) { //Нужно быть админом, чтобы получить код
		w.WriteHeader(http.StatusBadRequest)
		u.Respond(w, u.Message(u.WrongEmail, "This email address is not on the allowed list"))
		return
	}

	emailConfirm := models.CreateEmailConfirm(req.Email)
	code, err := u.SendEmailCode(emailConfirm.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.Respond(w, u.Message(u.InternalError, err.Error()))
		return
	}
	emailConfirm = models.UpdateEmailConfirm(emailConfirm, code)
	u.Respond(w, dto.EmailConfirmToDto(emailConfirm))
}

var ConfirmEmail = func(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var req struct {
		Code uint
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		u.Respond(w, u.Message(u.InvalidBody, "Invalid request body"))
		return
	}

	emailConfirm := models.GetEmailConfirm(id)
	if emailConfirm == nil {
		w.WriteHeader(http.StatusNotFound)
		u.Respond(w, u.Message(u.EmailConfirmNotFound, "Email confirmation with this id doesn't exist"))
		return
	}
	if emailConfirm.Code != req.Code {
		w.WriteHeader(http.StatusBadRequest)
		u.Respond(w, u.Message(u.WrongConfirmationCode, "Wrong confirmation code"))
		return
	}

	sessionId := uuid.New().String()
	tk := new(models.Token)
	tk.Email = emailConfirm.Email
	tk.SessionId = sessionId
	tk.ExpiresAt = time.Now().UTC().Add(5 * time.Hour).Unix()

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))

	session, _ := models.AddSession(sessionId, tokenString, emailConfirm.Email)

	models.DeleteEmailConfirm(emailConfirm.Email)
	u.Respond(w, dto.SessionToDto(session))
}
