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
const HomeAutoCol = "home-automation"

var opt = option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))

type Fulfillment struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/fulfillment", fulfillmentHandler).Methods("POST")

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

func handleTypeCheck(ok bool, msg string) {
	if !ok {
		fmt.Println(msg)
	}
}

func handleLighting(octx []interface{}) {
	if len(octx) == 0 {
		fmt.Println("No output context detected")
		return
	}
	context, _ := octx[0].(map[string]interface{})
	params, _ := context["parameters"].(map[string]interface{})
	room, ok := params["MyRooms.original"].(string)
	handleTypeCheck(ok, "Room is not OK!")
	state, ok := params["State.original"].(string)
	handleTypeCheck(ok, "State is not OK!")

	ctx, client := setupFirestore()
	defer client.Close()
	_, err := client.Collection(HomeAutoCol).Doc("lighting").Update(
		ctx,
		[]firestore.Update{{Path: room, Value: state}},
	)
	panicOnError(err)
}

func parseIntent(intent map[string]interface{}, octx []interface{}) {
	if intent["displayName"] == "Lighting" {
		handleLighting(octx)
	}
}

func fulfillmentHandler(w http.ResponseWriter, r *http.Request) {
	var bodyJson map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&bodyJson)
	if err != nil {
		fmt.Println("Error:", err)
	}

	ctx, client := setupFirestore()
	defer client.Close()
	jsonBytes, err := json.Marshal(bodyJson)
	panicOnError(err)
	_, _, err = client.Collection(FulfillmentCol).Add(
		ctx,
		map[string]interface{}{
			"fulfillment": string(jsonBytes),
			"timestamp":   getNowInMillisecond(),
		},
	)
	panicOnError(err)

	qresult, _ := bodyJson["queryResult"].(map[string]interface{})

	intent, _ := qresult["intent"].(map[string]interface{})
	octx, _ := qresult["outputContexts"].([]interface{})

	parseIntent(intent, octx)

	endcodePostAndWrite(w, &map[string]interface{}{
		"greeting": "hello dialogflow",
	})
}
