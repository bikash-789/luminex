package healthcheck

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func NewHandler() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			log.Println("Error writing response", err)
		}
	}).Methods("GET")
	return r
}
