package main

import (
	"log"
	"net/http"
	"zymm/internal/api"
	"zymm/internal/config"
	"zymm/internal/db"
)

/*
	go run ./cmd/server/main.go

*/

func main() {

	config.Init()

	mux := http.NewServeMux()

	db.InitDB()

	api.UserHandlerDelegate(mux)
	api.AuthHandlerDelegate(mux)
	api.MembershipHandlerDelegate(mux)
	api.AttendanceApiPrefixHandlerDelegate(mux)
	api.FeedbackHandlerDelegate(mux)
	api.NotificationHandlerDelegate(mux)
	api.GymHandlerDelegate(mux)
	api.EmployeeHandlerDelegate(mux)
	api.MessagesHandlerDelegate(mux)

	log.Println("🚀 Server running on http://localhost:" + config.GetAppPortInString())
	err := http.ListenAndServe(":"+config.GetAppPortInString(), mux)
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
	// clear the exited only
	docker rm $(docker ps -a -q --filter status=exited --filter ancestor=uditnair90/zymm-backend:latest)

	docker rm $(docker ps -a -q --filter ancestor=uditnair90/zymm-backend:latest)
*/

/*
 	docker run -v ~/secrets/.env:/app/.env -d --pull=always --quiet --name uditnair90_zymm-backend --publish 5000:5000 uditnair90/zymm-backend:latest

	// for debug ( not daemon )
 	docker run -v ~/secrets/.env:/app/.env --pull=always --name uditnair90_zymm-backend --publish 5000:5000 uditnair90/zymm-backend:latest
*/
