package jwt

import (
	"context"
	"fmt"
	"health/lib/ctxkey"
	"net/http"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)


var PublicMethods = map[string]bool{
	"/auth.AuthService/Login":                true,
	"/auth.AuthService/Register":             true,
	"/auth.AuthService/ChangePassword":       true,
	"/auth.AuthService/VerifyChangePassword": true,
	"/auth.AuthService/VerifyRegister":       true,
}

func JWTServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if PublicMethods[info.FullMethod] {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	authHeader := md.Get("Authorization")
	if len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing Authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader[0], "Bearer ")

	userId, err := ParseToken(tokenString)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	ctx = context.WithValue(ctx, ctxkey.UserKey, userId)

	return handler(ctx, req)
}
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Получаем токен из заголовка
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Формат: "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			http.Error(w, "invalid token format", http.StatusUnauthorized)
			return
		}

		// Парсим токен
		userID, err := ParseToken(tokenString) // Ваша функция парсинга
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// Сохраняем userID в контекст
		ctx := context.WithValue(r.Context(), ctxkey.UserKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
