package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"pantry-pal/pantry/api"
	"pantry-pal/pantry/database"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	jwtSecret := os.Getenv("JWT_SECRET")

	newDB, dbError := sql.Open("sqlite3", dbURL)
	if dbError != nil {
		log.Fatal("Unable to open database. Closing app. Error: ", dbError)
	}
	defer api.CloseDB(newDB)

	config := api.Config{
		Db:       database.New(newDB),
		Port:     port,
		DBUrl:    dbURL,
		Renderer: api.MyTemplates("static/*.html"),
		Env:      "dev",
		Secret:   jwtSecret,
	}

	httpServerMux := http.NewServeMux()

	httpServerMux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))

	// Login
	httpServerMux.Handle("GET /login", http.HandlerFunc(config.Index))
	httpServerMux.Handle("POST /login", http.HandlerFunc(config.Login))
	httpServerMux.Handle("GET /signup", http.HandlerFunc(config.SignUp))
	httpServerMux.Handle("POST /signup", http.HandlerFunc(config.CreateUser))
	httpServerMux.Handle("GET /home", http.HandlerFunc(config.Home))
	httpServerMux.Handle("GET /logout", http.HandlerFunc(config.Logout))

	// Users endpoints
	httpServerMux.Handle("GET /user", http.HandlerFunc(config.GetUserInfo))
	httpServerMux.Handle("DELETE /user/{UserID}", http.HandlerFunc(config.DeleteUser))
	httpServerMux.Handle("POST /user/admin/{UserID}", http.HandlerFunc(config.AddUserAdmin))
	httpServerMux.Handle("DELETE /user/admin/{UserID}", http.HandlerFunc(config.RevokeUserAdmin))
	httpServerMux.Handle("POST /user/name", http.HandlerFunc(config.UpdateUserName))
	httpServerMux.Handle("POST /user/password", http.HandlerFunc(config.UpdateUserPassword))

	// Pantry endpoints
	httpServerMux.Handle("POST /pantry", http.HandlerFunc(config.HandleNewItem))
	httpServerMux.Handle("GET /pantry", http.HandlerFunc(config.GetAllPantryItems))
	httpServerMux.Handle("POST /pantry/addone/{ItemID}", http.HandlerFunc(config.HandleAddOnePantry))
	httpServerMux.Handle("POST /pantry/removeone/{ItemID}", http.HandlerFunc(config.HandleRemoveOnePantry))
	httpServerMux.Handle("GET /expiring", http.HandlerFunc(config.RenderExpiringSoon))
	httpServerMux.Handle("DELETE /pantry/{ItemID}", http.HandlerFunc(config.DeleteItem))

	// Shopping Cart endpoints
	httpServerMux.Handle("GET /shopping", http.HandlerFunc(config.RenderShoppingCart))
	httpServerMux.Handle("POST /shopping", http.HandlerFunc(config.AddItemShopping))
	httpServerMux.Handle("POST /shopping/addone", http.HandlerFunc(config.AddOneItemShopping))
	httpServerMux.Handle("POST /shopping/removeone", http.HandlerFunc(config.RemoveOneItemShopping))
	httpServerMux.Handle("DELETE /shopping/{itemName}", http.HandlerFunc(config.RemoveItemShopping))

	httpServer := http.Server{
		Handler:           httpServerMux,
		Addr:              ":" + port,
		ReadHeaderTimeout: 60 * time.Second,
	}
	fmt.Println("Server started on port", port)

	_ = httpServer.ListenAndServe()
}
