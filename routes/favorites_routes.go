package routes

import (
	"upload-service/controllers"

	"github.com/gorilla/mux"
)

func FavoritesRoutes(router *mux.Router) {
	// router.HandleFunc("/uploadmicro/v1/addContentToFavorites/{UserID}/{ContentID}/{AlbumTitle}", controllers.AddContentToFavorites()).Methods("POST")
    
	router.HandleFunc("/uploadmicro/v1/addContentToFavorites/{UserID}/{ContentID}", controllers.AddContentToFavorites()).Methods("POST")


	// router.HandleFunc("/uploadmicro/v1/removeContentFromFavorites/{UserID}/{ContentID}", controllers.RemoveContentFromFavorites()).Methods("DELETE")

   router.HandleFunc("/uploadmicro/v1/removeContentFromFavorites/{UserID}/{ContentID}/{AlbumID}", controllers.RemoveContentFromFavorites()).Methods("DELETE")


	// router.HandleFunc("/uploadmicro/v1/createNewAlbum/{UserID}/{AlbumTitle}", controllers.CreateNewAlbum()).Methods("POST")

	router.HandleFunc("/uploadmicro/v1/createNewAlbum/{UserID}", controllers.CreateNewAlbum()).Methods("POST")


	
	//router.HandleFunc("/uploadmicro/v1/removeAlbum/{UserID}/{AlbumTitle}", controllers.RemoveAlbum()).Methods("POST")

	router.HandleFunc("/uploadmicro/v1/removeAlbum/{UserID}/{AlbumID}", controllers.RemoveAlbum()).Methods("DELETE")


	//router.HandleFunc("/uploadmicro/v1/moveFavorite/{UserID}/{ContentID}/{FromAblum}/{ToAlbum}", controllers.MoveFavorite()).Methods("POST")

	router.HandleFunc("/uploadmicro/v1/moveFavorite/{UserID}/{ContentID}/{FromAlbumID}/{ToAlbumID}", controllers.MoveFavorite()).Methods("POST")

	router.HandleFunc("/uploadmicro/v1/getUserFavorites/{UserID}", controllers.GetUserFavorites()).Methods("GET")

	router.HandleFunc("/uploadmicro/v1/getUserAlbums/{UserID}", controllers.GetUserAlbums()).Methods("GET")
}
