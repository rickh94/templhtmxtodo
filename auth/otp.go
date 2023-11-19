package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"templtodo3/sessions"
	"time"
)

func generateRandomCode() (string, error) {
	const codeLength = 6
	const charset = "0123456789"

	randomCode := make([]byte, codeLength)
	charsetLength := big.NewInt(int64(len(charset)))

	for i := 0; i < codeLength; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		randomCode[i] = charset[randomIndex.Int64()]
	}

	return string(randomCode), nil
}

func GenerateOTP(c context.Context) (string, error) {
	code, err := generateRandomCode()
	if err != nil {
		log.Printf("failed to generate code: %v\n", err)
		return "", fmt.Errorf("failed to generate code")
	}
	SaveOTP(c, code)

	return code, nil
}

func SaveOTP(c context.Context, code string) {
	sessions.PutEncryptedString(c, "code", code)
	sessions.SessionManager.Put(c, "codeCreated", time.Now().Unix())
}

func CheckOTP(c context.Context, submittedCode string) bool {
	expectedCode := sessions.GetEncryptedString(c, "code")
	created := sessions.SessionManager.GetInt64(c, "codeCreated")
	createdTime := time.Unix(created, 0)

	if time.Since(createdTime) > 5*time.Minute {
		return false
	}
	if submittedCode == expectedCode {
		sessions.SessionManager.Remove(c, "code")
		sessions.SessionManager.Remove(c, "codeCreated")
		return true
	}
	return false
}
