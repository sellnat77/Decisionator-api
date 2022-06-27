package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sellnat77/Decisionator-api/internal/handlers"

	"github.com/joho/godotenv"
)

func main() {
	currEnv, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := getEnvVar("PORT", currEnv)
	portString := fmt.Sprintf(":%s", port)
	http.HandleFunc("/helloworld", handlers.HelloWorld)
	http.HandleFunc("/signup", handlers.SignUp)
	http.HandleFunc("/signin", handlers.SignIn)

	fmt.Printf("Starting server on %s", portString)
	http.ListenAndServe(portString, nil)
}

func getEnvVar(key string, currEnv map[string]string) string {
	envVar := os.Getenv(key)
	if len(envVar) == 0 {
		return currEnv[key]
	}
	return envVar
}
