# syntax=docker/dockerfile:1
FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /decisionator-api

EXPOSE $PORT
EXPOSE $JWT_KEY
EXPOSE $DB_HOST
EXPOSE $DB_PORT
EXPOSE $DB_USER
EXPOSE $DB_PASSS
EXPOSE $DB_NAME

CMD ["/decisionator-api"]
