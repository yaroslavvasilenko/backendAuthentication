package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func InsertUser(client *mongo.Client) *mongo.Collection {
	collection := client.Database("new_database").Collection("cust")
	return collection
}

func FindMongo(collection *mongo.Collection) Trainer {
	filter := Trainer{"Ash", 10, "Pallet Town"}
	var call Trainer
	err := collection.FindOne(context.TODO(), filter).Decode(&call)
	if err != nil {
		log.Fatal(err)
	}
	return call
}

type Trainer struct {
	Name string
	Age  int
	City string
}
