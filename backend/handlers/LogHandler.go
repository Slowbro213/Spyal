package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type LogRequest struct {
	Level string `json:"level"`
	Msg   string `json:"msg"`
}

type LogHandler struct {
	l *zap.Logger
}

func NewLogHandler(l *zap.Logger) *LogHandler {
	return &LogHandler{l: l}
}


func (lh *LogHandler) LogFrontend(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight (CORS OPTIONS)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var logReq LogRequest
	if err := json.NewDecoder(r.Body).Decode(&logReq); err != nil {
		http.Error(w, "Bad Request: invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	logReq.Msg = "Frontend: " + logReq.Msg
	switch logReq.Level {
	case "info":
		lh.l.Info(logReq.Msg)
	case "warn":
		lh.l.Warn(logReq.Msg)
	case "error":
		lh.l.Error(logReq.Msg)
	default:
		http.Error(w, "Bad Request: invalid level", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

