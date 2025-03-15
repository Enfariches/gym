package api

import (
	"context"
	"health/src/health/models"
	u "health/src/health/utils"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)


var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" {
			response = u.Message(u.MissingToken, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		dividedHeader := strings.Split(tokenHeader, " ")
		if len(dividedHeader) != 2 {
			response = u.Message(u.InvalidToken, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		tokenPart := dividedHeader[1] //Получаем вторую часть токена
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("TOKEN_SECRET")), nil
		})

		if err != nil {
			response = u.Message(u.InvalidToken, "Malformed authentication token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		if !token.Valid {
			response = u.Message(u.InvalidToken, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		session := models.GetSession(token.Raw)
		if session == nil {
			response = u.Message(u.InvalidToken, "Not logged in")
			w.WriteHeader(http.StatusForbidden)
			u.Respond(w, response)
			return
		}

		log.Println("User %", tk.Email)
		ctx := context.WithValue(r.Context(), "session", session)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
