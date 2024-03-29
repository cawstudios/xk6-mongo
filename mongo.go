package xk6_mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	k6modules "go.k6.io/k6/js/modules"
)

// Register the extension on module initialization, available to
// import from JS as "k6/x/mongo".
func init() {
	k6modules.Register("k6/x/mongo", new(Mongo))
}

// Mongo is the k6 extension for a Mongo client.
type Mongo struct{}

// Client is the Mongo client wrapper.
type Client struct {
	client *mongo.Client
}

// NewClient represents the Client constructor (i.e. `new mongo.Client()`) and
// returns a new Mongo client object.
// connURI -> mongodb://username:password@address:port/db?connect=direct
func (*Mongo) NewClient(connURI string) interface{} {

	clientOptions := options.Client().ApplyURI(connURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	return &Client{client: client}
}

const filter_is string = "filter is "

func (c *Client) Insert(database string, collection string, doc map[string]string) error {
	db := c.client.Database(database)
	col := db.Collection(collection)
	_, err := col.InsertOne(context.TODO(), doc)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) InsertMany(database string, collection string, docs []any) error {
	db := c.client.Database(database)
	col := db.Collection(collection)
	_, err := col.InsertMany(context.TODO(), docs)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Find(database string, collection string, filter interface{}, limit int64) []bson.M {
	db := c.client.Database(database)
	col := db.Collection(collection)
	if limit > 0 {
		opts := options.Find().SetLimit(limit)
		cur, err := col.Find(context.TODO(), filter, opts)
		if err != nil {
			log.Fatal(err)
		}
		var results []bson.M
		if err = cur.All(context.TODO(), &results); err != nil {
			panic(err)
		}
		return results
	} else {
		cur, err := col.Find(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}
		var results []bson.M
		if err = cur.All(context.TODO(), &results); err != nil {
			panic(err)
		}
		return results
	}
}

func (c *Client) FindOne(database string, collection string, filter map[string]string) bson.M {
	db := c.client.Database(database)
	col := db.Collection(collection)
	var result bson.M
	opts := options.FindOne().SetSort(bson.D{{"_id", 1}})
	err := col.FindOne(context.TODO(), filter, opts).Decode(&result)
	if err == mongo.ErrNoDocuments {
		log.Printf("No document was found for filter %v", filter)
		return nil
	}
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func (c *Client) UpdateOne(database string, collection string, filter interface{}, data map[string]string) error {
	// var result bson.M
	db := c.client.Database(database)
	col := db.Collection(collection)
	update := bson.D{{"$set", data}}
	_, err := col.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	// opts := options.FindOne().SetSort(bson.D{{"_id", 1}})
	// err = col.FindOne(context.TODO(), filter, opts).Decode(&result)
	// if err == mongo.ErrNoDocuments {
	// 	log.Printf("No document was found for filter %v", filter)
	// 	return nil
	// }
	return nil
}

func (c *Client) FindAll(database string, collection string) []bson.M {
	db := c.client.Database(database)
	col := db.Collection(collection)
	cur, err := col.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	var results []bson.M
	if err = cur.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	return results
}

func (c *Client) DeleteOne(database string, collection string, filter map[string]string) error {
	db := c.client.Database(database)
	col := db.Collection(collection)
	opts := options.Delete().SetHint(bson.D{{"_id", 1}})
	_, err := col.DeleteOne(context.TODO(), filter, opts)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (c *Client) DeleteMany(database string, collection string, filter map[string]string) error {
	db := c.client.Database(database)
	col := db.Collection(collection)
	opts := options.Delete().SetHint(bson.D{{"_id", 1}})
	_, err := col.DeleteMany(context.TODO(), filter, opts)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (c *Client) DropCollection(database string, collection string) error {
	db := c.client.Database(database)
	col := db.Collection(collection)
	err := col.Drop(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
