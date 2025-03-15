package api

import (
	"github.com/gorilla/handlers"
	"log"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return handlers.LoggingHandler(log.Writer(), next)
}
