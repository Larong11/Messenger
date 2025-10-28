package handlers

import (
	"encoding/json"
	"net/http"
	packageuserus "server/application/use_cases/user"
	upgradeerrors "server/internal/errors"
)

type UserHandler struct {
	registerUserUseCases *packageuserus.RegisterUserUseCases
}

func NewUserHandler(registerUserUseCases *packageuserus.RegisterUserUseCases) *UserHandler {
	return &UserHandler{
		registerUserUseCases,
	}
}
func (h *UserHandler) CheckUserName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Структура для чтения JSON
	var req struct {
		Username string `json:"username"`
	}

	// Декодируем JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	username := req.Username
	available, err := h.registerUserUseCases.CheckUserName(ctx, username)
	if err != nil {
		upgradeerrors.HandleHTTPError(w, err)
		return
	}
	resp := map[string]bool{"available": available}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (h *UserHandler) CheckEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
	var req struct {
		Email string `json:"email"`
	}
	// Декодируем JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	available, err := h.registerUserUseCases.CheckEmail(ctx, req.Email)
	if err != nil {
		upgradeerrors.HandleHTTPError(w, err)
		return
	}
	resp := map[string]bool{"available": available}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
	var req struct {
		FirstName        string `json:"first_name"`
		LastName         string `json:"last_name"`
		UserName         string `json:"username"`
		Email            string `json:"email"`
		Password         string `json:"password"`
		AvatarURL        string `json:"avatar_url"`
		VerificationCode string `json:"verification_code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	ID, err := h.registerUserUseCases.RegisterUser(ctx, req.FirstName, req.LastName, req.UserName, req.Email, req.Password, req.AvatarURL, req.VerificationCode)
	if err != nil {
		upgradeerrors.HandleHTTPError(w, err)
		return
	}
	resp := map[string]int{"id": *ID}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (h *UserHandler) RequestVerificationCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	err := h.registerUserUseCases.RequestVerificationCode(ctx, req.Email)
	if err != nil {
		upgradeerrors.HandleHTTPError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
