package main

import (
	firestore "cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"fmt"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
    "time"
)

const FulfillmentCol = "fulfillments"

var opt = option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))

type Fulfillment struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/getFulfillments", getFulfillments).Methods("GET")
	router.HandleFunc("/getFulfillment/{id}", getFulfillment).Methods("GET")
	router.HandleFunc("/createFulfilment", createFulfilment).Methods("POST")
	router.HandleFunc("/updateFulfilment/{id}", updateFulfilment).Methods("PUT")
	router.HandleFunc("/fulfillment", endpointStdOut).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}

func getNowInMillisecond() int64 {
    return time.Now().UnixNano() / 1000000
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func endcodePostAndWrite(w http.ResponseWriter, i interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(i)
}

func setupFirestore() (context.Context, *firestore.Client) {
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, opt)
	panicOnError(err)

	client, err := app.Firestore(ctx)
	panicOnError(err)

	return ctx, client
}

func getFulfillments(w http.ResponseWriter, r *http.Request) {
	ctx, client := setupFirestore()
	defer client.Close()

	docs, err := client.Collection(FulfillmentCol).DocumentRefs(ctx).GetAll()
	panicOnError(err)

	if len(docs) == 0 {
		endcodePostAndWrite(w, &[]Fulfillment{})
		return
	}

	posts := make([]Fulfillment, len(docs))
	for index, doc := range docs {
		snapshot, _ := doc.Get(ctx)
		data := snapshot.Data()
		title, _ := data["title"].(string)
		body, _ := data["body"].(string)
		posts[index] = Fulfillment{
			ID:    doc.ID,
			Title: title,
			Body:  body,
		}
	}

	endcodePostAndWrite(w, posts)
}

func getFulfillment(w http.ResponseWriter, r *http.Request) {
	ctx, client := setupFirestore()
	defer client.Close()

	params := mux.Vars(r)
	id := params["id"]
	doc, err := client.Collection(FulfillmentCol).Doc(id).Get(ctx)
	if err != nil {
		log.Fatalln("No fulfillments found", err)
		endcodePostAndWrite(w, &Fulfillment{})
		return
	}
	fulfillment := doc.Data()
	endcodePostAndWrite(w, fulfillment)
}

func createFulfilment(w http.ResponseWriter, r *http.Request) {
	var fulfillment Fulfillment
	err := json.NewDecoder(r.Body).Decode(&fulfillment)
	panicOnError(err)

	ctx, client := setupFirestore()
	defer client.Close()

	doc, _, err := client.Collection(FulfillmentCol).Add(
		ctx,
		map[string]interface{}{
			"title": fulfillment.Title,
			"body":  fulfillment.Body,
		},
	)
	panicOnError(err)

	newpost, err := doc.Get(ctx)
	panicOnError(err)
	data := newpost.Data()
	endcodePostAndWrite(w, data)
}

func updateFulfilment(w http.ResponseWriter, r *http.Request) {
	var fulfillment Fulfillment
	err := json.NewDecoder(r.Body).Decode(&fulfillment)
	panicOnError(err)

	ctx, client := setupFirestore()
	params := mux.Vars(r)
	id := params["id"]
	wresult, err := client.Collection(FulfillmentCol).Doc(id).Set(
		ctx,
		map[string]interface{}{
			"title": fulfillment.Title,
			"body":  fulfillment.Body,
		},
		firestore.MergeAll,
	)
	panicOnError(err)

	endcodePostAndWrite(w, &map[string]interface{}{
		"success":   true,
		"timestamp": wresult.UpdateTime,
	})
}

func endpointStdOut(w http.ResponseWriter, r *http.Request) {
	var bodyJson map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&bodyJson)
	if err != nil {
		fmt.Println("Error:", err)
	}

	jsonBytes, err := json.Marshal(bodyJson)
	panicOnError(err)

	ctx, client := setupFirestore()
	defer client.Close()

	_, _, err = client.Collection(FulfillmentCol).Add(
		ctx,
		map[string]interface{}{
			"fulfillment": string(jsonBytes),
            "timestamp": getNowInMillisecond(),
		},
	)
	panicOnError(err)

	endcodePostAndWrite(w, &map[string]interface{}{
		"greeting": "hello dialogflow",
	})
}
