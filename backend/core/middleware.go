package core

import ( 
	"net/http"
)


type Middleware struct {
	http.ResponseWriter
	StatusCode int
	Hijacked bool
}
