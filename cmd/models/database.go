package models

import "go.mongodb.org/mongo-driver/mongo"

// Database contains mongo.Database
type Database struct {
	database *mongo.Database
}

type Client struct {
	client *mongo.Client
}

type Session struct {
	collection *mongo.Collection
	filter     interface{}
	limit      int64
	project    interface{}
	skip       int64
	sort       interface{}
}

// Collection contains mongo.Collection
type Collection struct {
	collection *mongo.Collection
}
