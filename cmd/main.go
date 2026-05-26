package main

import (
	"log"
	"os"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := server.New(":" + port)
	log.Printf("sports statistics site is available on http://localhost:%s", port)
	log.Fatal(srv.ListenAndServe())
}
