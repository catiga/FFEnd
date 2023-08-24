package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	password := "123456"
	hashPassByte := sha256.Sum256([]byte(password))
	hashPass := hex.EncodeToString(hashPassByte[:])

	fmt.Println(hashPass)
}
