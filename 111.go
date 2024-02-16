package main
import (	
	"encoding/json"
	"fmt"	"log"
	"net/http")
func main() {
	http.HandleFunc("/", handleGet)	log.Fatal(http.ListenAndServe(":8080", nil))
}
func handleGet(w http.ResponseWriter, r *http.Request) {	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)		return
	}
	response := map[string]interface{}{
		"message": "OK",	}
	w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)		return
	}}