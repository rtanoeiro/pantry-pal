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
	httpServerMux.Handle("POST /reset", http.HandlerFunc(config.ResetUsers))

	//Login
	httpServerMux.Handle("GET /", http.FileServer(http.Dir("static")))
	httpServerMux.HandleFunc("GET /signup", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/signup.html")
	})

	// After Login Pages
	httpServerMux.HandleFunc("GET /home", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/home.html")
	})
	httpServerMux.Handle("POST /login", http.HandlerFunc(config.Login))

	// Users endpoints
	httpServerMux.Handle("POST /users", http.HandlerFunc(config.CreateUser))
	httpServerMux.Handle("GET /users/", http.HandlerFunc(config.GetUserInfo))
	httpServerMux.Handle("POST /admin", http.HandlerFunc(config.UserAdmin))
	httpServerMux.Handle("PATCH /users/email", http.HandlerFunc(config.UpdateUserEmail))
	httpServerMux.Handle("PATCH /users/name", http.HandlerFunc(config.UpdateUserName))
	httpServerMux.Handle("PATCH /users/password", http.HandlerFunc(config.UpdateUserPassword))

	// Pantry endpoints
	httpServerMux.Handle("POST /pantry", http.HandlerFunc(config.HandleNewItem))
	httpServerMux.Handle("GET /pantry/{itemName}", http.HandlerFunc(config.GetItemByName))
	httpServerMux.Handle("GET /pantry/", http.HandlerFunc(config.GetAllPantryItems))
	httpServerMux.Handle("DELETE /pantry/{itemID}", http.HandlerFunc(config.DeleteItem))

	httpServer := http.Server{
		Handler:           httpServerMux,
		Addr:              ":" + port,
		ReadHeaderTimeout: 60 * time.Second,
	}
	fmt.Println("Server started on port", port)

	_ = httpServer.ListenAndServe()
}
