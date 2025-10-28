package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"zymm/internal/config"

	_ "github.com/denisenkom/go-mssqldb"
)

var DB *sql.DB

func InitDB() {

	envWiseHost := "DB_HOST"
	envWisePass := "DB_PASS"

	if config.IsProduction() {
		envWiseHost = "LIVE_DB_HOST"
		envWisePass = "LIVE_DB_PASS"
	}

	user := "sa"
	pass := os.Getenv(envWisePass)
	host := os.Getenv(envWiseHost)
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	if user == "" || pass == "" || host == "" || port == "" || name == "" {
		log.Fatal("❌ Database environment variables not set")
	}

	connString := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s&encrypt=true&trustservercertificate=true",
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
