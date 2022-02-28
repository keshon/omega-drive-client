package utils

import (
	"log"
	"os"
	"time"
)

func ExpiryTimer(expiryDate string) {
	// Setup a start and end time
	createdAt := time.Now()
	expiresAt, _ := time.Parse(time.RFC822, expiryDate+" 00:00 UTC")

	// Get delta diff
	diff := expiresAt.Sub(createdAt)

	// Exit
	if diff.Hours() <= 0 {
		log.Println("[utils/expiry_timer.go][ExpiryTimer] Timer expired. Latest date is " + expiryDate)
		os.Exit(1)
	}
}
