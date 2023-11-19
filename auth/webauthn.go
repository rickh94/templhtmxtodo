package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"templtodo3/database"
	"templtodo3/sessions"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

func BeginPasskeyRegistration(c context.Context, user *database.User) (*protocol.CredentialCreation, error) {
	options, session, err := WebAuthnManager.BeginRegistration(user)
	if err != nil {
		return nil, err
	}
	// serialize the session to json and save it with the session manager
	sessionJson, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}
	sessions.PutEncryptedString(c, "webauthnRegistration", string(sessionJson))

	return options, nil
}

func FinishPasskeyRegistration(r *http.Request, user *database.User) error {
	sessionJson := sessions.GetEncryptedString(r.Context(), "webauthnRegistration")
	var session webauthn.SessionData
	err := json.Unmarshal([]byte(sessionJson), &session)
	if err != nil {
		return err
	}

	credential, err := WebAuthnManager.FinishRegistration(user, session, r)
	if err != nil {
		return err
	}
	return user.AddCredential(credential)
}

func BeginPasskeyLogin(r *http.Request, user *database.User) (*protocol.CredentialAssertion, error) {
	options, session, err := WebAuthnManager.BeginLogin(user)
	if err != nil {
		return nil, err
	}

	sessionJson, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}
	sessions.PutEncryptedString(r.Context(), "webauthnLogin", string(sessionJson))

	return options, nil
}

func FinishPasskeyLogin(r *http.Request, user *database.User) error {
	sessionJson := sessions.GetEncryptedString(r.Context(), "webauthnLogin")
	var session webauthn.SessionData
	err := json.Unmarshal([]byte(sessionJson), &session)
	if err != nil {
		log.Default().Printf("failed to unmarshal session: %v\n", err)
		return err
	}

	_, err = WebAuthnManager.FinishLogin(user, session, r)
	if err != nil {
		log.Default().Printf("failed to login: %v\n", err)
		return err
	}

	return nil
}
