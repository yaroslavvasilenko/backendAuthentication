package main

import (
	"crypto/rsa"
	"crypto/x509"
	"github.com/yaroslavvasilenko/backendAuthentication/database"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"log"
	"net/http"
	"os"
)

// We have some guid already = we send it in our request:
// localhost:8080/signin/?user-id=82e177bd14364bfea2425f63888e15f1

type Application struct {
	ServerMongo *mongo.Client
	UserAuth    *mongo.Collection
	Secret      *rsa.PrivateKey
}

func main() {
	secret, _ := x509.ParsePKCS1PrivateKey(downloadKey())
	clientMongo := database.ConnectMongoBD()
	register := database.CollectionMongo(clientMongo)
	App := &Application{
		ServerMongo: clientMongo,
		UserAuth:    register,
		Secret:      secret,
	}

	http.HandleFunc("/", App.firstRoute)
	http.HandleFunc("/insert", App.secondRoute)
	// start the server on port 8000
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func downloadKey() []byte {
	file, _ := os.Open("key")
	defer file.Close()
	key, _ := io.ReadAll(file)
	return key
}
