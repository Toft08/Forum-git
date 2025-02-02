package main

import (
	"forum/database"
	"forum/web"
	"log"
	"net/http"
)

func main() {

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	db := database.InitDB()
	defer db.Close()

	database.MakeTables(db)

	http.HandleFunc("/", web.PageHandler)

	log.Println("Server is running on http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting the server:", err)
	}
}
