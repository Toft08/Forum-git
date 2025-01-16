# Forum Project Checklist

## 1. Setup the Environment
- [ ] Initialize a new Go project: `go mod init forum`
- [ ] Set up **Docker** with a `Dockerfile` and `docker-compose.yml` for SQLite.
- [ ] Create a folder structure:
- [ ] Add dependencies: `sqlite3`, `bcrypt`, and `UUID`.

---

## 2. Database Design
- **Plan an Entity-Relationship Diagram (ERD)** for these entities:
- Users: `id`, `email`, `username`, `password`
- Posts: `id`, `title`, `content`, `author_id`, `created_at`
- Comments: `id`, `content`, `post_id`, `author_id`, `created_at`
- Categories: `id`, `name`
- Post_Categories: `post_id`, `category_id`
- Likes: `id`, `user_id`, `post_or_comment_id`, `type`
- Write SQL queries:
- [ ] `CREATE TABLE` for the above.
- [ ] Use `INSERT` for test data.
- [ ] `SELECT` to retrieve posts, comments, and likes.

---

## 3. User Authentication
- [ ] **Register**:
- Form: Email, Username, Password (hashed with `bcrypt`).
- Check for duplicate email.
- Store data in SQLite.
- [ ] **Login**:
- Check email and password.
- Create a cookie-based session using UUID.
- Handle cookie expiration.
- [ ] **Middleware**:
- Restrict access to certain routes for logged-in users only.

---

## 4. Forum Features
- **Posts**:
- [ ] Create posts (logged-in users only) with categories.
- [ ] Display posts to all users.
- [ ] Allow filtering by categories.
- **Comments**:
- [ ] Add comments (logged-in users only).
- [ ] Display comments below posts.
- **Likes/Dislikes**:
- [ ] Registered users can like/dislike posts and comments.
- [ ] Show the count of likes and dislikes to all users.
- **Filters**:
- [ ] Categories (e.g., subforums).
- [ ] Posts created by the logged-in user.
- [ ] Posts liked by the logged-in user.

---

## 5. Frontend
- Use only HTML and CSS:
- Create templates (`templates` folder) for:
  - [ ] Login and Register pages.
  - [ ] Main page (list of posts with filters).
  - [ ] Post details (with comments and likes).
  - [ ] Create/Edit post forms.
- Serve static files (CSS) from the `/static` folder.

---

## 6. Error Handling
- [ ] Return appropriate HTTP status codes (e.g., 404 for not found, 401 for unauthorized).
- [ ] Handle database and server errors gracefully.

---

## 7. Testing
- [ ] Write unit tests for:
- Authentication logic.
- Database queries.
- Middleware.
- [ ] Use Go's built-in `testing` package.

---

## 8. Final Steps
- [ ] Containerize the application:
- Use Docker to run both the web app and SQLite.
- [ ] Add documentation to the README:
- How to run the app.
- Database schema and ERD.
- Features and usage instructions.
