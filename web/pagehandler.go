package web

import (
	"database/sql"
	"forum/database"
	"html/template"
	"log"
	"net/http"
)

var db *sql.DB

var tmpl = template.Must(template.ParseGlob("templates/*.html"))

func PageHandler(w http.ResponseWriter, r *http.Request) {

	db = database.InitDB()
	defer db.Close()

	switch r.URL.Path {
	case "/":
		HomePage(w, r)
	case "/login":
		Login(w, r)
	case "/signup":
		SignUp(w, r)

	}
}

// renderTemplate handles the rendering of HTML templates with provided data
func renderTemplate(w http.ResponseWriter, t string, data interface{}) {

	err := tmpl.ExecuteTemplate(w, t+".html", data)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
