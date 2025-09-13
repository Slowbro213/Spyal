package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"spyal/auth"
	"spyal/core"
	"spyal/models"
	"spyal/repos"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	core.Handler
	userRepository repos.UserRepository
}

func NewUserHandler(l *zap.Logger, ur repos.UserRepository) *UserHandler {
	return &UserHandler{
		Handler:        core.Handler{Log: l},
		userRepository: ur,
	}
}

func (uh *UserHandler) LoginOrRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	if username == "" || password == "" {
		http.Error(w, "Username and Password required", http.StatusBadRequest)
		return
	}

	user, err := uh.userRepository.GetByUsername(ctx, username)
	if err != nil && !errors.Is(err, repos.ErrUserNotFound) {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		uh.Log.Error("failed to fetch user", zap.Error(err))
		return
	}

	if user == nil {
		passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Sorry, something went wrong!", http.StatusInternalServerError)
			uh.Log.Error("error hashing password", zap.Error(err))
			return
		}

		user = &models.User{
			Username: username,
			Password: string(passHash),
		}

		if err := uh.userRepository.Create(ctx, user); err != nil {
			http.Error(w, "Sorry, something went wrong!", http.StatusInternalServerError)
			uh.Log.Error("error creating user", zap.Error(err))
			return
		}
	} else {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
	}

	token := auth.CreateToken(user.ID, user.Username, auth.TokenTTL*time.Second)

	prod := os.Getenv("ENV") == "production"
	http.SetCookie(w, &http.Cookie{
		Name:     "auth",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   prod,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(auth.TokenTTL * time.Second),
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := fmt.Sprintf(`{"ttl":%d}`,int(auth.TokenTTL))
	if _, err = w.Write([]byte(resp)); err != nil {
		uh.Log.Error("error writing response", zap.Error(err))
	}
}

func (uh *UserHandler) Logout(w http.ResponseWriter, _ *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`{"message":"logged out successfully"}`)); err != nil {
		uh.Log.Error("error writing response", zap.Error(err))
	}
}
