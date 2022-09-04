package util

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var currEnv *map[string]string

func getCurrentEnvironment() *map[string]string {
	if currEnv == nil {
		readEnv, err := godotenv.Read()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		currEnv = &readEnv
	}

	return currEnv
}

func GetEnvVar(key string) string {
	envVar := os.Getenv(key)
	if len(envVar) == 0 {
		tmpVar := *getCurrentEnvironment()
		return tmpVar[key]
	}
	return envVar
}
