package auth

import (
	"context"
	"fmt"

	"templtodo3/database"
	"templtodo3/sessions"
)

func LoginUser(c context.Context, user *database.User) error {
	sessions.PutEncryptedString(c, "userID", user.ID)
	return nil
}

func LogoutUser(c context.Context) {
	sessions.SessionManager.Remove(c, "userID")
}

func GetUserID(c context.Context) string {
	return sessions.GetEncryptedString(c, "userID")
}

func GetUser(c context.Context) (*database.User, error) {
	userID := GetUserID(c)
	if userID == "" {
		return nil, fmt.Errorf("user not logged in")
	}
	return database.GetUserByID(userID)
}
