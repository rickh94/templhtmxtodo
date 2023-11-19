package routes

import (
	"net/http"
	"templtodo3/components"
	"templtodo3/database"
	"templtodo3/pages"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/mavolin/go-htmx"
)

func TodoList(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*database.User)
	todos, err := database.GetUserTodos(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	csrfToken := csrf.Token(r)
	HxRender(w, r, pages.TodoListPage(todos, csrfToken, ""))
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*database.User)
	r.ParseForm()
	todoText := r.Form.Get("todo")
	todo, err := database.CreateTodo(todoText, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	hxRequest := htmx.Request(r)

	csrfToken := csrf.Token(r)
	var component templ.Component
	if hxRequest == nil {
		todos, err := database.GetUserTodos(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		component = components.Page(pages.TodoListPage(todos, csrfToken, ""))
	} else {
		component = components.TodoListItem(*todo, csrfToken, false)
	}

	w.Header().Set("Content-Type", "text/html")
	component.Render(r.Context(), w)

	return
}

func EditTodo(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*database.User)
	todoID := chi.URLParam(r, "id")
	todo, err := database.GetTodoByID(todoID, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	hxRequest := htmx.Request(r)

	csrfToken := csrf.Token(r)
	var component templ.Component
	if hxRequest == nil {
		todos, err := database.GetUserTodos(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		component = components.Page(pages.TodoListPage(todos, csrfToken, todoID))
	} else {
		component = components.TodoItemForm(*todo, csrfToken)
	}

	w.Header().Set("Content-Type", "text/html")
	htmx.TriggerAfterSettle(r, "FocusInput", FocusInputEvent{
		ID: "todo-input-" + todo.ID,
	})
	component.Render(r.Context(), w)
	return
}

func GetTodo(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*database.User)
	todoID := chi.URLParam(r, "id")
	todo, err := database.GetTodoByID(todoID, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	hxRequest := htmx.Request(r)
	var component templ.Component
	if hxRequest == nil {
		csrfToken := csrf.Token(r)
		todos, err := database.GetUserTodos(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		component = components.Page(pages.TodoListPage(todos, csrfToken, ""))
	} else {
		component = components.TodoItemText(*todo)
	}

	w.Header().Set("Content-Type", "text/html")

	component.Render(r.Context(), w)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*database.User)
	todoID := chi.URLParam(r, "id")
	r.ParseForm()
	todoText := r.Form.Get("text")
	todo, err := database.UpdateTodoText(todoID, todoText, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	hxRequest := htmx.Request(r)
	var component templ.Component
	if hxRequest == nil {
		csrfToken := csrf.Token(r)
		todos, err := database.GetUserTodos(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		component = components.Page(pages.TodoListPage(todos, csrfToken, ""))
	} else {
		component = components.TodoItemText(*todo)
	}

	w.Header().Set("Content-Type", "text/html")

	htmx.TriggerAfterSettle(r, "ShowAlert", ShowAlertEvent{
		Message:  "Successfully updated your todo",
		Title:    "Saved!",
		Variant:  "success",
		Duration: 3000,
	})
	component.Render(r.Context(), w)
	return
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	todoID := chi.URLParam(r, "id")
	user := r.Context().Value("user").(*database.User)
	err := database.DeleteTodo(todoID, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func CompleteTodo(w http.ResponseWriter, r *http.Request) {
	todoID := chi.URLParam(r, "id")
	user := r.Context().Value("user").(*database.User)
	todo, err := database.CompleteTodo(todoID, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	hxRequest := htmx.Request(r)

	csrfToken := csrf.Token(r)
	var component templ.Component
	if hxRequest == nil {
		todos, err := database.GetUserTodos(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		component = components.Page(pages.TodoListPage(todos, csrfToken, ""))
	} else {
		component = components.TodoCompletedButton(*todo, csrfToken)
	}

	w.Header().Set("Content-Type", "text/html")
	component.Render(r.Context(), w)
	return
}

func UnCompleteTodo(w http.ResponseWriter, r *http.Request) {
	todoID := chi.URLParam(r, "id")
	user := r.Context().Value("user").(*database.User)
	todo, err := database.UnCompleteTodo(todoID, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	hxRequest := htmx.Request(r)

	csrfToken := csrf.Token(r)
	var component templ.Component
	if hxRequest == nil {
		todos, err := database.GetUserTodos(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		component = components.Page(pages.TodoListPage(todos, csrfToken, ""))
	} else {
		component = components.TodoCompletedButton(*todo, csrfToken)
	}

	w.Header().Set("Content-Type", "text/html")
	component.Render(r.Context(), w)
	return
}
