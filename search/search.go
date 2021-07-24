package search

import (
	"Movieli/config"
	"Movieli/controllers"
	"Movieli/service"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)
var t *template.Template
const OMDBURL = "http://www.omdbapi.com/?" + config.OMDBAPI

//var Wd, _ = os.Getwd()

var funcMap = template.FuncMap{
	"delete": deleteFromDB,
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	var arrayDB []controllers.Movie
	executeTemplate(w, t, arrayDB)
	APIRequestHandler(w, r)
}

func MovieSearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		APIRequestHandler(w, r)
	}
}

func LoadTemplate() *template.Template {
	return template.Must(template.New("index.html").Funcs(funcMap).ParseFiles("public/template/index.html"))
}

func executeTemplate(w http.ResponseWriter, t *template.Template, arrayDB []controllers.Movie) {
	t = LoadTemplate()
	arrayDB = dbList()
	err := t.Execute(w, arrayDB)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
}

func readerToString(response http.Response) string {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(response.Body)
	if err != nil {
		return ""
	}
	return buf.String()
}

func APIRequestHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("k")
	if search != "" {
		search = url.QueryEscape(search)
		//fmt.Printf("Get: %s\n", search)
		//fmt.Printf("Get: %s\n", OMDBURL + search)
		response, _ := http.Get(OMDBURL + search)

		jsonByteArray, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("Error (responseBuffer): %s", err)
		}
		//fmt.Printf("jsonByteArray: \n %s\n", jsonByteArray)
		jsonStringMap := make(map[string]string)
		json.Unmarshal(jsonByteArray, &jsonStringMap)
		if !dbTitleIteranceCheck(jsonStringMap) {
			return
		}

		client := &http.Client{}
		req, err := http.NewRequest("POST", "http://" + config.SERVERADDR + "/movie/" , bytes.NewBuffer(jsonByteArray))
		if err != nil {
			fmt.Println(err)
		}

		req.Header.Add("Accept", "text/html")
		req.Header.Add("User-Agent", "MSIE/15.0")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		io.Copy(os.Stdout, resp.Body)
	}
}

func dbList() []controllers.Movie {
	db := service.DBConnect()

	defer db.Close()
	rows, err := db.Query("SELECT * FROM Movieli")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	var movies []controllers.Movie

	for rows.Next() {
		var mr controllers.Movie
		err := rows.Scan(&mr.Id, &mr.Poster, &mr.Title, &mr.Year, &mr.Director)
		if err != nil{
			fmt.Println(err)
			continue
		}
		movies = append(movies, mr)
	}
	return movies
}

func dbTitleIteranceCheck(jsonStringMap map[string]string) bool {
	moviesDB := dbList()
	//fmt.Printf("movies DB: \n %s\n", moviesDB)
	if jsonStringMap["Title"] == "" { return false }
	for i := 0; i < len(moviesDB); i++ {
		if jsonStringMap["Title"] == moviesDB[i].Title && jsonStringMap["Year"] == moviesDB[i].Year {
			return false
		}
	}
	return true
}

func deleteFromDB(idFromTemplate string) bool {
	client := &http.Client{}
	idQuery := strings.NewReader(idFromTemplate)
	req, err := http.NewRequest("DELETE", "http://" + config.SERVERADDR + "/movie/", idQuery)
	if err != nil {
		fmt.Println(err)
		return false
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
	return true
}





