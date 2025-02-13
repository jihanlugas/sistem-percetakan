package utils

import (
	"net/mail"
)

// return bool

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsAvailablePreload(data string, preloads []string) bool {
	for _, preload := range preloads {
		if preload == data {
			return true
		}
	}
	return false
}
