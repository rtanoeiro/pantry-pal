package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	envError := godotenv.Load()

	if envError != nil {
		log.Fatal("Unable to load environment variables")
	}

	httpServerMux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("static"))
	httpServerMux.Handle("/", fileServer)

	httpServer := http.Server{
		Handler: httpServerMux,
		Addr:    ":8080",
	}
	httpServer.ListenAndServe()

}
