package helpers

import (
	"golang.org/x/crypto/bcrypt"
)

// hash password digunakan untuk menyamarkan password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// membandingkan password hash dengan password hash user yang disimpan
func ComparePassword(password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}
