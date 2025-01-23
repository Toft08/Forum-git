package web

import (
	"log"
	"net/http"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}
func IsLoggedIn(r *http.Request) (bool, int) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Println("No session ID cookie found")
		return false, 0
	}
	log.Println("Session ID:", cookie.Value)
	
	var userID int
	err = db.QueryRow("SELECT user_id FROM Session WHERE id = ?", cookie.Value).Scan(&userID)
	if err != nil {
		return false, 0
	}

	return true, userID
}
