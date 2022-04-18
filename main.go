package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Struct for all articles
type Article struct {
	Id      string `json:"id"`
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

var Articles []Article

//Home Page using ResponseWriter, greeting user with message, as well as notification in dev terminal relaying endpoint hit
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "WELCOME TO THE HOME PAGE")
	fmt.Println("Endpoint Hit: homePage")
}

//Helper function that will handle all requests from user and utilizing respective request functions
func handleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/all", returnAllArticles)
	myRouter.HandleFunc("/article/{id}", returnSingleArticle)
	myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
	myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

//All Articles
func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	json.NewEncoder(w).Encode(Articles)
}

//Single article by id
func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	for _, article := range Articles {
		if article.Id == key {
			json.NewEncoder(w).Encode(article)
		}
	}
}

//Create article using unmarshal in order to append to Articles
func createNewArticle(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article Article

	json.Unmarshal(reqBody, &article)

	Articles = append(Articles, article)

	json.NewEncoder(w).Encode(article)
}

//Deleting article by looping Articles for matching id
func deleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	for index, article := range Articles {
		if article.Id == id {
			Articles = append(Articles[:index], Articles[index+1:]...)
		}
	}
}

//Updating article mathcing id
func updateArticle(w http.ResponseWriter, r *http.Request) {
	articleId := mux.Vars(r)["id"]

	var updatedArticle Article

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the title and description only in order to update")
	}
	json.Unmarshal(reqBody, &updatedArticle)

	for i, singleArticle := range Articles {
		if singleArticle.Id == articleId {
			singleArticle.Title = updatedArticle.Title
			singleArticle.Desc = updatedArticle.Desc
			Articles = append(Articles[:i], singleArticle)
			json.NewEncoder(w).Encode(singleArticle)
		}
	}
}

//Main simply holds moc data and is being passed handleRequest in order to keep clean and DRY
func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	Articles = []Article{
		{Id: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content"},
		{Id: "2", Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
	}
	handleRequest()
}
