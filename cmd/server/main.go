package main

import (
	"log"
	"net/http"
	"os"
	"zymm/internal/api"
	"zymm/internal/db"

	"github.com/joho/godotenv"
)

/*
	go run ./cmd/server/main.go

*/

func main() {
	// load variables from .env into the environment
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, falling back to system env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("❌ PORT not set in environment. Exiting...")
	}

	mux := http.NewServeMux()

	/*
	 register routes
	 http://localhost:10000/<handler>
	 ex: http://localhost:10000/api/
	*/

	db.InitDB()

	api.UserHandlerDelegate(mux)
	api.AuthHandlerDelegate(mux)

	log.Println("🚀 Server running on http://localhost:" + port)
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
