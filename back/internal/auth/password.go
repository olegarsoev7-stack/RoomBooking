package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword создает хеш из обычного пароля
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword сравнивает обычный пароль с его захешированной версией
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
