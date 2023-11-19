package main

import (
	"fmt"
	"net/http"
	"time"

	"templtodo3/auth"
	"templtodo3/config"
	"templtodo3/database"
	"templtodo3/routes"
	"templtodo3/sender"
	"templtodo3/sessions"
	"templtodo3/static"

	"github.com/benbjohnson/hashfs"
	"github.com/go-chi/chi/v5"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/mavolin/go-htmx"
)

//go:generate templ generate
//go:generate bun run build
func main() {

	// SETUP
	config.Init() // Always call this first
	database.Init()
	sessions.Init()
	auth.Init()
	sender.Init()

	r := chi.NewRouter()

	// MIDDLEWARE
	r.Use(middleware.Logger)
	r.Use(htmx.NewMiddleware())
	r.Use(csrf.Protect([]byte(config.AppConfig.SecretKey), csrf.Secure(true)))
	r.Use(middleware.Compress(5, "application/json", "text/html", "text/css", "application/javascript"))
	// Just long enough for preload to matter
	r.Use(middleware.SetHeader("Cache-Control", "max-age=300"))
	r.Use(middleware.SetHeader("Vary", "HX-Request"))
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(sessions.SessionManager.LoadAndSave)

	// the static router is kinda funny, works better to give it its own servemux then mount that
	s := http.NewServeMux()
	s.Handle("/", http.StripPrefix("/static", hashfs.FileServer(static.HashStatic)))
	r.Mount("/static", s)

	// ROUTES
	r.Get("/", routes.Index)

	// AUTH ROUTES
	r.Route("/auth", func(r chi.Router) {
		r.Get("/login", routes.StartLogin)
		r.Post("/login", routes.ContinueLogin)
		r.Post("/code", routes.CompleteCodeLogin)
		r.Get("/code", routes.ForceCodeLogin)
		r.Get("/logout", routes.LogoutUser)
		r.Post("/passkey/register", routes.RegisterPasskey)
		r.Post("/passkey/login", routes.CompletePasskeySignin)
		r.With(routes.LoginRequired).Get("/me", routes.UserInfo)
	})

	r.With(routes.LoginRequired).Route("/todos", func(r chi.Router) {
		r.Get("/", routes.TodoList)
		r.Post("/", routes.CreateTodo)
		r.Get("/{id}", routes.GetTodo)
		r.Delete("/{id}", routes.DeleteTodo)
		r.Get("/{id}/edit", routes.EditTodo)
		r.Post("/{id}/edit", routes.UpdateTodo)
		r.Post("/{id}/complete", routes.CompleteTodo)
		r.Post("/{id}/uncomplete", routes.UnCompleteTodo)
	})

	fmt.Println("Listening on :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
	}
}
