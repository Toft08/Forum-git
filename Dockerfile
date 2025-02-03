FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o forum

WORKDIR /app

LABEL project="forum"
LABEL description="Forum website for the golang course"

EXPOSE 8080
CMD [ "/app/forum" ]