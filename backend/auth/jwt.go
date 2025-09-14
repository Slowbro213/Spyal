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
	TokenTTL = 60*60
)

//nolint
var tokenSecret []byte

type ctxKey string

const usernameKey ctxKey = "username"

//nolint
func init() {
	secret := os.Getenv("TOKEN_SECRET")
	fmt.Println("secret is : "+secret)
	if secret != "" {
		tokenSecret = []byte(secret)
		return
	}

	fmt.Println("secret env is empty")
	tokenLength := 32
	b := make([]byte, tokenLength)
	if _, err := rand.Read(b); err != nil {
		log.Fatalf("failed to generate token secret: %s", err)
	}
	tokenSecret = []byte(base64.RawURLEncoding.EncodeToString(b))
}

func CreateToken(id int64 ,username string, ttl time.Duration) string {
	exp := time.Now().Add(ttl).Unix()
	payload := fmt.Sprintf("%s.%d.%d", username, id, exp)

	mac := hmac.New(sha256.New, tokenSecret)
	mac.Write([]byte(payload))
	signature := mac.Sum(nil)

	return fmt.Sprintf("%s.%s", base64.RawURLEncoding.EncodeToString([]byte(payload)), base64.RawURLEncoding.EncodeToString(signature))
}

func VerifyToken(token string) (int, string, bool) {
	parts := strings.Split(token, ".")
	tokenLength := 2
	if len(parts) != tokenLength {
		return -1, "", false
	}


	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return -1, "", false
	}

	signature, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return -1, "", false
	}

	mac := hmac.New(sha256.New, tokenSecret)
	mac.Write(payloadBytes)
	expectedSig := mac.Sum(nil)

	if !hmac.Equal(signature, expectedSig) {
		return -1, "", false
	}

	infoLength := 3
	parts2 := strings.Split(string(payloadBytes), ".")
	if len(parts2) != infoLength {
		return -1, "", false
	}

	username := parts2[0]
	id, err := strconv.Atoi(parts2[1])
	if err != nil {
		return -1, "", false
	}
	expUnix, err := strconv.ParseInt(parts2[2], 10, 64)
	if err != nil {
		return -1, "", false
	}

	if time.Now().Unix() > expUnix {
		return -1, "", false
	}

	return id, username, true
}

func UsernameFromContext(ctx context.Context) (string, bool) {
	u, ok := ctx.Value(usernameKey).(string)
	return u, ok
}
