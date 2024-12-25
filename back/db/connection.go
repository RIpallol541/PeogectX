package db

import (
	"context"
	"log"
	
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Client *mongo.Client

// Connect initializes MongoDB connection
func Connect() {
	var err error
	Client, err = mongo.Connect(options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
		log.Fatal(err)
	}
}

// Disconnect closes MongoDB connection
func Disconnect() {
	if err := Client.Disconnect(context.TODO()); err != nil {
		log.Fatal(err)
	}
}
