package web

import (
	"log"
	"net/http"
	"time"
)

// Logout logs out the user by deleting the session from the database and setting the session cookie to expire
func Logout(w http.ResponseWriter, r *http.Request, data *PageDetails) {

	cookie, err := r.Cookie("session_id")
	if err == nil {
		// Delete session from database
		_, err := db.Exec("DELETE FROM Session WHERE id = ?", cookie.Value)
		if err != nil {
			log.Println("Error deleting session:", err)
		}
	}
	// Expire the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	data.LoggedIn = false

	http.Redirect(w, r, "/", http.StatusFound)
}
