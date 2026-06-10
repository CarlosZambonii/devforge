package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/CarlosZambonii/devforge/internal/service"
)

type URLHandler struct {
	service *service.URLService
}

func NewURLHandler(service *service.URLService) *URLHandler {
	return &URLHandler{service: service}
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	Code string `json:"code"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{"body invalido"})
		return
	}

	code, err := h.service.Shorten(context.Background(), req.URL)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, shortenResponse{code})
}

func (h *URLHandler) Resolve(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	rawURL, err := h.service.Resolve(context.Background(), code)
	if err != nil {
		writeJSON(w, http.StatusNotFound, errorResponse{"code nao encontrado"})
		return
	}

	http.Redirect(w, r, rawURL, http.StatusFound)
}

func (h *URLHandler) Delete(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	if err := h.service.Delete(context.Background(), code); err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
