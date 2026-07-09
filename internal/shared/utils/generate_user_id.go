package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

func GenerateUserId(username string) (string, error) {
	cleanUsername := strings.ToLower(strings.ReplaceAll(username, " ", ""))

	max := big.NewInt(10000)
	randomNumber, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", errors.New("failed to generate user id")
	}

	return fmt.Sprintf("usr_%s_%04d", cleanUsername, randomNumber.Int64()), nil
}