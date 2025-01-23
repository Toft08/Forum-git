package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func main() {
	db = initDB()
	defer db.Close()

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/", homePage)
	http.HandleFunc("/signup", signUp)
	http.HandleFunc("/login", login)
	http.HandleFunc("/create-post", createPost)
	http.HandleFunc("/logout", logout)

	log.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// renderTemplate handles the rendering of HTML templates with provided data
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles("templates/" + tmpl + ".html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

// homePage renders the index page
func homePage(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Query().Get("message")

	// Fetch posts from the database
	rows, err := db.Query("SELECT title, content FROM Post ORDER BY created_at DESC")
	if err != nil {
		log.Println("Error fetching posts:", err)
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Build a list of posts
	var posts []map[string]string
	for rows.Next() {
		var title, content string
		rows.Scan(&title, &content)
		posts = append(posts, map[string]string{
			"Title":   title,
			"Content": content,
		})
	}

	// Pass posts to template
	renderTemplate(w, "index", map[string]interface{}{
		"Message": message,
		"Posts":   posts,
	})
}

// signUp handles both GET and POST requests for user registration
func signUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "signup", nil)
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Hash the password before storing in database
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Attempt to insert new user into database
		_, err = db.Exec("INSERT INTO User (username, email, password, created_at) VALUES (?, ?, ?, ?)", username, email, hashedPassword, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			log.Println("Error inserting user:", err)
			http.Error(w, "Email already exists", http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// login handles both GET and POST requests for user authentication
func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "login", nil)
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Query database for user's hashed password using their email
		var hashedPassword string
		var userID int
		err := db.QueryRow("SELECT id, password FROM User WHERE username = ?", username).Scan(&userID, &hashedPassword)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// Verify submitted password matches stored hash
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			log.Println("Password does not match hash", err)
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		sessionID := uuid.NewString()

		_, err = db.Exec("INSERT INTO Session (id, user_id, created_at) VALUES (?,?,?)",
			sessionID, userID, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			log.Println("Error creating session:", err)
			http.Error(w, "Failed to create session", http.StatusInternalServerError)
			return
		}
		// Set session ID as a cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
		})
		renderTemplate(w, "index", map[string]interface{}{
			"IsLoggedIn": isLoggedIn,
		})

		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
}

func createPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Render the "Create Post" form
		renderTemplate(w, "create-post", nil)
		return
	}

	if r.Method == http.MethodPost {
		// Parse form values
		title := r.FormValue("title")
		content := r.FormValue("content")
		// Get session ID from cookie
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "Unauthorized: No session found", http.StatusUnauthorized)
			return
		}

		// Retrieve user ID from session
		var userID int
		err = db.QueryRow("SELECT user_id FROM Session WHERE id = ?", cookie.Value).Scan(&userID)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
			return
		}

		// Insert post into the database
		_, err = db.Exec("INSERT INTO Post (title, content, user_id, created_at) VALUES (?, ?, ?, ?)",
			title, content, userID, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			log.Println("Error creating post:", err)
			http.Error(w, "Failed to create post", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}
func isLoggedIn(r *http.Request) bool {
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		return false
	}

	return true
}
