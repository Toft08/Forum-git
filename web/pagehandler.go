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
	case "/logout":
		Logout(w, r)
	case "/create-post":
		CreatePost(w, r)
	case "/errorhandler":
		errorHandler(w, "error", "error", http.StatusNotFound)
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

func errorHandler(w http.ResponseWriter, errorMessage string, tmpl string, statusCode int) {
	log.Printf("Response status: %d\n", statusCode)
	t, err := template.ParseFiles("templates/"+tmpl+".html", "templates/addons.html")
	if err != nil {
		http.Error(w, errorMessage, http.StatusInternalServerError)
		log.Println("Error parsing template:", err)
		return
	}
	t.Execute(w, t)
	// htmlFileAddress := "templates/" + htmlFileName
	// tmpl, err := template.ParseFiles(htmlFileAddress, "templates/header.html", "templates/navBar.html")
	// if err != nil {
	// 	log.Println(err)
	// 	http.Error(w, errorMessage, statusCode)
	// 	return
	// }
	// err = renderTemplate(w, "error", nil)
	// if err != nil {
	// 	if statusCode == http.StatusInternalServerError {
	// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	}
	// }
}
