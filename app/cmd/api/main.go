package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/CarlosZambonii/devforge/internal/crypto"
	"github.com/CarlosZambonii/devforge/internal/handler"
	"github.com/CarlosZambonii/devforge/internal/repository"
	"github.com/CarlosZambonii/devforge/internal/service"
)

func main() {
	key := []byte("12345678901234567890123456789012")

	aes, err := crypto.NewAESCrypto(key)
	if err != nil {
		log.Fatalf("erro ao criar crypto: %v", err)
	}

	repo := repository.NewURLRepository("localhost:6379")
	svc := service.NewURLService(repo, aes)
	h := handler.NewURLHandler(svc)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"status":"ok"}`)
	})

	mux.HandleFunc("POST /shorten", h.Shorten)
	mux.HandleFunc("GET /{code}", h.Resolve)
	mux.HandleFunc("DELETE /{code}", h.Delete)

	log.Println("DevForge API rodando na porta 8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
