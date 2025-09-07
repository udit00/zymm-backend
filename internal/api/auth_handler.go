package api

import (
	"encoding/json"
	"net/http"
)

func AuthHandlerDelegate(mux *http.ServeMux) {
	mux.HandleFunc("/authHandler1", _AuthHandler)
	mux.HandleFunc("/authHandler2", _AuthHandler2)
}

func _AuthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"status": "Test 1"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func _AuthHandler2(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"status": "Test 2"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
