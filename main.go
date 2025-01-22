package main

import (
	"forum/web"
	"log"
	"net/http"
)

func main() {

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/", web.PageHandler)

	log.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
