package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v4"
	"github.com/sellnat77/Decisionator-api/internal/models"
	"github.com/sellnat77/Decisionator-api/internal/util"
)

var jwtKey = []byte("my_secret_key")

func SignUp(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	if req.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var userCreds models.UserCredentials

	err := json.NewDecoder(req.Body).Decode(&userCreds)
	if err != nil {
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		return
	}
	log.Printf("Body: %s, %s", userCreds.Username, userCreds.Password)

	loginValid := checkPassword(w, &userCreds)
	if !loginValid {
		fmt.Fprint(w, "Invalid login credentials")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokenString, expiration := createJWTToken(w, &userCreds)
	if tokenString == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     `token`,
		Value:    *tokenString,
		Expires:  expiration,
		Path:     "/",
		Secure:   false,
		HttpOnly: false,
		SameSite: http.SameSiteNoneMode,
	})
	fmt.Fprintf(w, "Header token is %s \nExpires %s", *tokenString, expiration)
}

func checkPassword(w http.ResponseWriter, credentials *models.UserCredentials) bool {
	log.Printf("User did not exist %s, creating them", credentials.Username)
	var passwordHash *string
	conn := util.GetConnection()
	dbTrx, err := conn.Begin(context.Background())
	if err != nil {
		log.Printf("Error beginning db transaction %s", err)
	}

	err = dbTrx.QueryRow(context.Background(), `
	SELECT password
	FROM users
 	WHERE username = $1
	LIMIT 1`, credentials.Username).Scan(&passwordHash)

	if err != nil {
		if err == pgx.ErrNoRows {
			log.Print("No rows from existence check")
		} else {
			log.Printf("Rolling back everything %s", err)
			dbTrx.Rollback(context.Background())
		}
	}
	dbTrx.Commit(context.Background())

	ok := false
	if passwordHash != nil {
		log.Printf("Password found is %s", *passwordHash)
		ok = util.CheckPasswordHash(credentials.Password, *passwordHash)
	} else {
		var err error
		credentials.Password, err = util.GeneratehashPassword(credentials.Password)
		log.Printf("Password generated is %s", credentials.Password)
		if err != nil {
			panic("Error generating hash")
		}
		ok = true
	}
	defer conn.Release()
	return ok
}

func createJWTToken(w http.ResponseWriter, credentials *models.UserCredentials) (*string, time.Time) {
	// Expire token in 5 minutes
	expirationTime := time.Now().Add(5 * time.Hour)

	claims := &models.Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, expirationTime
	}
	createUser(credentials)
	return &tokenString, expirationTime
}

func createUser(credentials *models.UserCredentials) int {
	pool := util.GetConnectionPool()
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Error acquiring db connection %s", err)
	}

	dbTrx, err := conn.Begin(context.Background())
	if err != nil {
		log.Printf("Error beginning db transaction %s", err)
	}
	exists := validateUserExists(dbTrx, credentials)
	if !exists {
		log.Printf("User did not exist %s, creating them", credentials.Username)
		_, err = dbTrx.Exec(context.Background(), `
			INSERT INTO users (
				username,
				password,
				email,
				created
				)
			VALUES (
				$1,
				$2,
				$3,
				$4
				) ON CONFLICT DO NOTHING;`, credentials.Username, credentials.Password, credentials.Email, time.Now())
		if err != nil {
			log.Printf("Rolling back everything %s", err)
			dbTrx.Rollback(context.Background())
		}
	}
	dbTrx.Commit(context.Background())
	defer conn.Release()
	return 1
}

func validateUserExists(dbTrx pgx.Tx, credentials *models.UserCredentials) bool {
	var id *int
	conn := util.GetConnection()
	dbTrx, err := conn.Begin(context.Background())
	if err != nil {
		log.Printf("Error beginning db transaction %s", err)
	}

	err = dbTrx.QueryRow(context.Background(), `
	SELECT id
	FROM users
 	WHERE username = $1
	LIMIT 1`, credentials.Username).Scan(&id)

	if err != nil {
		if err == pgx.ErrNoRows {
			log.Print("No rows from existence check")
		} else {
			log.Printf("Rolling back everything %s", err)
			dbTrx.Rollback(context.Background())
		}
	}
	dbTrx.Commit(context.Background())
	if id == nil {
		log.Printf("Id fetched was nil")
		return false
	}
	log.Printf("id fetched is %d", *id)
	return true
}
