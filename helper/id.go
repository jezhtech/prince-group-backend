package helper

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const base = int64(len(charset))

// Epoch base: Jan 1, 2025
var epochBase = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC).Unix()

func toBaseN(num int64, length int) string {
	result := make([]byte, length)
	for i := length - 1; i >= 0; i-- {
		result[i] = charset[num%base]
		num /= base
	}
	return string(result)
}

func GenerateUserID() string {
	now := time.Now()
	formattedDate := now.Format("060102") // YYMMDD

	// Add seconds since midnight to increase uniqueness
	secondsSinceMidnight := now.Hour()*3600 + now.Minute()*60 + now.Second()
	timePart := fmt.Sprintf("%04d", secondsSinceMidnight%10000) // Ensure 4-digit

	// Generate 2 random alphanumeric characters (36^2 = 1296 variations)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomPart := make([]byte, 2)
	for i := range randomPart {
		randomPart[i] = charset[r.Intn(len(charset))]
	}

	return fmt.Sprintf("PG-%s%s-%s", formattedDate, timePart, string(randomPart))
}

func CheckUserID(userID string) bool {
	return len(userID) == 13 && strings.HasPrefix(userID, "PG-")
}

func GenerateBookingNumber() string {
	// Current time since epochBase in seconds
	seconds := time.Now().Unix() - epochBase

	// Add randomness in the lower bits
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomBits := int64(r.Intn(1296)) // 36^2 possibilities (2 base36 digits)

	// Final number combines time and randomness
	combined := seconds*1296 + randomBits // max 36^6 = 2,176,782,336 total space

	return toBaseN(combined, 6)
}
