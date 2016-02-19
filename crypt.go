package main

import (
	"crypto/sha256"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

// HashPassword PBKDF2 for encrypting passwords
func HashPassword(password string) string {
	salt := []byte{0x11, 0x22, 0x33, 0x44}
	hashed := pbkdf2.Key([]byte(password), salt, 4096, sha256.Size, sha256.New)
	return fmt.Sprintf("%x", hashed)
}
