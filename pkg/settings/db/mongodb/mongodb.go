package mongodb

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mongo constants
const (
	mongoURL = "mongodb://localhost:27017"
	// mongoURL = "mongodb://tiger:tiger123@ds263078.mlab.com:63078/erp-user-service?retryWrites=false&w=majority"
)

// Mongo variables
var (
	Client *mongo.Client
)

// InitMongoClient initiates the connection to mongodb server
func InitMongoClient() {
	var err error
	Client, err = mongo.NewClient(options.Client().ApplyURI(mongoURL))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = Client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	log.Println("Database successfully connected")
}
