package web

import (
	"net/http"
)

// homePage renders the index page
func HomePage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", nil)
}
