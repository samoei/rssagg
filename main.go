package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/samoei/rssagg/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT not found in the environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB URL not found in the environment")
	}

	driverName := "postgres"
	dataSourceName := dbURL

	conn, err := sql.Open(driverName, dataSourceName)

	if err != nil {
		log.Fatal("Could not connect to the database", err)
	}

	apiConfig := apiConfig{
		DB: database.New(conn),
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthcheck", handleHealthCheck)
	v1Router.Get("/err", handleErr)
	v1Router.Post("/users", apiConfig.handlerCreateUser)
	v1Router.Get("/users", apiConfig.middlewareAuth(apiConfig.handleGetUser))

	v1Router.Post("/feeds", apiConfig.middlewareAuth(apiConfig.handlerCreateFeed))
	v1Router.Get("/feeds", apiConfig.handlerGetFeeds)
	v1Router.Post("/feed-follows", apiConfig.middlewareAuth(apiConfig.handleCreateFeedFollow))
	v1Router.Get("/feed-follows", apiConfig.middlewareAuth(apiConfig.handleGetFeedFollows))
	v1Router.Delete("/feed-follows/{feed-follow-id}", apiConfig.middlewareAuth(apiConfig.handleDeleteFeedFollow))

	router.Mount("/v1", v1Router)
	server := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	fmt.Printf("Server starting on port %v", portString)
	err = server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}

}
