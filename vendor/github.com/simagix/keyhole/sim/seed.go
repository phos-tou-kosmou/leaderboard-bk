// Copyright 2018 Kuei-chun Chen. All rights reserved.

package sim

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/simagix/gox"
	"github.com/simagix/keyhole/sim/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Feeder seeds feeder
type Feeder struct {
	collection   string
	database     string
	file         string
	isDrop       bool
	showProgress bool
	total        int
}

// Model - robot model
type Model struct {
	ID          string `json:"_id" bson:"_id"`
	Name        string
	Description string
	Year        int
}

// Task - robot task
type Task struct {
	For         string `json:"for" bson:"for"`
	MinutesUsed int    `json:"minutesUsed" bson:"minutesUsed"`
}

// Robot -
type Robot struct {
	ID         string  `json:"_id" bson:"_id"`
	ModelID    string  `json:"modelId,omitempty" bson:"modelId,omitempty"`
	Notes      string  `json:"notes" bson:"notes"`
	BatteryPct float32 `json:"batteryPct,omitempty" bson:"batteryPct,omitempty"`
	Tasks      []Task  `json:"tasks" bson:"tasks"`
}

// NewFeeder establish seeding parameters
func NewFeeder() *Feeder {
	return &Feeder{isDrop: false, total: 1000, showProgress: true}
}

// SetCollection set collection
func (f *Feeder) SetCollection(collection string) {
	f.collection = collection
}

// SetDatabase set database
func (f *Feeder) SetDatabase(database string) {
	f.database = database
}

// SetFile set file
func (f *Feeder) SetFile(file string) {
	f.file = file
}

// SetIsDrop set isDrop
func (f *Feeder) SetIsDrop(isDrop bool) {
	f.isDrop = isDrop
}

// SetShowProgress set showProgress
func (f *Feeder) SetShowProgress(showProgress bool) {
	f.showProgress = showProgress
}

// SetTotal set total
func (f *Feeder) SetTotal(total int) {
	f.total = total
}

// SeedData seeds all demo data
func (f *Feeder) SeedData(client *mongo.Client) error {
	if f.file == "" {
		return f.SeedAllDemoData(client)
	} else if f.collection == "" {
		return errors.New("usage: keyhole --uri connection_uri --seed [--file filename --collection collection_name]")
	}
	return f.seedFromTemplate(client)
}

// SeedAllDemoData - seed data for demo
//  models: {
//    "_id": string
//   "name": string,
//   "description": string
//   "year": integer
// }
//
// robots: {
//   "_id": string
//   "modelId": string
//   "notes": string,
//   "batteryPct": float,
//   "tasks": [{"for": string, "minutesUsed": integer}]
// }
func (f *Feeder) SeedAllDemoData(client *mongo.Client) error {
	var err error
	if err = f.SeedFavorites(client); err != nil {
		return err
	}
	if err = f.seedRobots(client); err != nil {
		return err
	}
	if err = f.seedNumbers(client); err != nil {
		return err
	}
	if err = f.SeedCars(client); err != nil {
		return err
	}
	return err
}

// SeedFavorites seeds demo data of collection favorites
func (f *Feeder) SeedFavorites(client *mongo.Client) error {
	var err error
	var ctx = context.Background()
	c := client.Database(f.database).Collection("lookups")
	favoritesCollection := client.Database(f.database).Collection("favorites")
	if f.isDrop {
		if err = c.Drop(ctx); err != nil {
			return err
		}
		if err = favoritesCollection.Drop(ctx); err != nil {
			return err
		}
	}

	for i := 0; i < 10; i++ {
		c.InsertOne(ctx, bson.M{"_id": i + 1000, "type": "sports", "name": util.Favorites.Sports[i]})
		c.InsertOne(ctx, bson.M{"_id": i + 1100, "type": "book", "name": util.Favorites.Books[i]})
		c.InsertOne(ctx, bson.M{"_id": i + 1200, "type": "movie", "name": util.Favorites.Movies[i]})
		c.InsertOne(ctx, bson.M{"_id": i + 1300, "type": "city", "name": util.Favorites.Cities[i]})
		c.InsertOne(ctx, bson.M{"_id": i + 1400, "type": "music", "name": util.Favorites.Music[i]})
	}
	favoritesCount := f.seedCollection(favoritesCollection, 2)
	fmt.Printf("Seeded favorites: %d\n", favoritesCount)
	return err
}

