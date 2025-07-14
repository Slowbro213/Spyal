package pages 

import (
	"net/http"
)


const (
	LayoutBase   = "layouts/base.html"
	LayoutEmpty  = "layouts/empty.html"

	PageCreate   = "pages/create.html"
	PageRemote   = "pages/remote.html"
	PageHome     = "pages/index.html"
	PageRoom     = "pages/room.html"

	CompRoom     = "components/room.html"
)



func IsFragment(r *http.Request) bool {
	return r.Header.Get("X-Smart-Link") == "true"
} 
