package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/CarlosZambonii/devforge/internal/crypto"
	"github.com/CarlosZambonii/devforge/internal/handler"
	"github.com/CarlosZambonii/devforge/internal/repository"
	"github.com/CarlosZambonii/devforge/internal/service"
	"github.com/CarlosZambonii/devforge/pkg/vault"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func main() {
	// New Relic (opcional: so ativa se a license key existir)
	var nrApp *newrelic.Application
	if licenseKey := os.Getenv("NEW_RELIC_LICENSE_KEY"); licenseKey != "" {
		app, err := newrelic.NewApplication(
			newrelic.ConfigAppName("DevForge API"),
			newrelic.ConfigLicense(licenseKey),
			newrelic.ConfigDistributedTracerEnabled(true),
		)
		if err != nil {
			log.Printf("aviso: New Relic nao iniciou: %v", err)
		} else {
			nrApp = app
			log.Println("New Relic ativo")
		}
	}

	// Vault
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

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	repo := repository.NewURLRepository(redisAddr)
	svc := service.NewURLService(repo, aes)
	h := handler.NewURLHandler(svc)

	mux := http.NewServeMux()

	// helper: registra rota com instrumentacao New Relic (se ativo)
	handle := func(pattern string, fn http.HandlerFunc) {
		if nrApp != nil {
			_, h := newrelic.WrapHandleFunc(nrApp, pattern, fn)
			mux.HandleFunc(pattern, h)
		} else {
			mux.HandleFunc(pattern, fn)
		}
	}

	handle("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"status":"ok"}`)
	})
	handle("POST /shorten", h.Shorten)
	handle("GET /{code}", h.Resolve)
	handle("DELETE /{code}", h.Delete)

	log.Println("DevForge API rodando na porta 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
