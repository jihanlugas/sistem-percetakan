package utils

import (
	"github.com/google/uuid"
)

func GetUniqueID() string {
	uuid := uuid.New().String()
	return uuid
}
