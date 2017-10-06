// based on http://www.alexedwards.net/blog/organising-database-access
package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Datastore interface {
	AllEmployees(int) ([]*Employee, error)
	//    AllDepartments() ([]*Department, error) # TODO: Department model when I've figured out Employee
}

type DB struct {
	*sql.DB
}

func new_db(path string) (*DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil // returns a pointer to the "db" member of the "DB" struct
}
