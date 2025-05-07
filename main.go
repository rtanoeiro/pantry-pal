package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"pantry-pal/pantry/api"
	"pantry-pal/pantry/database"
	"time"

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
	jwtSecret := os.Getenv("JWT_SECRET")

	newDB, dbError := sql.Open("libsql", dbURL)
	if dbError != nil {
		log.Fatal("Unable to open database. Closing app. Error: ", dbError)
	}
	defer newDB.Close()

	config := api.Config{
		Db:     database.New(newDB),
		Env:    "dev",
		Secret: jwtSecret,
	}

	fmt.Println("Connected to the database successfully")
	httpServerMux := http.NewServeMux()

	// Reset - Used in dev for testing
	httpServerMux.Handle("POST /api/reset", http.HandlerFunc(config.ResetUsers))

	//Login
	httpServerMux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/login.html")
	})
	httpServerMux.HandleFunc("GET /api/signup", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/signup.html")
	})

	httpServerMux.HandleFunc("GET /home", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/home.html")
	})
	httpServerMux.Handle("POST /api/login", http.HandlerFunc(config.Login))

	// Users endpoints
	httpServerMux.Handle("POST /api/users", http.HandlerFunc(config.CreateUser))
	httpServerMux.Handle("GET /api/users/", http.HandlerFunc(config.GetUserInfo))
	httpServerMux.Handle("POST /api/admin", http.HandlerFunc(config.UserAdmin))
	httpServerMux.Handle("PATCH /api/users/email", http.HandlerFunc(config.UpdateUserEmail))
	httpServerMux.Handle("PATCH /api/users/name", http.HandlerFunc(config.UpdateUserName))
	httpServerMux.Handle("PATCH /api/users/password", http.HandlerFunc(config.UpdateUserPassword))

	// Pantry endpoints
	httpServerMux.Handle("POST /api/pantry", http.HandlerFunc(config.HandleNewItem))
	httpServerMux.Handle("GET /api/pantry/{itemName}", http.HandlerFunc(config.GetItemByName))
	httpServerMux.Handle("GET /api/pantry/", http.HandlerFunc(config.GetAllPantryItems))
	httpServerMux.Handle("DELETE /api/pantry/{itemID}", http.HandlerFunc(config.DeleteItem))

	httpServer := http.Server{
		Handler:           httpServerMux,
		Addr:              ":" + port,
		ReadHeaderTimeout: 60 * time.Second,
	}
	fmt.Println("Server started on port", port)

	_ = httpServer.ListenAndServe()
}
