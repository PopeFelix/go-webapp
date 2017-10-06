package models

import (
	"fmt"
	//	"github.com/fatih/structs"
	"log"
	"os"
	"testing"
)

const path string = "/home/took/go/src/go-webapp/share/db/employees.db" // FIXME: How can I make this relative to the app directory?

var TheDB *DB

func TestMain(m *testing.M) {
	db, err := new_db(path)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	TheDB = db
	os.Exit(m.Run())
}

func TestAllEmployees(t *testing.T) {
	employees, err := TheDB.AllEmployees(10)
	if err != nil {
		t.Error(err)
	}
	if len(employees) != 10 {
		t.Errorf("Expected 10 employee records, got %d", len(employees))
	}
	log.Printf("OK %d employees", len(employees))
	for idx, employee := range employees {
		if employee.number <= 0 {
			t.Errorf("Got invalid employee number for employee %d/%d", idx, len(employees))
		}
		if employee.firstName == "" {
			t.Errorf("Got empty first name for employee %d/%d", idx, len(employees))
		}
		if employee.lastName == "" {
			t.Errorf("Got empty last name for employee %d/%d", idx, len(employees))
		}
		if employee.gender == "" {
			t.Errorf("Got empty gender for employee %d/%d", idx, len(employees))
		}
		if employee.hireDate == "" {
			t.Errorf("Got empty hire date for employee %d/%d", idx, len(employees))
		}
	}
}

func TestSearchEmployees(t *testing.T) {
	t.Run("Employee exists", func(t *testing.T) {
		conditions := map[string]string{"lastName": "Zykh"} // All employees named "Zykh". There are 148.
		employees, err := TheDB.SearchEmployees(conditions)
		if err != nil {
			t.Error(err)
		}
		expectedCount := 148
		if len(employees) != expectedCount {
			t.Errorf("Expected %d employee records, got %d", expectedCount, len(employees))
		}
		for _, employee := range employees {
			if employee.lastName != "Zykh" {
				t.Errorf("Got employee with unexpected last name \"%s\"", employee.lastName)
			}
		}
	})

	t.Run("Employee does not exist", func(t *testing.T) {
		conditions := map[string]string{"lastName": "Gambolputty ... auf Uln"}
		employees, err := TheDB.SearchEmployees(conditions)
		if err != nil {
			t.Error(err)
		}
		if len(employees) != 0 {
			t.Errorf("Expected 0 employee records, got %d", len(employees))
		}
	})

	t.Run("Invalid field", func(t *testing.T) {
		fieldName := "dingleDangle"
		conditions := map[string]string{fieldName: "banana"}
		employees, err := TheDB.SearchEmployees(conditions)
		if err == nil {
			t.Error("No error returned for invalid field")
		}
		s := err.Error()
		expected := fmt.Sprintf("Unrecognized field \"%s\"", fieldName)
		if s != expected {
			t.Errorf("Expected error message \"%s\", got \"%s\"", expected, s)
		}
		if len(employees) != 0 {
			t.Errorf("Expected 0 employee records, got %d", len(employees))
		}
	})

	t.Run("All supported fields", func(t *testing.T) {
		fields := [6]string{"number", "birthDate", "firstName", "lastName", "gender", "hireDate"}

		for _, field := range fields {
			conditions := map[string]string{field: "argle bargle"}
			_, err := TheDB.SearchEmployees(conditions)
			if err != nil {
				t.Errorf("Caught error with field \"%s\"")
				t.Error(err)
			}
		}
	})
}

func TestFindEmployee(t *testing.T) {
	//10001|1953-09-02|Georgi|Facello|M|1986-06-26

	t.Run("Existing employee", func(t *testing.T) {
		empNumber := 10001
		employee, err := TheDB.FindEmployee(empNumber)
		if err != nil {
			t.Error(err)
		}
		expected := Employee{number: empNumber, birthDate: "1953-09-02T00:00:00Z", firstName: "Georgi", lastName: "Facello", gender: "M", hireDate: "1986-06-26T00:00:00Z"}
		if *employee != expected {
			//			for _, field := range structs.Fields(expected) {
			//				expectedVal := field.Value()
			//
			//				expectedVal := structs.Field(field).Value()
			//				if expected.field != employee.field {
			//					t.Errorf("\"%s\" != \"%s\"", employee.field, expected.field)
			//				}
			//			}
			t.Errorf("Employee \"%s\" does not match expected value \"%s\"", *employee, expected) // FIXME: compare values
		}
	})
}
