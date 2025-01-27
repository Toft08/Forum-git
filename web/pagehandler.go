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

	data := PageDetails{}

	db = database.InitDB()
	defer db.Close()

	switch r.URL.Path {
	case "/":
		HomePage(w, r, &data)
	case "/login":
		Login(w, r, &data)
	case "/signup":
		SignUp(w, r, &data)
	case "/logout":
		Logout(w, r, &data)
	case "/create-post":
		CreatePost(w, r, &data)

	default:
		if strings.HasPrefix(r.URL.Path, "/post") {
			PostHandler(w, r, &data)
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
	log.Printf("Error message: %s\n", errorMessage)

	w.WriteHeader(statusCode)

	err := tmpl.ExecuteTemplate(w, "error.html", map[string]string{
		"ErrorMessage": errorMessage,
	})
	if err != nil {
		// If rendering the template fails
		log.Printf("Error executing template: %s\n", err)
		http.Error(w, errorMessage, statusCode)
		return
	}
}
