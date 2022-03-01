package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"time"
)

type singleton struct {
	ctx context.Context
}

func (s singleton) insert(dbName string, collectionName string, value interface{}) (*mongo.InsertOneResult, error) {
	client := s.getConnection()

	db := client.Database(dbName)

	collection := db.Collection(collectionName)

	return collection.InsertOne(GetDbInstance().getContext(), value)
}

// defined type with interface
type Db interface {
	// here will be methods
	getConnection() *mongo.Client

	getContext() context.Context

	insert(dbName string, collectionName string, value interface{}) (*mongo.InsertOneResult, error)
}

// declare variable
var instance *singleton = nil

// Get only one object
func GetDbInstance() Db {
	if instance == nil {
		instance = new(singleton)
	}
	return instance
}

func (s singleton) getContext() context.Context {
	return s.ctx
}

func (s singleton) getConnection() *mongo.Client {
	dbUri := os.Getenv("DB_URI")

	s.ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(s.ctx, options.Client().ApplyURI(dbUri))

	checkErr(err)

	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	return client
}
