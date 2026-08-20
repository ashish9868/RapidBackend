package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"math/big"

	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(p string) string {
	h, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(h)
}

func GenerateRandomPassword(length int) string {
	const (
		lower   = "abcdefghijklmnopqrstuvwxyz"
		upper   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		numbers = "0123456789"
		special = "!@#$%^&*-+[]{},."
	)

	if length < 8 {
		length = 8
	}

	all := lower + upper + numbers + special
	password := make([]byte, length)

	sets := []string{lower, upper, numbers, special}

	// Guarantee each character category.
	for i, set := range sets {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(set))))
		password[i] = set[n.Int64()]
	}

	// Fill remaining characters.
	for i := 4; i < length; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(all))))
		password[i] = all[n.Int64()]
	}

	// Shuffle the password.
	for i := len(password) - 1; i > 0; i-- {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		j := int(n.Int64())

		password[i], password[j] = password[j], password[i]
	}

	return string(password)
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
