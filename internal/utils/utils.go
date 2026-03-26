package utils

import (
	"crypto/rand"
	"log"
	"math/big"
)

// SecureRandomInt returns a cryptographically secure random int in [0, limit).
func SecureRandomInt(limit int64) int {
	n, err := rand.Int(rand.Reader, big.NewInt(limit))
	if err != nil {
		log.Fatalf("[ERROR] Failed to generate secure random int: %v", err)
	}
	return int(n.Int64())
}

// GetRemainingDoor returns the door that is neither choice nor montyChoice.
func GetRemainingDoor(choice, montyChoice, numDoors int) int {
	for i := 0; i < numDoors; i++ {
		if i != choice && i != montyChoice {
			return i
		}
	}
	return choice
}
