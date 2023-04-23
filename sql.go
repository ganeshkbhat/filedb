package sql

// package sql

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"go/ast"
	"go/parser"
	"go/token"
)

// https://chat.openai.com/c/68356039-c336-4610-bb98-16b4f2a89bcb

/*
Struct for reading and working with CSV Imports
*/

type csvData struct {
	header []string
	rows   [][]string
}

/*
err := execQuery("/path/to/database.db", "CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT, email TEXT)")
if err != nil {
	// handle error
}
*/

/**
err := execQuery("/path/to/database.db", "ALTER TABLE users ADD COLUMN phone TEXT")
if err != nil {
	// handle error
}
*/

/**
err := execQuery("/path/to/database.db", "DROP TABLE IF EXISTS users")
if err != nil {
	// handle error
}
*/

func ExecQuery(dbPath string, query string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

/*
err := dumpTable("/path/to/database.db", "users")
if err != nil {
	// handle error
}
*/

func DumpTable(dbPath, tableName string) error {
	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Prepare SELECT statement
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
	if err != nil {
		return err
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// Prepare slice for row values
	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}

	// Loop through rows
	for rows.Next() {
		// Scan row into slice of pointers to interface{}
		if err := rows.Scan(values...); err != nil {
			return err
		}

		// Print row values
		for i, column := range columns {
			value := *(values[i].(*interface{}))
			fmt.Printf("%s: %v\n", column, value)
		}
		fmt.Println()
	}

	// Check for errors after loop
	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

/*
err := importSQLDump("/path/to/database.db", "/path/to/dump.sql")
if err != nil {
	// handle error
}
*/

func ImportSQLDump(dbPath, dumpPath string) error {
	// Build sqlite3 command
	cmd := exec.Command("sqlite3", dbPath)

	// Redirect input from dump file
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	defer in.Close()
	dumpData, err := os.ReadFile(dumpPath)
	if err != nil {
		return err
	}
	go func() {
		defer in.Close()
		fmt.Fprint(in, string(dumpData))
	}()

	// Run command and wait for it to finish
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("import failed: %v\n%s", err, out)
	}

	return nil
}

func ReadCSV(filename string) (*csvData, error) {
	// Open the CSV file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Parse the CSV data
	reader := csv.NewReader(file)
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Return the CSV data
	return &csvData{header, rows}, nil
}

func ValidateCSV(csvData *csvData) error {
	if len(csvData.header) == 0 {
		// Check if the CSV has a header row
		return errors.New("CSV file is missing header row")
	}

	// Check if all rows have the same number of columns as the header row
	for i, row := range csvData.rows {
		if len(row) != len(csvData.header) {
			return fmt.Errorf("CSV file row %d has incorrect number of columns", i+1)
		}
	}

	// Validation passed
	return nil
}

func ImportCSV(filename, tableName string, db *sql.DB) error {
	// Read the CSV data
	csvData, err := ReadCSV(filename)
	if err != nil {
		return err
	}

	// Validate the CSV structure
	err = ValidateCSV(csvData)
	if err != nil {
		return err
	}

	// Prepare the INSERT statement
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		StrJoin(csvData.header, ", "),
		StrJoin(MakePlaceholders(len(csvData.header)), ", "),
	)
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Insert the CSV data
	for _, row := range csvData.rows {
		_, err := stmt.Exec(row)
		if err != nil {
			return err
		}
	}

	// Import successful
	return nil
}

func StrJoin(strs []string, sep string) string {
	return "\"" + sep + "\""
}

func MakePlaceholders(count int) []string {
	placeholders := make([]string, count)
	for i := 0; i < count; i++ {
		placeholders[i] = "?"
	}
	return placeholders
}

func ExecTransactionQuery(dbPath string, query string, args ...interface{}) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)
	if err != nil {
		return err
	}

	return nil
}

/**
func main() {
    // Open the database connection
    db, err := sql.Open("sqlite3", "mydatabase.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Define the query and arguments for bulk update
    query := "INSERT INTO mytable (name, value) VALUES (?, ?)"

    args := [][]interface{}{
        {"apple", 1, 1},
        {"banana", 2, 2},
        {"cherry", 3, 3},
    }

	err := execBulkQuery("example.db", query, data)
    if err != nil {
        log.Fatal(err)
    }
}
*/

