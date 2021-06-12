package main

import (
	"Movieli/controllers"
	"Movieli/search"
	"log"
	"net/http"
	"time"
)




func main() {
	handler := http.NewServeMux() // Маршрутизация
	//handler.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("/public/"))))
	handler.HandleFunc("/", search.IndexHandler)
	handler.HandleFunc("/hello/", controllers.Logger(controllers.BasicAuth(controllers.HelloHandler)))
	handler.HandleFunc("/movie/", controllers.Logger(controllers.MovieHandler))
	handler.HandleFunc("/movies/", controllers.Logger(controllers.MoviesHandler))


	server := http.Server{
		Addr: "localhost:8080",
		Handler: handler,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 2^20 128kByte
	}

	log.Printf("Listening on http://%s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}