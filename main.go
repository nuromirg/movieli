package main

import (
	"Movieli/config"
	"Movieli/controllers"
	"Movieli/search"
	"Movieli/service"
	"log"
	"net/http"
	"time"
)

func main() {
	service.InitDB(config.DBNAME)
	handler := http.NewServeMux()
	//handler.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("/public/"))))
	handler.HandleFunc("/", search.IndexHandler)
	handler.HandleFunc("/search/", controllers.Logger(search.MovieSearchHandler))
	handler.HandleFunc("/movie/", controllers.Logger(controllers.MovieHandler))
	handler.HandleFunc("/movie/delete", controllers.Logger(controllers.MovieDeleteHandler))
	handler.HandleFunc("/movies/", controllers.Logger(controllers.MoviesHandler))


	server := http.Server{
		Addr: config.SERVERADDR,
		Handler: handler,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 2^20 128kByte
	}

	log.Printf("Listening on http://%s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}