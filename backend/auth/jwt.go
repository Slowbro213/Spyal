package auth

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	TokenTTL = 100
)

//nolint
var tokenSecret []byte

type ctxKey string

const usernameKey ctxKey = "username"

//nolint
func init() {
	secret := os.Getenv("TOKEN_SECRET")
	if secret != "" {
		tokenSecret = []byte(secret)
		return
	}

	tokenLength := 32
	b := make([]byte, tokenLength)
	if _, err := rand.Read(b); err != nil {
		log.Fatalf("failed to generate token secret: %s",err)
	}
	tokenSecret = []byte(base64.RawURLEncoding.EncodeToString(b))
}


func CreateToken(username string, ttl time.Duration) string {
	exp := time.Now().Add(ttl).Unix()
	payload := fmt.Sprintf("%s.%d", username, exp)

	mac := hmac.New(sha256.New, tokenSecret)
	mac.Write([]byte(payload))
	signature := mac.Sum(nil)

	return fmt.Sprintf("%s.%s", base64.RawURLEncoding.EncodeToString([]byte(payload)), base64.RawURLEncoding.EncodeToString(signature))
}

func VerifyToken(token string) (string, bool) {
	parts := strings.Split(token, ".")
	tokenLength := 2
	if len(parts) != tokenLength {
		return "", false
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", false
	}

	signature, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", false
	}

	mac := hmac.New(sha256.New, tokenSecret)
	mac.Write(payloadBytes)
	expectedSig := mac.Sum(nil)

	if !hmac.Equal(signature, expectedSig) {
		return "", false
	}

	parts2 := strings.Split(string(payloadBytes), ".")
	if len(parts2) != tokenLength {
		return "", false
	}

	username := parts2[0]
	expUnix, err := strconv.ParseInt(parts2[1], 10, 64)
	if err != nil {
		return "", false
	}

	if time.Now().Unix() > expUnix {
		return "", false
	}

	return username, true
}

func WithUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, usernameKey, username)
}

func UsernameFromContext(ctx context.Context) (string, bool) {
	u, ok := ctx.Value(usernameKey).(string)
	return u, ok
}
