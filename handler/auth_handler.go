package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"first-project/models"
	"first-project/service"

	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Post("/register", h.handleRegister)
	r.Post("/login", h.handleLogin)
}

func (h *AuthHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}

	newUser, err := h.service.Register(user)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			h.respondWithError(w, http.StatusConflict, "Email sudah terdaftar")
		} else if strings.Contains(err.Error(), "minimal 6 karakter") {
			h.respondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			h.respondWithError(w, http.StatusInternalServerError, "Gagal mendaftarkan user")
		}
		return
	}

	h.respondWithJSON(w, http.StatusCreated, newUser)
}

func (h *AuthHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}

	token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		h.respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *AuthHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

func (h *AuthHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}
