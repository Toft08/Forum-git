-- USER Table
CREATE TABLE User (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TEXT NOT NULL
);

-- POST Table
CREATE TABLE Post (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES User (id) ON DELETE CASCADE
);

-- COMMENT Table
CREATE TABLE Comment (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT NOT NULL,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TEXT NOT NULL,
    FOREIGN KEY (post_id) REFERENCES Post (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES User (id) ON DELETE CASCADE
);

-- CATEGORY Table
CREATE TABLE Category (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);

-- POST_CATEGORY Table
CREATE TABLE Post_Category (
    -- Might need to change the primary key to a composite key like the one below.
    -- PRIMARY KEY (category_id, post_id),
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    category_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    FOREIGN KEY (category_id) REFERENCES Category (id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES Post (id) ON DELETE CASCADE
);

-- LIKE Table
CREATE TABLE Like (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    post_id INTEGER,
    comment_id INTEGER,
    created_at TEXT NOT NULL,
    type INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES User (id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES Post (id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES Comment (id) ON DELETE CASCADE
);

-- Create SESSION table
CREATE TABLE Session (
    id TEXT PRIMARY KEY, -- Unique session ID (UUID)
    user_id INTEGER NOT NULL,
    created_at TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES USER (id)
);

INSERT INTO Category (name) VALUES ('General');
INSERT INTO Category (name) VALUES ('Tutorial');
INSERT INTO Category (name) VALUES ('Question');

INSERT INTO User (email, username, password, created_at) VALUES 
('admin@example.com', 'admin', 'hashedpassword', datetime('now'));

INSERT INTO Post (title, content, user_id, created_at) VALUES 
('Welcome to the forum', 'This is the first post!', 1, datetime('now'));