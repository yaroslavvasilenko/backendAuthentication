package main

import (
	"github.com/yaroslavvasilenko/backendAuthentication/database"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

// We have some guid already = we send it in our request:
// localhost:8080/signin/?user-id=82e177bd14364bfea2425f63888e15f1
type applecation struct {
	serverMongo *mongo.Client
	userAuth    *mongo.Collection
}

func main() {
	clientMongo := database.Dbcall()
	registr := database.InsertUser(clientMongo)
	app := &applecation{
		serverMongo: clientMongo,
		userAuth:    registr,
	}
	// "Signin" and "Welcome" are the handlers that we will implement
	http.HandleFunc("/", app.handleGuid)
	// start the server on port 8000
	log.Fatal(http.ListenAndServe(":8080", nil))

}
