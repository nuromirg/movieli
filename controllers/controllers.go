package controllers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Resp struct {
	Message interface{}
	Error string
}

type Movie struct {
	Id string `json:"id"`
	Poster string `json:"poster"` //or it should be url.URL ..?
	Title string `json:"title"`
	Year string	`json:"year"`
	Director string `json:"director"`
}

type MovieStorage struct {
	movies []Movie
}

func HelloHandler(w http.ResponseWriter, request *http.Request) {

	name := strings.Replace(request.URL.Path, "/hello/", "", 1)

	response := Resp{
		Message: fmt.Sprintf("Hello, %s! Glad to see you.", name),
	}
	responseJson, _ := json.Marshal(response) //convert to json structure
	w.WriteHeader(http.StatusOK)

	w.Write(responseJson)
}

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		log.Printf("server [net/http] method [%s] connection from [%v]", request.Method, request.RemoteAddr)
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, request)
	}
}
//middleware func
func BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	//login:pass
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(auth) != 2 || auth[0] != "Basic" { //see documentation
			http.Error(w, "Authorization failed", http.StatusUnauthorized)
			return
		}
		hashed, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(hashed), ":", 2)
		log.Printf("pair %+v", pair)
		if len(pair) != 2 || !bauth(pair[0], pair[1]) {
			http.Error(w, "Authorization failed", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func bauth(username, password string) bool {
	if username == "test" && password == "test" {
		return true
	}
	return false
}

func MovieHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		handlerGetMovie(w, r)
	} else if r.Method == http.MethodPost {
		handlerAddMovie(w, r)
	} else if r.Method == http.MethodDelete {
		handlerDeleteMovie(w, r)
	} else if r.Method == http.MethodPut {
		handlerUpdateMovie(w, r)
	}
}

func MoviesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		handlerGetMovie(w, r)
	}
	w.WriteHeader(http.StatusOK)
	response := Resp{
		Message: movieStorage.GetMovies(),
	}
	moviesJson, _ := json.Marshal(response)
	w.Write(moviesJson)
}

func handlerGetMovie(w http.ResponseWriter, r *http.Request){
	var response Resp
	id := strings.Replace(r.URL.Path, "/movie/", "", 1)
	movie := movieStorage.FindMovieById(id)
	if movie == nil {
		w.WriteHeader(http.StatusNotFound)
		response.Error = fmt.Sprintf("")
		responseJson, _ := json.Marshal(response)
		w.Write(responseJson)
	}
	response.Message = movie
	responseJson, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
	return
}

func handlerAddMovie(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	var movie Movie
	var response Resp
	err := decoder.Decode(&movie)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Error = err.Error()
		responseJson, _ := json.Marshal(response)
		w.Write(responseJson)
		return
	}
	err = movieStorage.AddMovies(movie)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Error = err.Error()
		responseJson, _ := json.Marshal(response)
		w.Write(responseJson)
		return
	}

	MoviesHandler(w, r) //result
}

func handlerDeleteMovie(w http.ResponseWriter, r *http.Request) {
	id := strings.Replace(r.URL.Path, "/movie/", "", 1)
	var response Resp
	err := movieStorage.DeleteMovie(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Error = err.Error()
		responseJson, _ := json.Marshal(response)
		w.Write(responseJson)
		return
	}
	MoviesHandler(w, r)
}

func handlerUpdateMovie(w http.ResponseWriter, r *http.Request) {
	id := strings.Replace(r.URL.Path, "/movie/", "", 1)
	decoder := json.NewDecoder(r.Body)
	var movie Movie
	var response Resp
	err := decoder.Decode(&movie)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Error = err.Error()
		responseJson, _ := json.Marshal(response)
		w.Write(responseJson)
		return
	}
	movie.Id = id
	err = movieStorage.UpdateMovie(movie)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response.Error = fmt.Sprintf("")
		responseJson, _ := json.Marshal(response)
		w.Write(responseJson)
		return
	}
	response.Message = movie
	responseJson, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}



var movieStorage = MovieStorage{
	movies: make([]Movie, 0),
}

func (s MovieStorage) FindMovieById(id string) *Movie {
	for _, movie := range s.movies {
		if movie.Id == id {
			return &movie
		}
	}
	return nil
}

func (s MovieStorage) GetMovies() []Movie {
	return s.movies
}
func (s *MovieStorage) AddMovies(movie Movie) error{
	for _, mv:= range s.movies {
		if mv.Id == movie.Id {
			return errors.New(fmt.Sprintf("Movie with id %s not found.", movie.Id))
		}
	}
	s.movies = append(s.movies, movie)
	return nil
}

func (s *MovieStorage) UpdateMovie(movie Movie) error {
	for i, mv := range s.movies {
		if mv.Id == movie.Id {
			s.movies[i] = movie
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Movie with id %s not found.", movie.Id))
}

func (s *MovieStorage) DeleteMovie(id string) error {
	for i, mv := range s.movies {
		if mv.Id == id {
			s.movies = append(s.movies[:i], s.movies[i+1:]...)
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Movie with id %s not found (was deleted).", id))
}
