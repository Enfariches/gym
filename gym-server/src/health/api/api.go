package api

import (
	"fmt"
	"health/src/health/controllers"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func StartApi() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	router.PathPrefix("/upload").Handler(http.StripPrefix("/upload/", http.FileServer(http.Dir("./upload"))))
	routerV1 := router.PathPrefix("/api/v1").Subrouter()
	authRouter := routerV1.PathPrefix("/").Subrouter()
	modRouter := authRouter.PathPrefix("/").Subrouter()

	authRouter.HandleFunc("/session", controllers.GetSession).Methods("GET")
	authRouter.HandleFunc("/session/logout", controllers.Logout).Methods("GET")

	modRouter.HandleFunc("/schedule/{day}", controllers.AddScheduleHandler).Methods("POST")
	modRouter.HandleFunc("/schedule", controllers.UpdateScheduleHandler).Methods("PUT")
	modRouter.HandleFunc("/schedule/{id}", controllers.DeleteScheduleHandler).Methods("DELETE")

	modRouter.HandleFunc("/videos/{id}", controllers.DeleteVideoHandler).Methods("DELETE")
	modRouter.HandleFunc("/videos/{id}", controllers.UpdateVideoHandler).Methods("PUT")
	modRouter.HandleFunc("/videos/{id}", controllers.ToggleVideoHandler).Methods("GET")
	modRouter.HandleFunc("/videos", controllers.UploadVideo).Methods("POST")

	modRouter.HandleFunc("/admin", controllers.CreateAdmin).Methods("POST") //Админский метод
	routerV1.HandleFunc("/email/send", controllers.SendConfirmationCode).Methods("POST")
	routerV1.HandleFunc("/email/{id}/confirm", controllers.ConfirmEmail).Methods("POST")

	routerV1.HandleFunc("/export/video/{id}", controllers.ExportVideo).Methods("GET")

	routerV1.HandleFunc("/schedule/{day}", controllers.ListScheduleHandler).Methods("GET")
	routerV1.HandleFunc("/schedule", controllers.ListAllScheduleHandler).Methods("GET")

	routerV1.HandleFunc("/stats", controllers.AddStatHandler).Methods("POST")
	routerV1.HandleFunc("/stats", controllers.ListStatsHandler).Methods("GET")
	routerV1.HandleFunc("/stats/details", controllers.ListStatsDetailsHandler).Methods("GET")

	routerV1.HandleFunc("/videos", controllers.GetVideosListHandler).Methods("GET")

	authRouter.Use(JwtAuthentication)
	router.Use(LoggingMiddleware)
	router.Use(PrivateNetworkMiddleware)

	allowedHeaders := handlers.AllowedHeaders([]string{"Origin", "Authorization", "Content-Length", "Content-Type"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "PUT", "POST", "DELETE"})
	origins := handlers.AllowedOrigins(strings.Split(os.Getenv("ALLOW_ORIGINS"), ","))
	exposedHeaders := handlers.ExposedHeaders([]string{"Content-Length", "Content-Type", "Access-Control-Allow-Origin"})
	credentials := handlers.AllowCredentials()

	err = http.ListenAndServe(":8083", handlers.CORS(exposedHeaders, allowedHeaders, allowedMethods, origins, credentials)(router))
	if err != nil {
		fmt.Println(err)
	}
}

var PrivateNetworkMiddleware = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Request-Private-Network", "true")
		w.Header().Add("Access-Control-Allow-Private-Network", "true")
		next.ServeHTTP(w, r)
	})
}
