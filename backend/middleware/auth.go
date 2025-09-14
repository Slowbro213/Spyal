package middleware

import (
	"context"
	"net/http"
	"net/url"

	"spyal/auth"
)

//nolint
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawToken := ""
		if c, err := r.Cookie("auth"); err == nil {
			rawToken = c.Value
		}
		if rawToken == "" {
			redirectToLogin(w, r)
			return
		}

		id, username, ok := auth.VerifyToken(rawToken)
		if !ok {
			redirectToLogin(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "username", username)
		ctx = context.WithValue(ctx, "id", int64(id))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	loginURL := "/login?next=" + url.QueryEscape(r.URL.RequestURI())
	http.SetCookie(w, &http.Cookie{
		Name:     "auth",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, loginURL, http.StatusSeeOther)
}
