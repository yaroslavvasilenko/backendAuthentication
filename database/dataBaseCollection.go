package database

import (
	"go.mongodb.org/mongo-driver/mongo"
)

func CollectionMongo(client *mongo.Client) *mongo.Collection {
	collection := client.Database("userAuth").Collection("users")
	return collection
}
