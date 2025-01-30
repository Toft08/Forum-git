package web

import (
	"encoding/json"
	"net/http"
)

func HandleLike(w http.ResponseWriter, r *http.Request) {
	// Check if user is logged in
	isLoggedIn, userID := VerifySession(r)
	if !isLoggedIn {
		http.Error(w, "Must be logged in to like posts", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var data struct {
		PostID int  `json:"post_id"`
		IsLike bool `json:"is_like"` // true for like, false for dislike
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if user already liked/disliked this post
	var existingID int
	err := db.QueryRow("SELECT id FROM Like WHERE user_id = ? AND post_id = ?",
		userID, data.PostID).Scan(&existingID)

	if err == nil {
		// Update existing like/dislike
		_, err = db.Exec("UPDATE Like SET like = ? WHERE id = ?",
			data.IsLike, existingID)
	} else {
		// Create new like/dislike
		_, err = db.Exec("INSERT INTO Like (user_id, post_id, like) VALUES (?, ?, ?)",
			userID, data.PostID, data.IsLike)
	}

	if err != nil {
		http.Error(w, "Failed to save like/dislike", http.StatusInternalServerError)
		return
	}

	// Return updated like counts
	var like, dislike int
	err = db.QueryRow("SELECT COUNT(*) FROM Like WHERE post_id = ? AND like = 1",
		data.PostID).Scan(&like)
	if err != nil {
		http.Error(w, "Failed to get like count", http.StatusInternalServerError)
		return
	}

	err = db.QueryRow("SELECT COUNT(*) FROM Like WHERE post_id = ? AND like = 0",
		data.PostID).Scan(&dislike)
	if err != nil {
		http.Error(w, "Failed to get dislike count", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"like":    like,
		"dislike": dislike,
	})
}
