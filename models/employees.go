package models

import (
	"bytes"
	"errors"
	"fmt"
	"log"
)

type Employee struct {
	number    int
	birthDate string
	firstName string
	lastName  string
	gender    string
	hireDate  string
	title     string
}

func fieldMap() func(string) string {
	innerMap := map[string]string{
		"number":    "emp_no",
		"birthDate": "birth_date",
		"firstName": "first_name",
		"lastName":  "last_name",
		"gender":    "gender",
		"hireDate":  "hire_date",
	}

	return func(key string) string {
		return innerMap[key]
	}
}

//func (e Employee) current_salary int {
//    // TODO
//}
func (db *DB) AllEmployees(limit int) ([]*Employee, error) {
	var query string

	if limit > 0 {
		query = fmt.Sprintf("SELECT * FROM employees LIMIT %d", limit)
	} else {
		query = "SELECT * FROM employees"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	employees := make([]*Employee, 0)
	for rows.Next() {
		e := new(Employee)
		err := rows.Scan(&e.number, &e.birthDate, &e.firstName, &e.lastName, &e.gender, &e.hireDate)
		if err != nil {
			return nil, err
		}
		employees = append(employees, e)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return employees, nil
}

func (db *DB) SearchEmployees(conditions map[string]string) ([]*Employee, error) {
	var query_buffer bytes.Buffer
	values := make([]interface{}, 0, len(conditions))

	query_buffer.WriteString("SELECT * FROM employees WHERE ")
	for key, value := range conditions {
		log.Printf("DEBUG: \"%s\" -> \"%s\"", key, value)

		field := fieldMap()(key)
		log.Printf("DEBUG: key \"%s\" maps to field \"%s\"", key, field)
		if field == "" { // throw exception
			return nil, errors.New("Unrecognized field \"" + key + "\"")
		}
		query_buffer.WriteString(field)
		query_buffer.WriteString("=? ") // FIXME this will only work for a single condition. Need to join them
		values = append(values, value)
	}

	log.Printf("DEBUG: Query: \"%s\"", query_buffer.String())

	for idx, value := range values {
		log.Printf("DEBUG: value %d: \"%s\"", idx, value)
	}

	rows, err := db.Query(query_buffer.String(), values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	employees := make([]*Employee, 0)
	for rows.Next() {
		e := new(Employee)
		err := rows.Scan(&e.number, &e.birthDate, &e.firstName, &e.lastName, &e.gender, &e.hireDate)
		if err != nil {
			return nil, err
		}
		employees = append(employees, e)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	log.Printf("DEBUG: returning %d employees", len(employees))
	return employees, nil
}

func (db *DB) FindEmployee(number int) (*Employee, error) {
	e := new(Employee)
	err := db.QueryRow("SELECT * FROM employees WHERE emp_no=?", number).Scan(&e.number, &e.birthDate, &e.firstName, &e.lastName, &e.gender, &e.hireDate)
	if err != nil {
		return nil, err
	}
	return e, nil
}
