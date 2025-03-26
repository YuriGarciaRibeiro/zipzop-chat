package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandomSalt() string {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		panic(err) // Em produção, trate este erro adequadamente
	}
	return hex.EncodeToString(salt)
}
