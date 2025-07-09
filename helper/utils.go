package helper

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

func GenerateUserID() string {
	now := time.Now()
	formattedDate := now.Format("20060102")
	return "PG-" + formattedDate + "-" + uuid.New().String()
}

func CheckUserID(userID string) bool {
	return len(userID) == 23 && strings.HasPrefix(userID, "PG-")
}
