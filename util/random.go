package util

import (
	"crypto/rand"
	"fmt"
	"strings"
)

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits  = "0123456789"
	special = "!@#$%^&*()_+-=[]{}|;:,.<>?"
)

func GenerateRandomString(length int, includeLetters, includeDigits, includeSpecial bool) (string, error) {
	var chars string
	if includeLetters {
		chars += letters
	}
	if includeDigits {
		chars += digits
	}
	if includeSpecial {
		chars += special
	}
	if chars == "" {
		return "", fmt.Errorf("no characters specified")
	}

	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	for _, v := range b {
		sb.WriteByte(chars[v%byte(len(chars))])
	}

	return sb.String(), nil
}
