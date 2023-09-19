package mongo

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Connection setup functions
func createConnString(username string, password string, address string, port int) string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%d/", username, password, address, port)
}

func createClient(connectionString string, timeout *context.Context) (*mongo.Client, error) {
	opts := options.Client()
	opts.ApplyURI(connectionString)
	opts.SetMaxPoolSize(100)
	opts.SetMinPoolSize(10)
	client, err := mongo.Connect(*timeout, opts)
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
