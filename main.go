package main

import (
	"encoding/json"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

var opt = option.WithCredentialsFile("path/to/serviceAccountKey.json")

type Post struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

var posts []Post

func main() {
	app, err := firebase.NewApp(context.Background(), nil, opt)
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8000"
	}
	router := mux.NewRouter()
	router.HandleFunc("/getPosts", getPosts).Methods("GET")
	router.HandleFunc("/getPost/{id}", getPost).Methods("GET")
	router.HandleFunc("/createPost", createPost).Methods("POST")
	router.HandleFunc("/updatePost/{id}", updatePost).Methods("PUT")

	posts = append(posts, Post{
		ID:    "1000000",
		Title: "Woot",
		Body:  "Something",
	})
	posts = append(posts, Post{
		ID:    "1000001",
		Title: "Naniii",
		Body:  "asdfaw890 a8s0df 8f0s8 ",
	})
	posts = append(posts, Post{
		ID:    "1000002",
		Title: "Waaaaahhhh",
		Body:  "-9df0asdf8sa8fs0",
	})

	fmt.Println("Listening on: " + port)
	http.ListenAndServe(port, router)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func endcodePostAndWrite(w http.ResponseWriter, i interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(i)
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	endcodePostAndWrite(w, posts)
}

func getPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range posts {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	endcodePostAndWrite(w, &Post{})
}

func createPost(w http.ResponseWriter, r *http.Request) {
	var post Post
	err := json.NewDecoder(r.Body).Decode(&post)
	handleError(err)

	post.ID = strconv.Itoa(rand.Intn(1000000))
	posts = append(posts, post)
	endcodePostAndWrite(w, post)
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println(params)
}
