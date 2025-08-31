package room

import (
	"crypto/rand"
	"encoding/hex"
)

const (
	digitNr = 6
)

func GenerateRoomID() (string, error) {
	b := make([]byte, digitNr) 
	_ , err := rand.Read(b)
	if err != nil {
		return "" , err
	}
	return hex.EncodeToString(b), nil
}
