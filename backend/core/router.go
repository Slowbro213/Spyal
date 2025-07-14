package core

import (
	"log"
	"net/http"
	"strings"
)

type Router struct {
	routes map[string]map[string]http.HandlerFunc // method -> path -> handler
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]map[string]http.HandlerFunc),
	}
}


func (r *Router) Get(path string, handler http.HandlerFunc) {
	r.register(http.MethodGet, path, handler)
}

func (r *Router) Post(path string, handler http.HandlerFunc) {
	r.register(http.MethodPost, path, handler)
}

func (r *Router) Put(path string, handler http.HandlerFunc) {
	r.register(http.MethodPut, path, handler)
}

func (r *Router) Patch(path string, handler http.HandlerFunc) {
	r.register(http.MethodPatch, path, handler)
}

func (r *Router) Delete(path string, handler http.HandlerFunc) {
	r.register(http.MethodDelete, path, handler)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	http.DefaultServeMux.ServeHTTP(w, req)
}

func (r *Router) register(method, path string, handler http.HandlerFunc) {
    if _, exists := r.routes[path]; !exists {
        r.routes[path] = make(map[string]http.HandlerFunc)
        http.DefaultServeMux.HandleFunc(path, r.createHandler(path))
    }

    if _, exists := r.routes[path][method]; exists {
        panic("duplicate handler for " + method + " " + path)
    }

    r.routes[path][method] = handler
}

func (r *Router) createHandler(path string) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        methodHandlers := r.routes[path]

        if req.Method == http.MethodOptions {
            r.handleOptions(w, methodHandlers)
            return
        }

        h, ok := methodHandlers[req.Method]
        if !ok {
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
            return
        }

        log.Printf("%s %s", req.Method, req.URL.Path)
        r.executeHandler(h, w, req)
    }
}

func (r *Router) handleOptions(w http.ResponseWriter, methodHandlers map[string]http.HandlerFunc) {
    methods := make([]string, 0, len(methodHandlers))
    for m := range methodHandlers {
        methods = append(methods, m)
    }
    w.Header().Set("Allow", strings.Join(methods, ", "))
    w.WriteHeader(http.StatusNoContent)
}

func (r *Router) executeHandler(h http.HandlerFunc, w http.ResponseWriter, req *http.Request) {
    defer func() {
        if rec := recover(); rec != nil {
            log.Printf("panic recovered: %v", rec)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        }
    }()
    h(w, req)
}
