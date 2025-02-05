package main

import (
	"forum/database"
	"forum/web"
	"log"
	"net/http"
)

func main() {

	http.Handle("/assets/", http.FileServer(http.Dir(".")))

	db := database.InitDB()
	defer db.Close()

	database.MakeTables(db)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		web.PageHandler(w, r, db)
	})

	log.Println("Server is running on http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting the server:", err)
	}
}