/**
func main() {
    // Open the database connection
    db, err := sql.Open("sqlite3", "mydatabase.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    // Define the query and arguments for bulk update
    query := "UPDATE mytable SET name = ?, value = ? WHERE id = ?"
    args := [][]interface{}{
        {"apple", 1, 1},
        {"banana", 2, 2},
        {"cherry", 3, 3},
    }
    // Execute the bulk update
    err = bulkQuery(db, query, args...)
    if err != nil {
        log.Fatal(err)
    }
}
*/

/**
func main() {
    // Open the database connection
    db, err := sql.Open("sqlite3", "mydatabase.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Define the query and arguments for bulk delete
    query := "DELETE FROM mytable WHERE id = ?"
    args := [][]interface{}{
        {1},
        {2},
        {3},
    }

    // Execute the bulk delete
    err = bulkQuery(db, query, args...)
    if err != nil {
        log.Fatal(err)
    }
}
*/

func ExecBulkQuery(dbPath string, query string, data [][]interface{}, args ...interface{}) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, row := range data {
		_, err := stmt.Exec(append(row, args...)...)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

/**
type MyRow struct {
    ID    int
    Name  string
    Value int
}

func main() {
    dbName := "mydatabase.db"
	// Open the database connection
    db, err := sql.Open("sqlite3", dbName)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Define the query and arguments for bulk update
    query := "UPDATE mytable SET name = ?, value = ? WHERE id = ?"
    args := []interface{}{"new_name", 100}
    data := []MyRow{
        {1, "name1", 10},
        {2, "name2", 20},
        {3, "name3", 30},
    }

    // Execute the bulk update
    err = bulkUpdate(dbName, query, data, args...)
    if err != nil {
        log.Fatal(err)
    }
}
*/

func BulkUpdate(dbName string, query string, data interface{}, args ...interface{}) error {
	// Open the database connection
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return err
	}
	defer db.Close()

	// Begin the transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Prepare the statement
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Iterate over the data and execute the query for each row
	rows := reflect.ValueOf(data)
	for i := 0; i < rows.Len(); i++ {
		row := rows.Index(i)
		_, err := stmt.Exec(append(row.Interface().([]interface{}), args...)...)
		if err != nil {
			return err
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

/**

package main

import (
	"fmt"
)

func main() {
	dbPath := "example.db"
	tableName := "mytable"
	columns := []string{"name", "age", "email"}
	data := [][]interface{}{
		{"John", 25, "john@example.com"},
		{"Jane", 30, "jane@example.com"},
		{"Bob", 40, "bob@example.com"},
	}

	err := bulkInsert("example.db", "users", columns, data)
	if err != nil {
		fmt.Println("Error inserting data:", err)
	} else {
		fmt.Println("Data inserted successfully")
	}
}
*/

func BulkInsert(dbPath, tableName string, columns []string, data [][]interface{}) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	columnsStr := strings.Join(columns, ", ")
	params := "?" + strings.Repeat(", ?", len(columns)-1)
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, columnsStr, params)

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, row := range data {
		_, err := stmt.Exec(row...)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

/**
	package main

	import (
		"fmt"
	)

	func main() {
		dbPath := "./mydb.db"
		tableName := "mytable"
		query := "DELETE FROM " + tableName + " WHERE name = ?"
		data := [][]interface{}{
			{"apple"},
			{"banana"},
			{"cherry"},
		}

		err := bulkDelete(dbPath, query, data)
		if err != nil {
			fmt.Println("Error deleting data:", err)
			return
		}

		fmt.Println("Data deleted successfully!")
	}
*/

func BulkDelete(dbName string, query string, data interface{}, args ...interface{}) error {
	// Open the database connection
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return err
	}
	defer db.Close()

	// Begin the transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Prepare the statement
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Iterate over the data and execute the query for each row
	rows := reflect.ValueOf(data)
	for i := 0; i < rows.Len(); i++ {
		row := rows.Index(i)
		_, err := stmt.Exec(append(row.Interface().([]interface{}), args...)...)
		if err != nil {
			return err
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
