package helper

import (
	"math/rand"
	"strings"
	"time"
)

func GenerateUserID() string {
	now := time.Now()
	formattedDate := now.Format("060102") // YYMMDD format

	// Generate 4 random alphanumeric characters
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	randomPart := make([]byte, 4)
	for i := range randomPart {
		randomPart[i] = charset[r.Intn(len(charset))]
	}

	return "PG-" + formattedDate + "-" + string(randomPart)
}

func CheckUserID(userID string) bool {
	return len(userID) == 13 && strings.HasPrefix(userID, "PG-")
}
