package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(p string) string {
	h, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(h)
}
func CheckPassword(hash, p string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(p)) == nil
}

func GetSha256Hash(str string) string {
	h := sha256.Sum256([]byte(str))
	return hex.EncodeToString(h[:])
}

func GenerateRandomHash() string {
	id := xid.New().String()
	bytes := make([]byte, 16) // 256-bit
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}

	raw := "tx_" + base64.RawURLEncoding.EncodeToString([]byte(id)) + "_" + base64.RawURLEncoding.EncodeToString(bytes)

	hash := GetSha256Hash(raw)

	// prefix := raw[:10]

	return hash
}
