package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"zymm/internal/repository"
)

func UserHandlerDelegate(mux *http.ServeMux) {
	mux.HandleFunc("/testHandler1", _TestHandler)
	mux.HandleFunc("/testHandler2", _TestHandler2)
}

func _TestHandler(w http.ResponseWriter, r *http.Request) {

	response := map[string]string{"status": "Test 1"}
	w.Header().Set("Content-Type", "application/json")
	user, userErr := repository.GetAllUsers()

	if userErr != nil {
		fmt.Errorf("ERROR FOUND")
	}
	if len(user) <= 0 {
		fmt.Println("ZERO NADA")
	}
	for i := 0; i < len(user); i++ {
		log.Println(user[i].Name + " - " + user[i].Pass)
	}
	json.NewEncoder(w).Encode(response)
}

func _TestHandler2(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"status": "Test 2"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
