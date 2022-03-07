package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Modelling the data based on the bson
type Creator struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	Name string `bson:"name,omitempty"`
	Description string `bson:"description,omitempty"`
}

type Video struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	Title string `bson:"title,omitempty"`
	Tags []string `bson:"tags,omitempty"`
	CreatorId interface{} `bson:"creator_id,omitempty"`
	Duration int32 `bson:"duration,omitempty"`
}

func createData(ctx context.Context, videoCollection *mongo.Collection, creatorCollection *mongo.Collection)  {
	// #1 Insert one
	// creatorResult, err := creatorCollection.InsertOne(ctx, bson.D{
	// 	{Key: "name", Value: "freeCodeCamp"},
	// 	{Key: "description", Value: "Learn anything for free here!"},
	// })
	// or
	freeCodeCampCreator := Creator{
		Name: "freeCodeCamp",
		Description: "Learn anything for free here!",
	}
	creatorResult, err := creatorCollection.InsertOne(ctx, freeCodeCampCreator)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(creatorResult.InsertedID)

	fmt.Println("---\t---\t---\t---\t---")
	
	// #1 Insert many
	// videoResult, err := videoCollection.InsertMany(ctx, []interface{}{
	// 	bson.D{
	// 		{"title", "React Tutorial for beginners"},
	// 		// The first (title) is a key, and the second is value
	// 		{"tags", bson.A{"programming", "tutorial", "tech"}},
	// 		// bson.A is represents array in MongoDB
	// 		{"duration", 300},
	// 		{"creator_id", creatorResult.InsertedID},
	// 	},
	// 	bson.D{
	// 		{"title", "Vue Tutorial for beginners"},
	// 		{"tags", bson.A{"programming", "tutorial", "tech"}},
	// 		{"duration", 280},
	// 		{"creator_id", creatorResult.InsertedID},
	// 	},
	// 	bson.D{
	// 		{"title", "Node.js Tutorial for beginners"},
	// 		{"tags", bson.A{"programming", "tutorial", "tech"}},
	// 		{"duration", 500},
	// 		{"creator_id", creatorResult.InsertedID},
	// 	},
	// })
	// or
	var freeCodeCampVideos [3]Video
	freeCodeCampVideos[0] = Video{
		Title: "React Tutorial for beginners",
		Tags: []string{"programming", "tutorial", "tech"},
		CreatorId: creatorResult.InsertedID,
		Duration: 300,
	}
	freeCodeCampVideos[1] = Video{
		Title: "Vue Tutorial for beginners",
		Tags: []string{"programming", "tutorial", "tech"},
		CreatorId: creatorResult.InsertedID,
		Duration: 280,
	}
	freeCodeCampVideos[2] = Video{
		Title: "Node.js Tutorial for beginners",
		Tags: []string{"programming", "tutorial", "tech"},
		CreatorId: creatorResult.InsertedID,
		Duration: 500,
	}
	fCCVideos := make([]interface{}, len(freeCodeCampVideos))
	for index, video := range freeCodeCampVideos {
		fCCVideos[index] = video
	}
	videoResult, err := videoCollection.InsertMany(ctx, fCCVideos)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(videoResult.InsertedIDs...)
}

