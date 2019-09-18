package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func main() {
	post := map[string]string{
		"title": "somethdfa asdfasdf af af fasdf32t4wefsd",
		"body":  "lorem50fasf fa fd fa fa 123124142",
	}

	// newPost(post)
    posts := getPosts()
    fmt.Println(posts)

    updatePost(1000000, post)
}

func handleError(err error) {
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
}

func getPosts() string {
	resp, err := http.Get("http://localhost:8000/getPosts")
	handleError(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	handleError(err)

	return string(body)
}

func newPost(post map[string]string) string {
	postByte, err := json.Marshal(post)
	handleError(err)

	resp, err := http.Post(
		"http://localhost:8000/createPost",
		"application/json",
		bytes.NewBuffer(postByte),
	)
	handleError(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	handleError(err)

    return string(body)
}

func updatePost(id int, post map[string]string) string {
	postByte, err := json.Marshal(post)
	handleError(err)
    
    reqString := "http://localhost:8000/updatePost/" + strconv.Itoa(id)

    client := &http.Client{}
	request, err := http.NewRequest(
		http.MethodPut,
        reqString,
		bytes.NewBuffer(postByte),
	)
	handleError(err)

    resp, err := client.Do(request)
	handleError(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	handleError(err)

	return string(body)
}
