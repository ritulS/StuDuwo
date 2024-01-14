package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Db_client *mongo.Client

func Get_Collection(name string) *mongo.Collection {
	return Db_client.Database("rental_db").Collection(name)
}

func Init_db() error {

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable.")
	}
	var err error
	Db_client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to MongoDB!")

	return nil
}

func Close_db() {
	err := Db_client.Disconnect(context.TODO())
	if err != nil {
		panic(err)
	}
}
