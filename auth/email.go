package auth

import (
	"context"
	"templtodo3/sessions"
)

func SaveEmail(c context.Context, email string) {
	sessions.PutEncryptedString(c, "email", email)
}

func GetEmail(c context.Context) string {
	return sessions.GetEncryptedString(c, "email")
}
