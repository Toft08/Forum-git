# Forum-git

Github repo for the Forum project at grit:lab

# Usage

First dowload the repo to you local machine with 'git clone'
give premissions to the shell script with 'chmod +x setup.sh'
run the script with './setup.sh'
now database is created and you can start the server with 'go run .'

![alt text](ERD.png)

# Authentication
Once the user logs in, they are given a UUID token in a session cookie. This token is used to authenticate the user. The token is stored in the database and is used to authenticate the user. When an user logs out, the token is deleted from the database. When a user is registering, we store is username and hashed password with bcrypt in the database.

# Frontend

# Backend