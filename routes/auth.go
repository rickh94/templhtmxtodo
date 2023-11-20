package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"templtodo3/auth"
	"templtodo3/database"
	"templtodo3/pages"
	"templtodo3/sender"

	"github.com/gorilla/csrf"
	"github.com/mavolin/go-htmx"
)

func StartLogin(w http.ResponseWriter, r *http.Request) {
	if RedirectLoggedIn(w, r, "/auth/me") {
		return
	}
	csrfToken := csrf.Token(r)
	nextLoc := r.URL.Query().Get("next")
	if nextLoc == "" {
		nextLoc = "/auth/me"
	}
	HxRender(w, r, pages.StartLoginPage(csrfToken, nextLoc))
}

func ContinueLogin(w http.ResponseWriter, r *http.Request) {
	if RedirectLoggedIn(w, r, "/auth/me") {
		return
	}
	r.ParseForm()

	userEmail := r.Form.Get("email")
	if userEmail == "" {
		w.WriteHeader(http.StatusBadRequest)
		// TODO: send back the form with an error
		return
	}

	nextLoc := r.Form.Get("next")
	if nextLoc == "" {
		nextLoc = "/nextLoc"
	}
	user, err := database.GetUserByEmail(userEmail)
	if err != nil || len(user.Credentials) == 0 {
		ContinueOtpSignIn(w, r, userEmail, nextLoc)
		return
	} else {
		ContinuePasskeySignin(w, r, user, nextLoc)
		return
	}

}

func ForceCodeLogin(w http.ResponseWriter, r *http.Request) {
	userEmail := auth.GetEmail(r.Context())
	if userEmail == "" {
		w.WriteHeader(http.StatusBadRequest)
		// TODO: send back the form with an error
		return
	}
	nextLoc := r.Form.Get("next")
	if nextLoc == "" {
		nextLoc = "/todos"
	}
	ContinueOtpSignIn(w, r, userEmail, nextLoc)
	return
}

func ContinueOtpSignIn(w http.ResponseWriter, r *http.Request, userEmail string, nextLoc string) {
	token := csrf.Token(r)

	code, err := auth.GenerateOTP(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		// TODO: send back the form with an error
		return
	}
	message := fmt.Sprintf("Your one-time login code is %s. It will expire in 5 minutes.", code)
	go sender.SendEmail(userEmail, "TemplTodo: Login Code", message)
	auth.SaveEmail(r.Context(), userEmail)
	htmx.TriggerAfterSettle(r, "ShowAlert", ShowAlertEvent{
		Message:  "Login Code sent to your email",
		Title:    "Code Sent!",
		Variant:  "success",
		Duration: 3000,
	})
	HxRender(w, r, pages.FinishCodeLoginPage(token, nextLoc))
	return
}

func ContinuePasskeySignin(w http.ResponseWriter, r *http.Request, user *database.User, nextLoc string) {
	token := csrf.Token(r)
	auth.SaveEmail(r.Context(), user.Email)
	options, err := auth.BeginPasskeyLogin(r, user)
	if err != nil {
		ContinueOtpSignIn(w, r, user.Email, nextLoc)
		return
	}
	HxRender(w, r, pages.FinishPasskeyLoginPage(options, token, nextLoc))
	return
}

func CompleteCodeLogin(w http.ResponseWriter, r *http.Request) {
	if RedirectLoggedIn(w, r, "/") {
		return
	}
	userEmail := auth.GetEmail(r.Context())
	if userEmail == "" {
		w.WriteHeader(http.StatusBadRequest)
		// TODO: redirect back to the main form with an error
		return
	}

	r.ParseForm()
	nextLoc := r.Form.Get("next")
	if nextLoc == "" {
		nextLoc = "/todos"
	}
	submittedCode := r.Form.Get("code")
	if auth.CheckOTP(r.Context(), submittedCode) {
		user, err := database.GetOrCreateUser(userEmail)
		if err != nil || user == nil {
			w.WriteHeader(http.StatusInternalServerError)
			// TODO: re-render the form with an error
			return
		}
		auth.LoginUser(r.Context(), user)
		Redirect(w, r, nextLoc)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		// TODO: re-render the form with an error
		return
	}

}

func CompletePasskeySignin(w http.ResponseWriter, r *http.Request) {
	userEmail := auth.GetEmail(r.Context())
	user, err := database.GetUserByEmail(userEmail)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		// TODO: re-render the form with an error
		return
	}

	if err := auth.FinishPasskeyLogin(r, user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Failed to log in"})
		return
	} else {
		auth.LoginUser(r.Context(), user)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "redirect": "/auth/me"})
		return
	}
}

func UserInfo(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetUser(r.Context())
	if err != nil {
		fmt.Println("Could not find user")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	credentialCount, err := database.CountUserCredentials(user)
	if err != nil {
		fmt.Println("Could not count credentials")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "text/html")

	registrationOptions, err := auth.BeginPasskeyRegistration(r.Context(), user)
	if err != nil {
		fmt.Println("Could not find registration options")
	}

	token := csrf.Token(r)
	component := pages.UserInfoPage(user, registrationOptions, token, fmt.Sprintf("%d", credentialCount))
	HxRender(w, r, component)
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	auth.LogoutUser(r.Context())
	Redirect(w, r, "/")
}

func RegisterPasskey(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetUser(r.Context())
	if err != nil {
		fmt.Println("Could not find user")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err = auth.FinishPasskeyRegistration(r, user); err != nil {
		log.Default().Printf("Error registering passkey for user %s: %v", user.Email, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not register passkey"))
		return
	}
	w.Write([]byte("OK"))
	w.WriteHeader(http.StatusCreated)
}

func DeletePasskeys(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetUser(r.Context())
	if err != nil {
		fmt.Println("Could not find user")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err = database.DeletePasskeys(user); err != nil {
		log.Default().Printf("Error deleting passkeys for user %s: %v", user.Email, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not delete passkeys"))
		return
	}
	htmx.TriggerAfterSettle(r, "ShowAlert", ShowAlertEvent{
		Message:  "All your passkeys have been deleted. Consider registering a new one!",
		Title:    "Passkeys Deleted!",
		Variant:  "success",
		Duration: 3000,
	})
	w.Write([]byte("OK"))
	w.WriteHeader(http.StatusOK)
}
