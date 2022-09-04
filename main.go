package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sellnat77/Decisionator-api/internal/handlers"
	"github.com/sellnat77/Decisionator-api/internal/models"
	"github.com/sellnat77/Decisionator-api/internal/util"
)

func main() {

	port := util.GetEnvVar("PORT")
	portString := fmt.Sprintf(":%s", port)
	initDatastores()

	http.HandleFunc("/helloworld", handlers.HelloWorld)
	http.HandleFunc("/signup", handlers.SignUp)
	http.HandleFunc("/signin", handlers.SignUp)
	http.HandleFunc("/meet", handlers.Meet)
	http.HandleFunc("/health", handlers.Healthcheck)
	http.HandleFunc("/", handlers.Preflight)

	log.Printf("Starting server on %s", portString)
	http.ListenAndServe(portString, nil)
}

func initDatastores() {
	dbHost := util.GetEnvVar("DB_HOST")
	dbPort := util.GetEnvVar("DB_PORT")
	dbUser := util.GetEnvVar("DB_USER")
	dbPass := util.GetEnvVar("DB_PASSS")
	dbName := util.GetEnvVar("DB_NAME")
	creds := models.DatastoreCredentials{Host: dbHost, Port: dbPort, User: dbUser, Password: dbPass, DBName: dbName}
	util.Initialize(creds)
	success := util.CreateUsersTable()
	if !success {
		log.Fatal("Users table not created")
	}
}
