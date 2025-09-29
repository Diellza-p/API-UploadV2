package main

import (
	"fmt"
	"log"
	"net/http"
	"upload-service/routes"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// Add debug route first
	router.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Debug endpoint working!")
	}).Methods("GET")

	// Add the favorites routes
	routes.FavoritesRoutes(router)

	// Add health check
	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	// Add a route to list all routes (for debugging)
	router.HandleFunc("/routes", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Available routes:")
		fmt.Fprintln(w, "GET /debug")
		fmt.Fprintln(w, "GET /healthz")
		fmt.Fprintln(w, "GET /uploadmicro/v1/getUserFavorites/{UserID}")
		fmt.Fprintln(w, "POST /uploadmicro/v1/addContentToFavorites/{UserID}/{ContentID}/{AlbumTitle}")
		fmt.Fprintln(w, "DELETE /uploadmicro/v1/removeContentFromFavorites/{UserID}/{ContentID}")
		fmt.Fprintln(w, "POST /uploadmicro/v1/createNewAlbum/{UserID}/{AlbumTitle}")
		fmt.Fprintln(w, "POST /uploadmicro/v1/removeAlbum/{UserID}/{AlbumTitle}")
		fmt.Fprintln(w, "POST /uploadmicro/v1/moveFavorite/{UserID}/{ContentID}/{FromAblum}/{ToAlbum}")
	}).Methods("GET")

	port := "3006"

	fmt.Printf("Debug server listening on port %s\n", port)
	fmt.Println("Test endpoints:")
	fmt.Println("  curl http://localhost:3006/debug")
	fmt.Println("  curl http://localhost:3006/healthz")
	fmt.Println("  curl http://localhost:3006/routes")
	fmt.Println("  curl http://localhost:3006/uploadmicro/v1/getUserFavorites/68d92bf74688878879ca9af1")

	log.Fatal(http.ListenAndServe(":"+port, router))
}
