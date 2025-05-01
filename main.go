package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"pantry-pal/pantry/api"
	"pantry-pal/pantry/database"

	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	envError := godotenv.Load()
	if envError != nil {
		log.Fatal("Unable to load environment variables")
	}

	dbURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	fmt.Println(dbURL)
	newDB, dbError := sql.Open("libsql", dbURL)
	if dbError != nil {
		log.Fatal("Unable to open database. Closing app. Error: ", dbError)
	}
	defer newDB.Close()

	errPing := newDB.Ping()
	if errPing != nil {
		log.Fatalf("Failed to ping the database: %v", errPing)
	}

	config := api.Config{
		Db:  database.New(newDB),
		Env: "dev",
	}

	fmt.Println("Connected to the database successfully")
	httpServerMux := http.NewServeMux()
	baseLoginPage := http.FileServer(http.Dir("static"))
	httpServerMux.Handle("/", baseLoginPage)
	httpServerMux.Handle("POST /api/users", http.HandlerFunc(config.CreateUser))
	httpServerMux.Handle("GET /api/users/{userInfo}", http.HandlerFunc(config.GetUserInfo))
	httpServerMux.Handle("GET /api/reset", http.HandlerFunc(config.ResetUsers))

	httpServer := http.Server{
		Handler: httpServerMux,
		Addr:    ":" + port,
	}
	fmt.Println("Server started on port", port)

	httpServer.ListenAndServe()
}
