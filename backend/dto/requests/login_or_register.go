package requests

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type LoginForm struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func NewLoginForm(r *http.Request) (LoginForm, error) {
	var f LoginForm
	ct := r.Header.Get("Content-Type")

	switch {
	case strings.Contains(ct, "application/json"):
		if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
			return f, errors.New("invalid JSON")
		}
	default:
		if err := r.ParseForm(); err != nil {
			return f, errors.New("cannot parse form")
		}
		f.Username = r.PostFormValue("username")
		f.Password = r.PostFormValue("password")
	}
	return f.Validate()
}

func (f LoginForm) Validate() (LoginForm, error) {
	f.Username = strings.TrimSpace(f.Username)
	f.Password = strings.TrimSpace(f.Password)

	switch {
	case f.Username == "":
		return f, errors.New("username required")
	case f.Password == "":
		return f, errors.New("password required")
	case len(f.Username) > 127:
		return f, errors.New("username too long (max 127)")
	case len(f.Password) < 6:
		return f, errors.New("password must be ≥ 6 characters")
	}
	return f, nil
}