func readData(ctx context.Context, videoCollection *mongo.Collection, creatorCollection *mongo.Collection)  {
	cursor, err := videoCollection.Find(ctx, bson.M{})
	if err != nil { log.Fatal(err.Error()) }
	defer cursor.Close(ctx)

	// #1 Get all documents from a collection, good for small dataset
	var videos []bson.M
	// or var videos []Videos
	err = cursor.All(ctx, &videos)
	if err != nil { log.Fatal(err.Error()) }

	for _, v := range videos {
		fmt.Println(v)
	}

	fmt.Println("---\t---\t---\t---\t---")

	// #2 Get all documents from a collection, good for big dataset
	for cursor.Next(ctx) {
		var video bson.M
		err = cursor.Decode(&video)
		if err != nil { log.Fatal(err.Error()) }
		fmt.Println(video)
	}

	fmt.Println("---\t---\t---\t---\t---")

	// #1 Get a single document from a collection
	var creator bson.M
	// or var creator Creator
	err = creatorCollection.FindOne(ctx, bson.M{}).Decode(&creator)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(creator)

	fmt.Println("---\t---\t---\t---\t---")

	// #1 Get several documents from a collection, with filter
	filterCursor, err := videoCollection.Find(ctx, bson.M{"tags": "CrashCourse"})
	if err != nil {
		log.Fatal(err.Error())
	}

	var videosFiltered []bson.M
	// or var videosFiltered Video
	err = filterCursor.All(ctx, &videosFiltered)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(videosFiltered)

	fmt.Println("---\t---\t---\t---\t---")

	// #1 Get sorted several documents from collection
	opts := options.Find()
	opts.SetSort(bson.D{{"duration", 1}})
	// 1 for ascending, -1 for descending
	sortCursor, err := videoCollection.Find(ctx, bson.D{
		{"duration", bson.D{
			{"$gt", 40},
		}},
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	var videosSorted []bson.M
	err = sortCursor.All(ctx, &videosSorted)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(videosSorted)
}

func updateData(ctx context.Context, videoCollection *mongo.Collection, creatorCollection *mongo.Collection)  {
	// #1 Update one document by the id
	id, err := primitive.ObjectIDFromHex("62201a0d854b848e2951ed61")
	if err != nil {
		log.Fatal(err.Error())
	}

	result, err := creatorCollection.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.D{
			{"$set", bson.D{
				{"description", "I love crash course"},
			}},
		},
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Updated %v Document(s)\n", result.ModifiedCount)

	fmt.Println("---\t---\t---\t---\t---")

	// #1 Update many documents
	creator_id, err := primitive.ObjectIDFromHex("6220d932b89c96822aa90acc")
	if err != nil {
		log.Fatal(err.Error())
	}

	result, err = videoCollection.UpdateMany(ctx, 
		bson.M{"creator_id": creator_id},
		bson.D{
			{"$set", bson.D{
				{"tags", bson.A{"tutorial", "freeCodeCamp", "js_framework"},},
			}},
		},
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Updated %v Document(s)\n", result.ModifiedCount)

	fmt.Println("---\t---\t---\t---\t---")

	// #1 Replace entire one document
	result, err = creatorCollection.ReplaceOne(ctx,
		bson.M{"name": "Traversy Media"},
		bson.M{
			"name": "Brad Traversy",
			"description": "Crash Course anything!",
		},
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Updated %v Document(s)\n", result.ModifiedCount)
}

func deleteData(ctx context.Context, videoCollection *mongo.Collection, creatorCollection *mongo.Collection)  {
	// #1 Delete many documents
	creator_id, err := primitive.ObjectIDFromHex("6220d932b89c96822aa90acc")
	if err != nil {
		log.Fatal(err.Error())
	}

	result, err := videoCollection.DeleteMany(ctx, bson.M{"creator_id": creator_id})
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Delete many documents successfull. %v", result.DeletedCount)

	// #1 Delete one document
	result, err = creatorCollection.DeleteOne(ctx, bson.M{"name": "freeCodeCamp"})
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Delete one document successfull. %v", result.DeletedCount)

	// #1 Delete/Drop the collection
	err = videoCollection.Drop(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	// Load MONGO_URI from .env
	err := godotenv.Load()
	if err != nil { log.Fatal(err.Error()) }

	// Try to initialize MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil { log.Fatal(err.Error()) }

	// Set up time out connection if has an error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Try to connect MongoDB
	err = client.Connect(ctx)
	if err != nil { log.Fatal(err.Error()) }
	defer client.Disconnect(ctx)

	// Testing connection
	// databases, err := client.ListDatabaseNames(ctx, bson.M{})
	// if err != nil {
	// 	log.Fatal(err.Error())
	// 	return
	// }
	// fmt.Println(databases)

	// Get the database and collections
	gomongodbDatabase := client.Database("go-mongodb")
	videoCollection := gomongodbDatabase.Collection("video")
	creatorCollection := gomongodbDatabase.Collection("creator")

	createData(ctx, videoCollection, creatorCollection)

	//readData(ctx, videoCollection, creatorCollection)

	//updateData(ctx, videoCollection, creatorCollection)

	//deleteData(ctx, videoCollection, creatorCollection)
	
}