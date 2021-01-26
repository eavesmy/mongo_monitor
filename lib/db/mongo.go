package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
)

var Mongo *mongo.Client

var URI = ""

var Collections = map[string]*mongo.Collection{}

func InitDB() *mongo.Client {

	client, err := mongo.NewClient(options.Client().ApplyURI(URI))

	if err != nil {
		panic(err)
	}

	Mongo = client

	return client
}

func Register(dbname, collectionname string, ctx context.Context) *mongo.Collection {

	_name := dbname + "." + collectionname

	if Collections[_name] != nil {
		return Collections[_name]
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(URI))

	if err != nil || client.Connect(ctx) != nil {
		panic(err)
	}

	c := client.Database(dbname).Collection(collectionname)

	Collections[_name] = c

	return c
}
