package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	router.HandleFunc("/uploadmicro/v1/feedback", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Feedback route works! Method:", r.Method)
	}).Methods("GET", "POST")

	fmt.Println("Test server starting on port 3006...")
	fmt.Println("Test endpoints:")
	fmt.Println("  curl http://localhost:3006/healthz")
	fmt.Println("  curl http://localhost:3006/uploadmicro/v1/feedback")

	http.ListenAndServe(":3006", router)
}
