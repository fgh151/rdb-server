package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"time"
)

type database struct {
	ctx    context.Context
	client *mongo.Client
}

func (s database) insert(dbName string, collectionName string, value interface{}) (*mongo.InsertOneResult, error) {
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
var instance *database = nil

// Get only one object
func GetDbInstance() Db {
	if instance == nil {
		instance = new(database)
	}
	return instance
}

func (s database) getContext() context.Context {
	return s.ctx
}

func (s database) getConnection() *mongo.Client {

	if s.client != nil {
		return s.client
	}

	dbUri := os.Getenv("DB_URI")

	s.ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)

	var err error

	s.client, err = mongo.Connect(s.ctx, options.Client().ApplyURI(dbUri))

	checkErr(err)

	// Ping the primary
	if err := s.client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	return s.client
}
