#!/bin/bash

# Build the Docker image
docker build -t forum-app .

# Run the Docker container
docker run -d --name forum-container -p 8080:8080 forum-app
