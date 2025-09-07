package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
)

var DB *sql.DB

func InitDB() {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	if user == "" || pass == "" || host == "" || port == "" || name == "" {
		log.Fatal("❌ Database environment variables not set")
	}

	connString := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
		user, pass, host, port, name)

	var err error
	DB, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatalf("❌ Error opening DB: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("❌ Could not connect to DB: %v", err)
	}

	log.Println("✅ Connected to MSSQL database")
}
