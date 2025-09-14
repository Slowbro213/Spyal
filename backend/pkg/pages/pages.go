package pages

import (
	"net/http"
	"os"
	"errors"
)

//nolint:gochecknoglobals
var (
	LayoutBase  string
	LayoutEmpty string

	PageCreate  string
	PageRemote  string
	PageHome    string
	PageLogin    string
	PageGames   string
	PageRoom    string

	CompRoom    string
)

func Init() error {
	viewsDir, err := getenv("VIEWS_DIR")
	if err != nil {
		return err
	}

	LayoutBase = viewsDir + "layouts/base.html"
	LayoutEmpty = viewsDir + "layouts/empty.html"

	PageCreate = viewsDir + "pages/create.html"
	PageRemote = viewsDir + "pages/remote.html"
	PageHome = viewsDir + "pages/index.html"
	PageLogin = viewsDir + "pages/login.html"
	PageGames = viewsDir + "pages/games.html"
	PageRoom = viewsDir + "pages/room.html"

	CompRoom = viewsDir + "components/room.html"


	return nil
}

func getenv(key string) (string,error) {
	if val := os.Getenv(key); val != "" {
		return val, nil
	}
	return "", errors.New("environment variable " + key + " not set")
}

func IsFragment(r *http.Request) bool {
	return r.Header.Get("X-Smart-Link") == "true"
}
