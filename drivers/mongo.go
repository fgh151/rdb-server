package drivers

import (
	"context"
	err2 "db-server/err"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
)

type database struct {
	ctx    context.Context
	client *mongo.Client
}

func (s database) Update(dbName string, collectionName string, id interface{}, value interface{}) (*mongo.UpdateResult, error) {
	client, _ := s.GetConnection()

	db := client.Database(dbName)

	collection := db.Collection(collectionName)

	return collection.UpdateByID(GetDbInstance().GetContext(), id, value)
}

func (s database) Delete(dbName string, collectionName string, id interface{}) (*mongo.DeleteResult, error) {
	client, _ := s.GetConnection()

	db := client.Database(dbName)

	collection := db.Collection(collectionName)

	return collection.DeleteOne(GetDbInstance().GetContext(), bson.M{"_id": id})
}

func (s database) Insert(dbName string, collectionName string, value interface{}) (*mongo.InsertOneResult, error) {
	client, _ := s.GetConnection()

	db := client.Database(dbName)

	collection := db.Collection(collectionName)

	return collection.InsertOne(GetDbInstance().GetContext(), value)
}

func (s database) Find(dbName string, collectionName string, filter interface{}) ([]*bson.D, error) {
	client, _ := s.GetConnection()

	db := client.Database(dbName)

	collection := db.Collection(collectionName)

	var ctx = GetDbInstance().GetContext()
	var res []*bson.D

	cur, err := collection.Find(ctx, filter)

	defer func() { _ = cur.Close(ctx) }()

	for cur.Next(ctx) {

		var d bson.D
		err := cur.Decode(&d)
		err2.PanicErr(err)

		res = append(res, &d)
	}

	return res, err
}

func (s database) List(dbName string, collectionName string) ([]*bson.D, error) {
	client, _ := s.GetConnection()

	db := client.Database(dbName)

	collection := db.Collection(collectionName)

	findOptions := options.Find()

	var ctx = GetDbInstance().GetContext()
	var res []*bson.D

	cur, err := collection.Find(ctx, bson.D{{}}, findOptions)

	defer func() { _ = cur.Close(ctx) }()

	for cur.Next(ctx) {

		var d bson.D
		err := cur.Decode(&d)
		err2.PanicErr(err)

		res = append(res, &d)
	}

	return res, err
}

type Db interface {
	GetConnection() (*mongo.Client, error)

	GetContext() context.Context

	Insert(dbName string, collectionName string, value interface{}) (*mongo.InsertOneResult, error)

	Find(dbName string, collectionName string, filter interface{}) ([]*bson.D, error)

	List(dbName string, collectionName string) ([]*bson.D, error)

	Update(dbName string, collectionName string, value interface{}) (*mongo.UpdateResult, error)
}

// declare variable
var instance *database = nil

func GetDbInstance() *database {
	if instance == nil {
		instance = new(database)
	}
	return instance
}

func (s database) GetContext() context.Context {
	return s.ctx
}

func (s database) GetConnection() (*mongo.Client, error) {

	if s.client != nil {
		return s.client, nil
	}

	dbUri := os.Getenv("DB_URI")

	s.ctx = context.TODO()

	var err error

	s.client, err = mongo.Connect(s.ctx, options.Client().ApplyURI(dbUri))

	// Ping the primary
	if err := s.client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	return s.client, err
}