package main

import (
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"
	"log"
	"net/http"
	"os"
	"song-library-test-task/internal/external"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	httptransport "song-library-test-task/internal/handler/http"
	"song-library-test-task/internal/handler/http/endpoints"
	"song-library-test-task/internal/models"
	"song-library-test-task/internal/repository/postgres"
	"song-library-test-task/internal/service"
)

// @title           Song Library API
// @version         1.0
// @description     This is an example service for managing songs.
// @host            localhost:8080
// @BasePath        /
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("[WARN] no .env file found")
	}

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASS", "")
	dbName := getEnv("DB_NAME", "songsdb")
	extAPI := getEnv("EXTERNAL_API_BASE_URL", "http://localhost:3000")

	// Connect to DB
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("[ERROR] Could not open DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("[ERROR] Could not connect to DB: %v", err)
	}
	log.Println("[INFO] Connected to Postgres")

	goose.SetBaseFS(nil)
	migrationsDir := "./db/migrations"

	// 2. Run the migrations
	if err := goose.Up(db, migrationsDir); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	log.Println("[INFO] Migrations applied successfully")
	// Initialize repository
	var repo models.SongRepository = postgres.NewSongRepository(db)

	// Initialize external client
	externalClient := external.NewMusicInfoClient(extAPI, 5*time.Second)

	// Initialize service
	svc := service.NewSongService(repo, externalClient)

	// Build endpoints
	eps := endpoints.MakeSongEndpoints(*svc)

	// Create HTTP handler
	handler := httptransport.NewHTTPHandler(eps)

	// Start server
	addr := ":8080"
	log.Printf("[INFO] Listening on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("[ERROR] %v", err)
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
