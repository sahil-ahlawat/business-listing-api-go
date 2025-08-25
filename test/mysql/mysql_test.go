// File: fitness/test/mysql/mysql_test.go
package myysql_test

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestMySQLConnection(t *testing.T) {
	dsn := "root:@tcp(localhost:3306)/mysql"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("MySQL connection open failed: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Fatalf("MySQL ping failed: %v", err)
	}

	// Optional test query
	var version string
	err = db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		t.Fatalf("MySQL query failed: %v", err)
	}

	fmt.Printf("Connected to MySQL version: %s\n", version)
}
