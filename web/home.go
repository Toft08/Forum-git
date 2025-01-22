package web

import "net/http"

// homePage renders the index page
func homePage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html", nil)
}