func (f *Feeder) seedRobots(client *mongo.Client) error {
	var err error
	var ctx = context.Background()
	modelsCollection := client.Database(f.database).Collection("models")
	robotsCollection := client.Database(f.database).Collection("robots")
	if f.isDrop {
		modelsCollection.Drop(ctx)
		robotsCollection.Drop(ctx)
	}

	for i := 1000; i < 1010; i++ {
		model := "model-" + fmt.Sprintf("%x", (rand.Intn(5000)+5000)*i)
		name := fmt.Sprintf("Robo %d-%x", i, rand.Intn(1000000))
		descr := fmt.Sprintf("%s %s", model, name)
		year := time.Now().Year() - rand.Intn(5)
		if _, err = modelsCollection.InsertOne(ctx, &Model{model, name, descr, year}); err != nil {
			log.Fatal(err)
		}

		for r := 0; r < 2+rand.Intn(20); r++ {
			id := "robot-" + fmt.Sprintf("%x", (rand.Intn(5000)+5000)*r)
			notes := fmt.Sprintf("%s %s", id, model)
			pct := rand.Float32()
			tasks := []Task{{"Business", 10 + rand.Intn(60)}, {"Home", 10 + rand.Intn(60)}}
			if _, err = robotsCollection.InsertOne(ctx, &Robot{id, model, notes, pct, tasks}); err != nil {
				continue
			}
		}
	}
	modelsCount, _ := modelsCollection.CountDocuments(ctx, bson.M{})
	robotsCount, _ := robotsCollection.CountDocuments(ctx, bson.M{})
	fmt.Printf("Seeded models: %d, robots: %d\n", modelsCount, robotsCount)
	return err
}

func (f *Feeder) seedNumbers(client *mongo.Client) error {
	var err error
	var ctx = context.Background()
	numbersCollection := client.Database(f.database).Collection("numbers")
	if f.isDrop {
		numbersCollection.Drop(ctx)
	}

	var docs []interface{}
	for n := 0; n < 1000; n++ {
		docs = append(docs, bson.M{"a": rand.Intn(100), "b": rand.Intn(100), "c": rand.Intn(100)})
	}
	if _, err = numbersCollection.InsertMany(ctx, docs); err != nil {
		log.Fatal(err)
	}
	// create index example
	indexView := numbersCollection.Indexes()
	idx := mongo.IndexModel{
		Keys: bson.D{{Key: "a", Value: 1}},
	}
	indexView.CreateOne(ctx, idx)
	idx = mongo.IndexModel{
		Keys: bson.D{{Key: "a", Value: 1}, {Key: "b", Value: 1}},
	}
	indexView.CreateOne(ctx, idx)
	idx = mongo.IndexModel{
		Keys: bson.D{{Key: "a", Value: 1}, {Key: "b", Value: -1}},
	}
	indexView.CreateOne(ctx, idx)
	idx = mongo.IndexModel{
		Keys: bson.D{{Key: "a", Value: 1}, {Key: "b", Value: 1}, {Key: "c", Value: 1}},
	}
	indexView.CreateOne(ctx, idx)

	numbersCount, _ := numbersCollection.CountDocuments(ctx, bson.M{})
	fmt.Printf("Seeded numbers: %d\n", numbersCount)
	return err
}

// SeedCars seeds cars collection
func (f *Feeder) SeedCars(client *mongo.Client) error {
	var err error
	var ctx = context.Background()
	carsCollection := client.Database(f.database).Collection("cars")
	dealersCollection := client.Database(f.database).Collection("dealers")
	employeesCollection := client.Database(f.database).Collection("employees")
	if f.isDrop {
		carsCollection.Drop(ctx)
		dealersCollection.Drop(ctx)
		employeesCollection.Drop(ctx)
	}

	// Upsert examples
	for i := 0; i < len(dealers); i++ {
		dealerID := fmt.Sprintf("DEALER-%d", 1+i)
		opts := options.Update()
		opts.SetUpsert(true)
		if _, err := dealersCollection.UpdateOne(ctx, bson.M{"_id": dealerID}, bson.M{"$set": bson.M{"name": dealers[i]}}, opts); err != nil {
			log.Fatal(err)
		}
	}

	var emp bson.M
	opts := options.Replace()
	opts.SetUpsert(true)
	var empID = int(1001)
	emp = getEmployee(empID, 0)
	empID++
	employeesCollection.ReplaceOne(ctx, bson.M{"_id": emp["_id"]}, emp, opts)
	for i := 0; i < 2; i++ {
		emp = getEmployee(empID, 1001)
		parent := empID
		employeesCollection.ReplaceOne(ctx, bson.M{"_id": emp["_id"]}, emp, opts)
		empID++
		for j := 0; j < 3; j++ {
			emp = getEmployee(empID, parent)
			pID := empID
			employeesCollection.ReplaceOne(ctx, bson.M{"_id": emp["_id"]}, emp, opts)
			empID++
			for k := 0; k < 5; k++ {
				emp = getEmployee(empID, pID)
				employeesCollection.ReplaceOne(ctx, bson.M{"_id": emp["_id"]}, emp, opts)
				empID++
			}
		}
	}

	// create index example
	indexView := carsCollection.Indexes()
	idx := mongo.IndexModel{
		Keys: bson.D{{Key: "filters.k", Value: 1}, {Key: "filters.v", Value: 1}},
	}
	indexView.CreateOne(ctx, idx)

	dealersCount, _ := dealersCollection.CountDocuments(ctx, bson.M{})
	carsCount := f.seedCollection(carsCollection, 1)
	fmt.Printf("Seeded cars: %d, dealers: %d\n", carsCount, dealersCount)
	return err
}

