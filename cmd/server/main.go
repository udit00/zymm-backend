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

	port := os.Getenv("ZYMM_PORT")
	if port == "" {
		log.Fatal("❌ ZYMM_PORT not set in environment. Exiting...")
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

/*  
	docker build . -t uditnair90/zymm-backend:latest
	docker push uditnair90/zymm-backend:latest
	docker pull uditnair90/zymm-backend:latest
	docker build . -t uditnair90/zymm-backend:latest && docker push uditnair90/zymm-backend:latest

	docker run -v ~/secrets/.env:/app/.env -d --pull=always --quiet --name uditnair90_zymm-backend --publish 5000:5000 uditnair90/zymm-backend:latest

	docker image prune -f     ###REMOVES UNUSED IMAGES###
	docker rm $(docker ps -a -q --filter status=exited --filter ancestor=uditnair90/zymm-backend:latest)
*/

/*
 	docker run -v ~/secrets/.env:/app/.env -d --pull=always --quiet --name uditnair90_zymm-backend --publish 5000:5000 uditnair90/zymm-backend:latest

	// for debug ( not daemon )
 	docker run -v ~/secrets/.env:/app/.env --pull=always --name uditnair90_zymm-backend --publish 5000:5000 uditnair90/zymm-backend:latest
*/
