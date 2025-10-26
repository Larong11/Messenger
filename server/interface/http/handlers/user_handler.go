package handlers

import (
	"encoding/json"
	"net/http"
	package_user_us "server/application/use_cases/user"
)

type UserHandler struct {
	registerUserUseCases *package_user_us.RegisterUserUseCases
}

func NewUserHandler(registerUserUseCases *package_user_us.RegisterUserUseCases) *UserHandler {
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
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		UserName  string `json:"username"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	ID, err := h.registerUserUseCases.RegisterUser(ctx, req.FirstName, req.LastName, req.UserName, req.Email, req.Password, req.AvatarURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
