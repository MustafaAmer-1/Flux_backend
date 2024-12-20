package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MustafaAmer-1/Flux/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()

	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Fatal("Port is not found in the environment")
	}

	DB_URL := os.Getenv("DB_URL")
	if DB_URL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	dbconn, err_db := sql.Open("postgres", DB_URL)
	if err_db != nil {
		log.Fatal("Can't connect to database", err_db)
	}

	dbQ := database.New(dbconn)
	apiCfg := apiConfig{
		DB: dbQ,
	}

	go startScraping(dbQ, 10, time.Minute)

	router := chi.NewRouter()
	router.Use(LoggingMiddleware)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerError)
	v1Router.Put("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Post("/users", apiCfg.handlerLoginUser)
	v1Router.Put("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)
	v1Router.Post("/feed_follow", apiCfg.middlewareAuth(apiCfg.handlerFollowFeed))
	v1Router.Get("/feed_follow", apiCfg.middlewareAuth(apiCfg.handlerGetFollowedFeeds))
	v1Router.Get("/feed_not_follow", apiCfg.middlewareAuth(apiCfg.handlerGetNotFollowedFeeds))
	v1Router.Delete("/feed_follow/{feedId}", apiCfg.middlewareAuth(apiCfg.handlerUnFollowFeed))
	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPosts))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + PORT,
	}

	log.Printf("Server started on port %v\n", PORT)
	log.Fatal(srv.ListenAndServe())
}
