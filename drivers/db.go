package drivers

import (
	"context"
	err2 "db-server/err"
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

func (s database) Insert(dbName string, collectionName string, value interface{}) (*mongo.InsertOneResult, error) {
	client := s.GetConnection()

	db := client.Database(dbName)

	collection := db.Collection(collectionName)

	return collection.InsertOne(GetDbInstance().GetContext(), value)
}

// defined type with interface
type Db interface {
	// here will be methods
	GetConnection() *mongo.Client

	GetContext() context.Context

	Insert(dbName string, collectionName string, value interface{}) (*mongo.InsertOneResult, error)
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

func (s database) GetContext() context.Context {
	return s.ctx
}

func (s database) GetConnection() *mongo.Client {

	if s.client != nil {
		return s.client
	}

	dbUri := os.Getenv("DB_URI")

	s.ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)

	var err error

	s.client, err = mongo.Connect(s.ctx, options.Client().ApplyURI(dbUri))

	err2.CheckErr(err)

	// Ping the primary
	if err := s.client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	return s.client
}
