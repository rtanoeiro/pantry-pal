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
		Db:       database.New(newDB),
		Renderer: api.MyTemplates(),
		Env:      "dev",
		Secret:   jwtSecret,
	}

	httpServerMux := http.NewServeMux()

	httpServerMux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))

	// Reset - Used in dev for testing
	httpServerMux.Handle("POST /reset", http.HandlerFunc(config.ResetUsers))

	//Login
	httpServerMux.Handle("/", http.HandlerFunc(config.Index))
	httpServerMux.Handle("POST /login", http.HandlerFunc(config.Login))
	httpServerMux.Handle("GET /signup", http.HandlerFunc(config.SignUp))
	httpServerMux.Handle("POST /signup", http.HandlerFunc(config.CreateUser))
	httpServerMux.Handle("GET /home", http.HandlerFunc(config.Home))
	httpServerMux.Handle("GET /logout", http.HandlerFunc(config.Logout))

	// Users endpoints
	httpServerMux.Handle("GET /user", http.HandlerFunc(config.GetUserInfo))
	httpServerMux.Handle("POST /admin", http.HandlerFunc(config.UserAdmin))
	httpServerMux.Handle("POST /user/email", http.HandlerFunc(config.UpdateUserEmail))
	httpServerMux.Handle("POST /user/name", http.HandlerFunc(config.UpdateUserName))
	httpServerMux.Handle("POST /user/password", http.HandlerFunc(config.UpdateUserPassword))

	// Pantry endpoints
	httpServerMux.Handle("POST /pantry", http.HandlerFunc(config.HandleNewItem))
	httpServerMux.Handle("GET /pantry/{itemName}", http.HandlerFunc(config.GetItemByName))
	httpServerMux.Handle("GET /pantry/", http.HandlerFunc(config.GetAllPantryItems))
	httpServerMux.Handle("GET /pantry_stats", http.HandlerFunc(config.GetPantryStats))
	httpServerMux.Handle("DELETE /pantry/{itemID}", http.HandlerFunc(config.DeleteItem))

	httpServer := http.Server{
		Handler:           httpServerMux,
		Addr:              ":" + port,
		ReadHeaderTimeout: 60 * time.Second,
	}
	fmt.Println("Server started on port", port)

	_ = httpServer.ListenAndServe()
}
