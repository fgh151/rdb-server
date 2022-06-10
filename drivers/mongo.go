package drivers

import (
	"context"
	err2 "db-server/err"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
)

type Database struct {
	ctx    context.Context
	client *mongo.Client
}

func (s Database) Update(dbName string, collectionName string, id interface{}, value interface{}) (*mongo.UpdateResult, error) {
	client, _ := s.GetConnection()

	db := client.Database(dbName)

	collection := db.Collection(collectionName)

	return collection.UpdateByID(GetDbInstance().GetContext(), id, value)
}

func (s Database) Delete(dbName string, collectionName string, id interface{}) (*mongo.DeleteResult, error) {
	client, _ := s.GetConnection()

	db := client.Database(dbName)

	collection := db.Collection(collectionName)

	return collection.DeleteOne(GetDbInstance().GetContext(), bson.M{"_id": id})
}

func (s Database) Insert(dbName string, collectionName string, value interface{}) (*mongo.InsertOneResult, error) {
	client, _ := s.GetConnection()

	db := client.Database(dbName)

	collection := db.Collection(collectionName)

	return collection.InsertOne(GetDbInstance().GetContext(), value)
}

func (s Database) Find(dbName string, collectionName string, filter interface{}, limit int64, skip int64) ([]*bson.D, error) {
	client, _ := s.GetConnection()

	db := client.Database(dbName)

	collection := db.Collection(collectionName)

	var ctx = GetDbInstance().GetContext()
	var res []*bson.D

	findOptions := options.Find()

	findOptions.Limit = &limit
	findOptions.Skip = &skip

	cur, err := collection.Find(ctx, filter, findOptions)

	defer func() { _ = cur.Close(ctx) }()

	for cur.Next(ctx) {

		var d bson.D
		err := cur.Decode(&d)
		err2.PanicErr(err)

		res = append(res, &d)
	}

	return res, err
}

func (s Database) List(dbName string, collectionName string, limit int64, skip int64, order int, sort string, filter bson.D) ([]*bson.D, int64, error) {

	client, _ := s.GetConnection()

	db := client.Database(dbName)

	collection := db.Collection(collectionName)

	findOptions := options.Find()

	findOptions.Limit = &limit
	findOptions.Skip = &skip
	findOptions.SetSort(bson.D{{sort, order}})

	//filter := bson.D{{}};

	var ctx = GetDbInstance().GetContext()
	var res []*bson.D

	cur, err := collection.Find(ctx, filter, findOptions)

	defer func() { _ = cur.Close(ctx) }()

	for cur.Next(ctx) {

		var d bson.D
		err := cur.Decode(&d)
		err2.PanicErr(err)

		res = append(res, &d)
	}

	count, err := collection.CountDocuments(ctx, bson.D{{}})

	return res, count, err
}

// Db Document oriented data base interface
type Db interface {
	GetConnection() (*mongo.Client, error)

	GetContext() context.Context

	Insert(dbName string, collectionName string, value interface{}) (*mongo.InsertOneResult, error)

	Find(dbName string, collectionName string, filter interface{}) ([]*bson.D, error)

	List(dbName string, collectionName string) ([]*bson.D, error)

	Update(dbName string, collectionName string, value interface{}) (*mongo.UpdateResult, error)
}

// declare variable
var instance *Database

// GetDbInstance Db get Document oriented data base
func GetDbInstance() *Database {
	if instance == nil {
		instance = new(Database)
	}
	return instance
}

func (s Database) GetContext() context.Context {
	return s.ctx
}

func (s Database) GetConnection() (*mongo.Client, error) {

	if s.client != nil {
		return s.client, nil
	}

	conn := NewMongoConnectionFromEnv()

	s.ctx = context.TODO()

	var err error

	s.client, err = mongo.Connect(s.ctx, options.Client().ApplyURI(conn.GetDsn()))

	// Ping the primary
	if err := s.client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	return s.client, err
}

// MongoConnection Struect to connect to mongo db
type MongoConnection struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
}

// NewMongoConnectionFromEnv Create new mongo connection
func NewMongoConnectionFromEnv() MongoConnection {
	return MongoConnection{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DbName:   os.Getenv("DB_DBNAME"),
	}
}

// GetDsn Get mongo db dsn
func (c MongoConnection) GetDsn() string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/?maxPoolSize=20&w=majority",
		c.User,
		c.Password,
		c.Host,
		c.Port,
	)
}

// GetMongoSort Convert sort order from sql db to mongo syntax
func GetMongoSort(sqlSort string, sqlOrder string) (int, string) {
	var sort int
	var order = sqlOrder
	if sqlSort == "ASC" {
		sort = 1
	} else {
		sort = -1
	}
	if sqlOrder == "id" {
		order = "_id"
	}

	return sort, order
}
