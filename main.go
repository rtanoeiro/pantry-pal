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
	httpServerMux.Handle("GET /api/reset", http.HandlerFunc(config.ResetUsers))

	httpServerMux.Handle("POST /api/users", http.HandlerFunc(config.CreateUser))
	httpServerMux.Handle("PATCH /api/users/{userID}/email", http.HandlerFunc(config.UpdateUserEmail))
	httpServerMux.Handle("PATCH /api/users/{userID}/name", http.HandlerFunc(config.UpdateUserName))
	httpServerMux.Handle("PATCH /api/users/{userID}/password", http.HandlerFunc(config.UpdateUserPassword))

	httpServerMux.Handle("GET /api/users/{userInfo}", http.HandlerFunc(config.GetUserInfo))

	httpServerMux.Handle("POST /api/pantry", http.HandlerFunc(config.HandleNewItem))
	httpServerMux.Handle("GET /api/pantry/{itemName}", http.HandlerFunc(config.GetItemByName))
	httpServerMux.Handle("GET /api/pantry/", http.HandlerFunc(config.GetAllPantryItems))
	httpServerMux.Handle("DELETE /api/pantry/{itemID}", http.HandlerFunc(config.DeleteItem))

	//httpServerMux.Handle("GET /api/pantry/{userID}", http.HandlerFunc(config.GetAllUserItems))

	httpServer := http.Server{
		Handler:           httpServerMux,
		Addr:              ":" + port,
		ReadHeaderTimeout: 60 * time.Second,
	}
	fmt.Println("Server started on port", port)

	_ = httpServer.ListenAndServe()
}
