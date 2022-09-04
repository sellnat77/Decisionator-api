package util

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sellnat77/Decisionator-api/internal/models"
)

var dbConnection *pgxpool.Pool

func Initialize(creds models.DatastoreCredentials) {
	var once sync.Once
	once.Do(func() {
		fmt.Println(creds.ConnectionStringPGX())
		dbConn, err := pgxpool.Connect(context.Background(), creds.ConnectionStringPGX())

		if err != nil {
			panic(err)
		}
		dbConnection = dbConn
	})
}

func GetConnectionPool() *pgxpool.Pool {
	return dbConnection
}

func GetConnection() *pgxpool.Conn {
	pool := GetConnectionPool()
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Error acquiring db connection %s", err)
	}
	return conn
}

func Healthcheck() bool {
	err := GetConnectionPool().Ping(context.Background())
	if err != nil {
		fmt.Printf("DB Healthcheck failed %s", err.Error())
		return false
	}
	fmt.Println("DB Healthcheck passed")
	return true
}

func QueryOne(field string, table string, filterField string, filter string, result *string) {
	pool := GetConnectionPool()
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Error acquiring db connection %s", err)
	}
	dbTrx, err := conn.Begin(context.Background())
	if err != nil {
		log.Printf("Error beginning db transaction %s", err)
	}

	query := fmt.Sprintf("SELECT %s FROM %s", field, table)
	log.Printf(query)

	err = dbTrx.QueryRow(context.Background(), query+`
 	WHERE $1 = $2
	LIMIT 1`, filterField, filter).Scan(&result)

	if err != nil {
		if err == pgx.ErrNoRows {
			log.Print("No rows from existence check")
		} else {
			log.Printf("Rolling back everything %s", err)
			dbTrx.Rollback(context.Background())
		}
	}
	dbTrx.Commit(context.Background())
	defer conn.Release()
}

func QueryOneInt(field string, table string, filterField string, filter string, result *int) {
	pool := GetConnectionPool()
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Error acquiring db connection %s", err)
	}
	dbTrx, err := conn.Begin(context.Background())
	if err != nil {
		log.Printf("Error beginning db transaction %s", err)
	}

	query := fmt.Sprintf("SELECT %s FROM %s", field, table)
	log.Printf(query)

	err = dbTrx.QueryRow(context.Background(), query+`
 	WHERE $1 = $2
	LIMIT 1`, filterField, filter).Scan(&result)

	if err != nil {
		if err == pgx.ErrNoRows {
			log.Print("No rows from existence check")
		} else {
			log.Printf("Rolling back everything %s", err)
			dbTrx.Rollback(context.Background())
		}
	}
	dbTrx.Commit(context.Background())
	defer conn.Release()
}

func QueryMany() {

}

func CreateUsersTable() bool {
	ok := true

	pool := GetConnectionPool()
	dbTrx, err := pool.Acquire(context.Background())
	if err != nil {
		fmt.Printf("Error acquiring db connection %s", err)
		ok = false
	}
	_, err = dbTrx.Exec(context.Background(), `
		DROP TABLE IF EXISTS users;
		CREATE TABLE IF NOT EXISTS users (
				ID serial PRIMARY KEY,
				username VARCHAR (254) UNIQUE NOT NULL,
				email VARCHAR (254) UNIQUE NOT NULL,
				password VARCHAR (60) NOT NULL,
				created TIMESTAMP NOT NULL
		);
		`)
	if err != nil {
		log.Fatal(err)
		fmt.Printf("Error creating table users %s", err)
		ok = false
	}

	defer dbTrx.Release()

	return ok
}
