package main

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

const mongoHost = "db"
const mongoPort = "27017"
const mongoDb = "test"
const mongoCollection = "importantData"

type DB interface {
	Close()
	FindImportantData() (*[]Document, error)
	FindFeaturedData() (*[]Document, error)
	FindOne(primitive.ObjectID) (*Document, error)
	VoteUp(primitive.ObjectID) error
}

type mongoDB struct {
	client *mongo.Client
}

func NewMongoDB() (*mongoDB, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	c, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + mongoHost + ":" + mongoPort))
	if err != nil {
		return nil, err
	}
	err = c.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = c.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	log.Printf("Initialized connection to mongo")
	return &mongoDB{c}, nil
}

func (m *mongoDB) Close() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	m.client.Disconnect(ctx)
}

func (m *mongoDB) FindImportantData() (*[]Document, error) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	cursor, err := m.client.Database(mongoDb).Collection(mongoCollection).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var results []Document
	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, err
	}

	return &results, nil
}

func (m *mongoDB) FindFeaturedData() (*[]Document, error) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	limit := int64(5)
	opts := &options.FindOptions{
		Sort:  bson.M{"upvotes": -1},
		Limit: &limit,
	}
	cursor, err := m.client.Database(mongoDb).Collection(mongoCollection).Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	var results []Document
	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, err
	}

	return &results, nil
}

func (m *mongoDB) FindOne(id primitive.ObjectID) (*Document, error) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	var doc Document
	err := m.client.Database(mongoDb).Collection(mongoCollection).FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		return nil, err
	}

	return &doc, nil
}

func (m *mongoDB) VoteUp(id primitive.ObjectID) error {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{"upvotes": 1}}

	coll := m.client.Database(mongoDb).Collection(mongoCollection)
	ur, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if ur.ModifiedCount <= 0 {
		return errors.New("no document was modified")
	}

	return nil
}

