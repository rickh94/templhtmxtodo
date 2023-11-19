package routes

import (
	"net/http"
	"templtodo3/pages"
)

func Index(w http.ResponseWriter, r *http.Request) {
	HxRender(w, r, pages.Index())
}
