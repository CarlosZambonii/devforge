package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/CarlosZambonii/devforge/internal/crypto"
	"github.com/CarlosZambonii/devforge/internal/handler"
	"github.com/CarlosZambonii/devforge/internal/repository"
	"github.com/CarlosZambonii/devforge/internal/service"
	"github.com/CarlosZambonii/devforge/pkg/vault"
)

func main() {
	// busca a chave AES do Vault (nao mais hardcoded)
	vaultClient, err := vault.NewClient()
	if err != nil {
		log.Fatalf("erro ao conectar no vault: %v", err)
	}

	key, err := vaultClient.GetAESKey()
	if err != nil {
		log.Fatalf("erro ao buscar chave AES no vault: %v", err)
	}

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
