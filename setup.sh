#!/bin/bash
# Install SQLite if needed
if ! command -v sqlite3 &> /dev/null; then
    echo "SQLite not found, installing..."
    sudo apt install sqlite3
fi

# Create the database from migration script
sqlite3 database.db < init.sql
echo "Database initialized!"
