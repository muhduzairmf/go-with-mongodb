package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load MONGO_URI from .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	// Try to initialize MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	// Set up time out connection if database has an error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Try to connect MongoDB
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer client.Disconnect(ctx)

	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	fmt.Println(databases)
}