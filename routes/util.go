package routes

import (
	"context"
	"log"
	"net/http"
	"net/url"

	"templtodo3/auth"
	"templtodo3/components"

	"github.com/a-h/templ"
	"github.com/mavolin/go-htmx"
)

func HxRender(w http.ResponseWriter, r *http.Request, component templ.Component) {
	hxRequest := htmx.Request(r)

	if hxRequest == nil {
		component = components.Page(component)
	}
	w.Header().Set("Content-Type", "text/html")
	component.Render(r.Context(), w)
}

func Redirect(w http.ResponseWriter, r *http.Request, url string) {
	if htmx.Request(r) != nil {
		htmx.PushURL(r, url)
		htmx.Redirect(r, url)
		return
	} else {
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

func LoginRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.GetUser(r.Context())
		if err != nil {
			location := r.URL.Path
			location = url.QueryEscape(location)
			Redirect(w, r, "/auth/login?next="+location)
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RedirectLoggedIn(w http.ResponseWriter, r *http.Request, url string) bool {
	if _, err := auth.GetUser(r.Context()); err == nil {
		Redirect(w, r, "/auth/me")
		return true
	}
	return false
}

func RedirectNotLoggedIn(w http.ResponseWriter, r *http.Request, url string) bool {
	if _, err := auth.GetUser(r.Context()); err != nil {
		log.Println(err)
		Redirect(w, r, "/auth/login")
		return true
	}
	return false
}

type ShowAlertEvent struct {
	Message  string `json:"message"`
	Title    string `json:"title"`
	Variant  string `json:"variant"`
	Duration int    `json:"duration"`
}

type FocusInputEvent struct {
	ID string `json:"id"`
}
