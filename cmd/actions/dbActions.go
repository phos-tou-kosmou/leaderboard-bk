package actions

import (
	"context"
	"encoding/json"
	"database/sql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "leaderboard-bk/cmd/models"
	_ "log"
	_ "time"
)

// Database contains mongo.Database
type DB interface {
	Beginx() (Tx, error)
	Close() error
}

type Client struct {
	client *sql.Conn
}

// Collection contains mongo.Collection
type Collection struct {
	collection *mongo.Collection
}

// Session -
type Session struct {
	collection *mongo.Collection
	filter     interface{}
	limit      int64
	project    interface{}
	skip       int64
	sort       interface{}
}

// Limit sets sorting
func (s *Session) Limit(limit int64) *Session {
	s.limit = limit
	return s
}

// Project sets sorting
func (s *Session) Project(project interface{}) *Session {
	s.project = project
	return s
}

// Skip sets sorting
func (s *Session) Skip(skip int64) *Session {
	s.skip = skip
	return s
}

// Sort sets sorting
func (s *Session) Sort(sort interface{}) *Session {
	s.sort = sort
	return s
}

// Decode returns all docs
func (s *Session) Decode(result interface{}) error {
	opts := options.Find()
	if s.limit > 0 {
		opts.SetLimit(s.limit)
	}
	if s.project != nil {
		opts.SetProjection(s.project)
	}
	if s.skip > 0 {
		opts.SetSkip(s.skip)
	}
	if s.sort != nil {
		opts.SetSort(s.sort)
	}
	cur, err := s.collection.Find(nil, s.filter, opts)
	if err != nil {
		return err
	}

	var docs []bson.M
	for cur.Next(nil) {
		var doc bson.M
		_ = cur.Decode(&doc)
		docs = append(docs, doc)
	}
	b, _ := json.Marshal(docs)
	_ = json.Unmarshal(b, result)
	return nil
}

type myDB struct {
	db DB
}

func NewDatabase(dsn string) (DB, error) {
	db, err := sql.Open("mysql", "theUser:thePassword@/theDbName")
	if err != nil {
		panic(err)
	}
	return db, nil
	//return &myDB{*db}, nil
}

// Database returns database
func (c *Client) Database(db string) *Database {
	return &Database{database: c.client.Database(db)}
}

// Collection returns database
func (d *Database) Collection(collection string) *Collection {
	return &Collection{collection: d.database.Collection(collection)}
}

// Find finds docs by given filter
func (c *Collection) Find(filter interface{}) *Session {
	return &Session{filter: filter, collection: c.collection}
}