var dealers = []string{"Atlanta Auto", "Buckhead Auto", "Johns Creek Auto"}
var brands = []string{"Audi", "BMW", "Chevrolet", "Ford", "Honda", "Mercedes-Benz", "Nissan", "Porsche", "Toyota", "Volkswagen"}
var styles = []string{"Sedan", "Coupe", "Convertible", "Minivan", "SUV", "Truck"}
var colors = []string{"Beige", "Black", "Blue", "Brown", "Gold", "Gray", "Green", "Orange", "Pink", "Purple", "Red", "Silver", "White", "Yellow"}

func getVehicle() bson.M {
	curYear := time.Now().Year()
	delta := rand.Intn(8)
	year := curYear - delta
	used := true
	if delta == 0 {
		used = false
	}
	brand := brands[rand.Intn(len(styles))]
	color := colors[rand.Intn(len(colors))]
	style := styles[rand.Intn(len(styles))]

	return bson.M{
		"dealer": fmt.Sprintf("DEALER-%d", 1+rand.Intn(len(dealers))),
		"brand":  brand,
		"color":  color,
		"style":  style,
		"year":   year,
		"used":   used,
		"filters": []bson.M{
			bson.M{"k": "brand", "v": brand},
			bson.M{"k": "color", "v": color},
			bson.M{"k": "style", "v": style},
			bson.M{"k": "year", "v": year},
			bson.M{"k": "used", "v": used}},
	}
}

func (f *Feeder) seedCollection(c *mongo.Collection, fnum int) int {
	var err error
	var ctx = context.Background()
	var bsize = 1000
	var remaining = f.total

	for i := 0; i < f.total; {
		num := bsize
		if remaining < bsize {
			num = remaining
		}
		var contentArray []interface{}
		for n := 0; n < num; n++ {
			if fnum == 1 {
				contentArray = append(contentArray, getVehicle())
			} else if fnum == 2 {
				contentArray = append(contentArray, util.GetDemoDoc())
			}
			i++
			remaining--
		}
		if _, err = c.InsertMany(ctx, contentArray); err != nil {
			panic(err)
		}
		if f.showProgress {
			fmt.Fprintf(os.Stderr, "\r%3.1f%% ", float64(100*i)/float64(f.total))
		}
	}
	if f.showProgress {
		fmt.Fprintf(os.Stderr, "\r100%%\r     \r")
	}
	cnt, _ := c.CountDocuments(ctx, bson.M{})
	return int(cnt)
}

// SeedFromTemplate seeds data from a template in a file
func (f *Feeder) seedFromTemplate(client *mongo.Client) error {
	var err error
	var ctx = context.Background()
	var bsize = 100
	var remaining = f.total
	var sdoc bson.M
	if sdoc, err = util.GetDocByTemplate(f.file, true); err != nil {
		return err
	}
	bytes, _ := json.Marshal(sdoc)
	doc := make(map[string]interface{})
	json.Unmarshal(bytes, &doc)
	collName := f.collection
	if collName == "" {
		collName = "examples"
	}
	log.Println("Seed data to collection", collName)
	c := client.Database(f.database).Collection(collName)
	if f.isDrop {
		c.Drop(ctx)
	}

	var wg = gox.NewWaitGroup(runtime.NumCPU())
	for threadNum := 0; threadNum < f.total; threadNum += bsize {
		wg.Add(1)
		num := bsize
		if remaining < bsize {
			num = remaining
		}
		remaining -= num
		if f.showProgress {
			fmt.Fprintf(os.Stderr, "\r%3.1f%% ", float64(100*(f.total-remaining))/float64(f.total))
		}
		go func(threadNum int, num int) {
			defer wg.Done()
			var contentArray []interface{}
			for n := 0; n < num; n++ {
				ndoc := make(map[string]interface{})
				util.RandomizeDocument(&ndoc, doc, false)
				contentArray = append(contentArray, ndoc)
			}
			opts := options.InsertMany()
			opts.SetOrdered(false) // ignore duplication errors
			c.InsertMany(ctx, contentArray, opts)
		}(threadNum, num)
	}
	wg.Wait()

	if f.showProgress {
		fmt.Fprintf(os.Stderr, "\r100%%   \n")
	}
	cnt, _ := c.CountDocuments(ctx, bson.M{})
	fmt.Printf("\rSeeded %s: %d, total count: %d\n", collName, f.total, cnt)
	return err
}

func getEmployee(id int, supervisor int) bson.M {
	dealerID := "DEALER-1"
	email := util.GetEmailAddress()
	s := strings.Split(strings.Split(email, "@")[0], ".")
	doc := bson.M{"_id": int32(id), "dealer": dealerID, "email": email, "name": s[0] + " " + s[2]}
	if supervisor != 0 {
		doc["manager"] = int32(supervisor)
	}
	return doc
}
