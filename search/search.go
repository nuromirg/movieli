package search

import (
	"Movieli/controllers"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
)

const OMDBURL = "http://www.omdbapi.com/?apikey=aeb7ea37&t="

//var Wd, _ = os.Getwd()

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	/*
	t := template.Must(template.New("index.html").Funcs(template.FuncMap {
		"key": Key,
	}).ParseFiles("public/template/index.html")) */
	t := template.Must(template.New("index.html").ParseFiles("public/template/index.html"))
	err := t.Execute(w, nil)
	if err != nil {
		return
	}
	//search := r.URL.Query().Get("k")
	//fmt.Printf("Get: %s\n", search)
	searchKey(r)
}

func searchKey(r *http.Request) controllers.Movie {
	search := r.URL.Query().Get("k")
	if search != "" {
		search := url.QueryEscape(search)
		fmt.Printf("Get: %s\n", search)
		fmt.Printf("Get: %s\n", OMDBURL + search)
		response, _ := http.Get(OMDBURL + search)
		jsonByteArray, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("Error: %s", err)
			return controllers.Movie{}
		}
		var listOfMovies controllers.Movie
		err = json.Unmarshal(jsonByteArray, &listOfMovies)
		if err != nil {
			fmt.Printf("Error while unmarshaling: %s", err)
			return controllers.Movie{}
		}
		//w.WriteHeader(http.StatusOK)
		//w.Write()

		fmt.Println(listOfMovies)
		return listOfMovies
	}
	return controllers.Movie{}
}









