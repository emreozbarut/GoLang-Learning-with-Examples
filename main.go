package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	hosts      = "localhost:27017"
	database   = "monGo"
	username   = ""
	password   = ""
	collection = "articles"
)

func initialiseMongo() (session *mgo.Session){

	info := &mgo.DialInfo{
		Addrs:    []string{hosts},
		Timeout:  60 * time.Second,
		Database: database,
		Username: username,
		Password: password,
	}

	session, err := mgo.DialWithInfo(info)
	if err != nil {
		panic(err)
	}

	return
}

type MongoSession struct {
	session *mgo.Session
}

type Article struct {
	Title string
	Desc string
	Content string
}

var mongoSession = MongoSession{}

func homePage(responseWriter http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		http.ServeFile(responseWriter, request, "./templates/form.html")

	case "POST":
		if err := request.ParseForm(); err != nil {
			fmt.Fprintf(responseWriter, "ParseForm() err: %v", err)
			return
		}
		fmt.Fprintf(responseWriter, "Post from website! r.PostFrom = %v\n", request.PostForm)

		article := Article{
			request.FormValue("title"),
			request.FormValue("description"),
			request.FormValue("content")}

		fmt.Fprintf(responseWriter, "Title: %s\n", article.Title)
		fmt.Fprintf(responseWriter, "Description: %s\n", article.Desc)
		fmt.Fprintf(responseWriter, "Content: %s\n", article.Content)

		json.NewEncoder(responseWriter).Encode(article)

	default:
		fmt.Fprintf(responseWriter, "POST and GET methods are included to API scope...")
	}
}

func updatePage(responseWriter http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		http.ServeFile(responseWriter, request, "./templates/update.html")

	case "POST":
		if err := request.ParseForm(); err != nil {
			fmt.Fprintf(responseWriter, "ParseForm() err: %v", err)
			return
		}
		fmt.Fprintf(responseWriter, "Post from website! r.PostFrom = %v\n", request.PostForm)

		article := Article{
			request.FormValue("title"),
			request.FormValue("description"),
			request.FormValue("content")}

		fmt.Fprintf(responseWriter, "Title: %s\n", article.Title)
		fmt.Fprintf(responseWriter, "Description: %s\n", article.Desc)
		fmt.Fprintf(responseWriter, "Content: %s\n", article.Content)

		json.NewEncoder(responseWriter).Encode(article)

	default:
		fmt.Fprintf(responseWriter, "POST and GET methods are included to API scope...")
	}
}

func allArticles(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Access-Control-Allow-Origin", "*")

	collection := mongoSession.session.DB(database).C(collection)

	results := []Article{}
	collection.Find(bson.M{"title": bson.RegEx{"", ""}}).All(&results)
	jsonString, err := json.Marshal(results)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(responseWriter, string(jsonString))
}

func updateArticle(responseWriter http.ResponseWriter, request *http.Request) {
	collection := mongoSession.session.DB(database).C(collection)

	var updateVar = request.FormValue("findByTitle")

	article := Article{
		request.FormValue("title"),
		request.FormValue("description"),
		request.FormValue("content")}

	err := collection.Update(bson.M{"title": updateVar}, article)
	if err != nil {
		panic(err)
	}

	jsonString, err := json.Marshal(article)
	if err != nil {
		http.Error(responseWriter, err.Error(), 500)
		return
	}

	responseWriter.Header().Set("content-type", "application/json")

	responseWriter.Write(jsonString)

}

func saveArticle(responseWriter http.ResponseWriter, request *http.Request) {
	collection := mongoSession.session.DB(database).C(collection)

	article := Article{
		request.FormValue("title"),
		request.FormValue("description"),
		request.FormValue("content")}

	err := collection.Insert(article)
	if err != nil {
		panic(err)
	}

	jsonString, err := json.Marshal(article)
	if err != nil {
		http.Error(responseWriter, err.Error(), 500)
		return
	}

	responseWriter.Header().Set("content-type", "application/json")

	responseWriter.Write(jsonString)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/update", updatePage)
	myRouter.HandleFunc("/articles", allArticles).Methods("GET")
	myRouter.HandleFunc("/articles", saveArticle).Methods("POST")
	myRouter.HandleFunc("/articles", updateArticle).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8083", myRouter))
}

func main() {
	mongoSession.session = initialiseMongo()
	handleRequests()
}