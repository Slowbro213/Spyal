package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"spyal/auth"
	"spyal/core"
	"spyal/dto/requests"
	"spyal/models"
	"spyal/pkg/pages"
	"spyal/pkg/utils/renderer"
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

func (uh *UserHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	props := map[string]any{"Next": r.URL.Query().Get("next")}
	_ = renderer.Render(w, pages.IsFragment(r), props,
		pages.LayoutBase,
		pages.PageLogin,
	)
}

func (uh *UserHandler) LoginOrRegister(w http.ResponseWriter, r *http.Request) {
	form, err := requests.NewLoginForm(r)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	username := form.Username
	password := form.Password

	user, err := uh.userRepository.GetByUsername(ctx, username)
	if err != nil && !errors.Is(err, repos.ErrUserNotFound) {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		uh.Log.Error("failed to fetch user", zap.Error(err))
		return
	}

	if user == nil {
		user, err = uh.createUser(ctx, username, password)
	} else {
		err = uh.authenticateUser(user, password)
	}
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	uh.loginUser(w, user)

	next := r.URL.Query().Get("next")
	if next == "" || !strings.HasPrefix(next, "/") {
		next = "/"
	}
	http.Redirect(w, r, next, http.StatusSeeOther)
}

func (uh *UserHandler) Logout(w http.ResponseWriter, _ *http.Request) {
	prod := os.Getenv("ENV") == "production"

	http.SetCookie(w, &http.Cookie{
		Name:     "auth",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   prod,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    "",
		Path:     "/",
		HttpOnly: false,
		Secure:   prod,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`{"message":"logged out successfully"}`)); err != nil {
		uh.Log.Error("error writing response", zap.Error(err))
	}
}

func (uh *UserHandler) createUser(ctx context.Context, username, password string) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		uh.Log.Error("hash error", zap.Error(err))
		return nil, errors.New("sorry, something went wrong")
	}
	u := &models.User{Username: username, Password: string(hash)}
	if err := uh.userRepository.Create(ctx, u); err != nil {
		uh.Log.Error("create error", zap.Error(err))
		return nil, errors.New("sorry, something went wrong")
	}
	return u, nil
}

func (uh *UserHandler) authenticateUser(u *models.User, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return errors.New("invalid credentials")
	}
	return nil
}

func (uh *UserHandler) loginUser(w http.ResponseWriter, u *models.User) {
	token := auth.CreateToken(u.ID, u.Username, auth.TokenTTL*time.Second)
	prod := os.Getenv("ENV") == "production"

	http.SetCookie(w, &http.Cookie{
		Name:     "auth",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   prod,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(auth.TokenTTL),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    u.Username,
		Path:     "/",
		HttpOnly: false,
		Secure:   prod,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(auth.TokenTTL),
	})
}

func writeJSONError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
