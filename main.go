package main

import (
	"log"
	"net/http"

	"SubGen/internal/config"
	"SubGen/internal/server"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Println("failed to load config.yaml:", err)
	}
	srv := server.New(cfg)
	mux := srv.Routes()
	log.Println("Server is running on port 7081")
	if err := http.ListenAndServe(":7081", mux); err != nil {
		log.Println("server error:", err)
	}
}
