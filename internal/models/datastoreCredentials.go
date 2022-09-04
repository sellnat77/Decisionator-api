package models

import "fmt"

type DatastoreCredentials struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func (c DatastoreCredentials) ConnectionStringPGX() string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", c.User, c.Password, c.Host, c.Port, c.DBName)
}
