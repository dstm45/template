// Package controllers contains all the controllers
package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/dstm45/template/pkg/services"
	"github.com/dstm45/template/pkg/utils"
)

type AuthController struct {
	S services.AuthService
}

func NewAuthController(service services.AuthService) *AuthController {
	controller := AuthController{
		S: service,
	}
	return &controller
}

type signinRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *AuthController) Signin(w http.ResponseWriter, r *http.Request) {
	var req signinRequest
	ctx := r.Context()
	logger := utils.GetLoggerFromContext(ctx)
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error("Erreur lors du décodage du payload json", "error", err)
		http.Error(w, "Erreur lors du décodage du payload json", http.StatusBadRequest)
		return
	}
	refreshCookie, accessCookie, err := c.S.SignIn(ctx, req.Email, req.Password)
	if err != nil {
		logger.Error("Erreur lors de la connection au compte de l'utilisateur", "error", err)
		http.Error(w, "Erreur lors de la connection au compte", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, refreshCookie)
	http.SetCookie(w, accessCookie)
	w.WriteHeader(http.StatusOK)
}

func (c *AuthController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	refreshCookie, accessCookie, err := c.S.RotateToken(ctx, r)
	if err != nil {
		http.Error(w, "Failed to refresh tokens", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, refreshCookie)
	http.SetCookie(w, accessCookie)
	w.WriteHeader(http.StatusOK)
}

func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}
	accessTokenString := cookie.Value
	c.S.DeleteToken(r.Context(), accessTokenString)
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	})
	w.WriteHeader(http.StatusOK)
}
