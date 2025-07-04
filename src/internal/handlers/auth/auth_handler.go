package auth

import (
	"backend2/internal/auth"
	"backend2/internal/utils"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	a *auth.AuthUsecase
}

func NewAuthHandler(a *auth.AuthUsecase) *AuthHandler {
	return &AuthHandler{a: a}
}

func (h *AuthHandler) GetToken(w http.ResponseWriter, r *http.Request) {

	uuid, _ := utils.GenerateUUID()
	userID := "user-" + uuid

	token, err := h.a.GenerateToken(userID)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	resp := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
