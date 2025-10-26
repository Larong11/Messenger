package handlers

import (
	"encoding/json"
	"net/http"
	"server/application/use_cases/user"
)

type UserHandler struct {
	registerUserUseCases *user.RegisterUserUseCases
}

func NewUserHandler(registerUserUseCases *user.RegisterUserUseCases) *UserHandler {
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
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	resp := map[string]bool{"available": available}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
