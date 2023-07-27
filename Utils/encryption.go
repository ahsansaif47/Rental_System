package utils

import "golang.org/x/crypto/bcrypt"

func Encrypt(password string) (string, error) {
	hBytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hBytes), err
}

func Compare_Encryption(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
