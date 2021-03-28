package controllers

import (
	"encoding/json"
	"fmt"
	"go-rest-api/config"
	"go-rest-api/models"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type DeleteResponse struct {
	Id string `json:"id"`
}

func rootPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Go Api Server")
	fmt.Println("Root endpoint is hooked!")
}

func fetchAllPosts(w http.ResponseWriter, r *http.Request) {
	var posts []models.Post
	models.GetAllPosts(&posts)
	responseBody, err := json.Marshal(posts)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}

func fetchSinglePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var post models.Post
	models.GetSinglePost(&post, id)
	responseBody, err := json.Marshal(post)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var post models.Post
	if err := json.Unmarshal(reqBody, &post); err != nil {
		log.Fatal(err)
	}
	models.InsertPost(&post)
	responseBody, err := json.Marshal(post)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	models.DeletePost(id)
	responseBody, err := json.Marshal(DeleteResponse{Id: id})
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	reqBody, _ := ioutil.ReadAll(r.Body)

	var updatePost models.Post
	if err := json.Unmarshal(reqBody, &updatePost); err != nil {
		log.Fatal(err)
	}

	models.UpdatePost(&updatePost, id)
	convertUintId, _ := strconv.ParseUint(id, 10, 64)
	updatePost.Model.ID = uint(convertUintId)
	responseBody, err := json.Marshal(updatePost)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}

func StartWebServer() error {
	fmt.Println("Rest API with Mux Routers")
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", rootPage)
	router.HandleFunc("/posts", fetchAllPosts).Methods("GET")
	router.HandleFunc("/post/{id}", fetchSinglePost).Methods("GET")

	router.HandleFunc("/post", createPost).Methods("POST")
	router.HandleFunc("/post/{id}", deletePost).Methods("DELETE")
	router.HandleFunc("/post/{id}", updatePost).Methods("PUT")

	return http.ListenAndServe(fmt.Sprintf(":%d", config.Config.ServerPort), router)
}
