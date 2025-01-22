package web

import (
	"database/sql"
	"forum/database"
	"html/template"
	"log"
	"net/http"
)

var db *sql.DB

var tmpl = template.Must(template.ParseGlob("../templates/*.html"))

func PageHandler(w http.ResponseWriter, r *http.Request) {

	db = database.InitDB()
	defer db.Close()

	switch r.URL.Path {
	case "/":
		homePage(w, r)
	case "/login":
		login(w, r)
	case "/signup":
		signUp(w, r)

	}
}

// renderTemplate handles the rendering of HTML templates with provided data
func renderTemplate(w http.ResponseWriter, t string, data interface{}) {

	err := tmpl.ExecuteTemplate(w, t, data)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
