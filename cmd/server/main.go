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

	// key := make([]byte, 64)
	// _, _ = rand.Read(key)
	// LogService.LogMessage("my key: " + base64.StdEncoding.EncodeToString(key))

	// load variables from .env into the environment
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, falling back to system env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("❌ PORT not set in environment. Exiting...")
	}

	mux := http.NewServeMux()

	db.InitDB()

	api.UserHandlerDelegate(mux)
	api.AuthHandlerDelegate(mux)

	log.Println("🚀 Server running on http://localhost:" + port)
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// docker build . -t uditnair90/api-padhai-golang:latest
// docker push uditnair90/api-padhai-golang:latest

/*  DOCKER BUILD AND PUSH
docker build . -t uditnair90/api-padhai-golang:latest && docker push uditnair90/api-padhai-golang:latest
*/

// docker pull uditnair90/api-padhai-golang:latest
// docker run -d --pull=always --quiet --name uditnair90_api-padhai-golang --env PORT=5000 --publish 5000:5000 uditnair90/api-padhai-golang:latest
// docker run -v ~/secrets/.env:/app/.env -d --pull=always --quiet --name uditnair90_api-padhai-golang --env PORT=5000 --publish 5000:5000 uditnair90/api-padhai-golang:latest
// docker image prune -f     ###REMOVES UNUSED IMAGES###
