package web

import (
	"database/sql"
	"forum/database"
	"html/template"
	"log"
	"net/http"
	"strings"
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
	case "/logout":
		Logout(w, r)
	case "/create-post":
		CreatePost(w, r)

	default:
		if strings.HasPrefix(r.URL.Path, "/post") {
			PostHandler(w, r)
		} else {
			ErrorHandler(w, "Page not found", "error", http.StatusNotFound)
		}
	}
}

// renderTemplate handles the rendering of HTML templates with provided data
func RenderTemplate(w http.ResponseWriter, t string, data interface{}) {

	err := tmpl.ExecuteTemplate(w, t+".html", data)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func ErrorHandler(w http.ResponseWriter, errorMessage string, tmpl string, statusCode int) {
	log.Printf("Response status: %d\n", statusCode)
	w.WriteHeader(statusCode) // Set the HTTP status code before rendering
	RenderTemplate(w, tmpl, map[string]string{
		"ErrorMessage": errorMessage,
	})
}
