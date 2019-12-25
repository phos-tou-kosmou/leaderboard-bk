package actions

import (
	"context"
	"encoding/json"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (c *Client) UpdateOne(t *testing.T) {
	var err error
	var client *mongo.Client
	var collection *mongo.Collection
	var ctx = context.Background()
	var doc = bson.M{"_id": primitive.NewObjectID(), "hometown": "Atlanta"}
	var result *mongo.UpdateResult
	client = getMongoClient()
	defer client.Disconnect(ctx)
	collection = client.Database(dbName).Collection(collectionName)
	if _, err = collection.InsertOne(ctx, doc); err != nil {
		t.Fatal(err)
	}
	var update bson.M
	_ = json.Unmarshal([]byte(`{ "$set": {"year": 1998}}`), &update)
	if result, err = collection.UpdateOne(ctx, bson.M{"_id": doc["_id"]}, update); err != nil {
		t.Fatal(err)
	}

	if result.ModifiedCount != 1 {
		t.Fatal("update failed, expected 1 but got", result.ModifiedCount)
	}
	_, _ = collection.DeleteMany(ctx, bson.M{"hometown": "Atlanta"})
}

func (c *Client) UpdateMany(t *testing.T) {
	var err error
	var client *mongo.Client
	var collection *mongo.Collection
	var ctx = context.Background()
	var docs []interface{}
	docs = append(docs, bson.M{"_id": primitive.NewObjectID(), "hometown": "Atlanta", "counter": int32(1)})
	docs = append(docs, bson.M{"_id": primitive.NewObjectID(), "hometown": "Atlanta", "counter": int32(2)})
	var result *mongo.UpdateResult
	client = getMongoClient()
	defer client.Disconnect(ctx)
	collection = client.Database(dbName).Collection(collectionName)
	if _, err = collection.InsertMany(ctx, docs); err != nil {
		t.Fatal(err)
	}
	var update bson.M
	_ = json.Unmarshal([]byte(`{ "$set": {"year": 1998}}`), &update)
	if result, err = collection.UpdateMany(ctx, bson.M{"hometown": "Atlanta"}, update); err != nil {
		t.Fatal(err)
	}
	if result.ModifiedCount != 2 {
		t.Fatal("update failed, expected 2 but got", result.ModifiedCount)
	}
	_, _ = collection.DeleteMany(ctx, bson.M{"hometown": "Atlanta"})
}
