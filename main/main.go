package main

import (
	"github.com/yaroslavvasilenko/backendAuthentication/database"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

// We have some guid already = we send it in our request:
// localhost:8080/signin/?user-id=82e177bd14364bfea2425f63888e15f1
type Applecation struct {
	ServerMongo *mongo.Client
	UserAuth    *mongo.Collection
}

func main() {
	clientMongo := database.Dbcall()
	registr := database.InsertUser(clientMongo)
	App := &Applecation{
		ServerMongo: clientMongo,
		UserAuth:    registr,
	}
	// "Signin" and "Welcome" are the handlers that we will implement
	http.HandleFunc("/", App.firstRoute)
	http.HandleFunc("/insert", App.sekondRoute)
	// start the server on port 8000
	log.Fatal(http.ListenAndServe(":8080", nil))

}
