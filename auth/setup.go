package auth

import (
	"fmt"
	"templtodo3/config"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

var WebAuthnManager *webauthn.WebAuthn

func Init() {
	wconfig := &webauthn.Config{
		RPDisplayName: config.AppConfig.DisplayName,
		RPID:          config.AppConfig.Hostname,
		RPOrigins:     []string{fmt.Sprintf("https://%s", config.AppConfig.Hostname)},
		AuthenticatorSelection: protocol.AuthenticatorSelection{
			RequireResidentKey: protocol.ResidentKeyNotRequired(),
			UserVerification:   protocol.VerificationRequired,
		},
		Timeouts: webauthn.TimeoutsConfig{
			Login: webauthn.TimeoutConfig{
				Enforce:    true,                 // Require the response from the client comes before the end of the timeout.
				Timeout:    time.Second * 60 * 5, // Standard timeout for login sessions.
				TimeoutUVD: time.Second * 60 * 5, // Timeout for login sessions which have user verification set to discouraged.
			},
			Registration: webauthn.TimeoutConfig{
				Enforce:    true,                  // Require the response from the client comes before the end of the timeout.
				Timeout:    time.Second * 60 * 10, // Standard timeout for registration sessions.
				TimeoutUVD: time.Second * 60 * 10, // Timeout for login sessions which have user verification set to discouraged.
			},
		},
	}

	var err error
	WebAuthnManager, err = webauthn.New(wconfig)
	if err != nil {
		panic(err)
	}
}
