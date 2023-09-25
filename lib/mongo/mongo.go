package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Connection setup functions
func createConnString(username string, password string, address string, port int) string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%d/", username, password, address, port)
}

func createClient(connectionString string, ctx *context.Context) (*mongo.Client, error) {
	opts := options.Client()
	opts.ApplyURI(connectionString)
	opts.SetMaxPoolSize(100)
	opts.SetMinPoolSize(10)
	opts.SetConnectTimeout(5 * time.Second)
	opts.SetTimeout(5 * time.Second)
	opts.SetMaxConnIdleTime(5 * time.Second)
	opts.SetSocketTimeout(5 * time.Second)
	client, err := mongo.Connect(*ctx, opts)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (r *Repository) getDatabases() ([]string, error) {
	databases, err := r.Db.ListDatabaseNames(*r.Ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	return databases, nil
}

func searchDb(databaseName string, database []string) string {
	for _, db := range database {
		if db == databaseName {
			return db
		}
		continue
	}
	return ""
}

func (r *Repository) isDbExist(databaseName string) (bool, error) {
	dbs, err := r.getDatabases()
	if err != nil {
		return false, err
	}
	db := searchDb(databaseName, dbs)
	if db == databaseName {
		log.Println("Database is exist.")
		r.Disconnect()
		return true, nil
	}
	log.Println("Database is not exist.")
	return false, nil
}

// Public functions
func NewRepository(client *mongo.Client, ctx *context.Context) *Repository {
	return &Repository{
		Db:  client,
		Ctx: ctx,
	}
}

func (r *Repository) Connect(username string, password string, address string, port int) (*mongo.Client, *context.Context, error) {
	connString := createConnString(username, password, address, port)
	ctx := context.Background()
	client, err := createClient(connString, &ctx)
	if err != nil {
		return nil, nil, err
	}
	return client, &ctx, nil
}

func (r *Repository) Disconnect() {
	r.Db.Disconnect(*r.Ctx)
}

func (r *Repository) PingDb() error {
	return r.Db.Ping(*r.Ctx, readpref.Primary())
}

func (r *Repository) CreateDb(dbName, collectionName string) error {
	isDbExist, _ := r.isDbExist(dbName)
	if !isDbExist {
		log.Println("Creating database...")
		err := r.Db.Database(dbName).CreateCollection(*r.Ctx, collectionName)
		if err != nil {
			return err
		}
		log.Println("Database created.")
	}
	return nil
}

func (r *Repository) GetAllDatabase() ([]string, error) {
	dbs, err := r.Db.ListDatabaseNames(*r.Ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	return dbs, nil
}

func (r *Repository) CreateDay(day *AdventCalendarDay, dbName, collectionName string) error {
	collection := r.Db.Database(dbName).Collection(collectionName)
	_, err := collection.InsertOne(*r.Ctx, day)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetDay(day uint8, dbName, collectionName string) (*AdventCalendarDay, error) {
	collection := r.Db.Database(dbName).Collection(collectionName)
	var result AdventCalendarDay
	err := collection.FindOne(*r.Ctx, bson.M{"day": day}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *Repository) GetAllDay(dbName, collectionName string) ([]*AdventCalendarDay, error) {
	collection := r.Db.Database(dbName).Collection(collectionName)
	var results []*AdventCalendarDay
	cur, err := collection.Find(*r.Ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	for cur.Next(*r.Ctx) {
		var elem AdventCalendarDay
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	cur.Close(*r.Ctx)
	if len(results) == 0 {
		return nil, mongo.ErrNoDocuments
	}
	return results, nil
}

func (r *Repository) UpdateDay(day *AdventCalendarDayUpdate, dbName, collectionName string) error {
	collection := r.Db.Database(dbName).Collection(collectionName)
	update := bson.M{"day": day.Day, "year": day.Year, "title": day.Title, "content": day.Content}
	_, err := collection.UpdateOne(*r.Ctx, bson.M{"_id": day.ID}, bson.M{"$set": update})
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteDay(day *DayIDRequest, dbName, collectionName string) error {
	collection := r.Db.Database(dbName).Collection(collectionName)
	_, err := collection.DeleteOne(*r.Ctx, bson.M{"_id": day.Id})
	if err != nil {
		return err
	}
	return nil
}
