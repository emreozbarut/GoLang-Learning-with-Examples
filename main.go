package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Article struct {
	Title string
	Desc string
	Content string
}

type Articles []Article

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

func allArticles(responseWriter http.ResponseWriter, request *http.Request) {
	articles := Articles{
		Article{Title:"Mock Title", Desc:"Mock Desc", Content:"Mock Content"},
	}

	fmt.Println("Endpoint invoked: All Articles Endpoint")
	json.NewEncoder(responseWriter).Encode(articles)
}

func saveArticle(responseWriter http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(responseWriter, "Save Article invoked...")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/articles", allArticles).Methods("GET")
	myRouter.HandleFunc("/articles", saveArticle).Methods("POST")

	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {
	handleRequests()
}