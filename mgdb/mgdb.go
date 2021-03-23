package mgdb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB - Db type.
type DB struct {
	DbURI  string
	Client *mongo.Client
	Db     *mongo.Database
}

// New - Initialise a new SkDB object. Connection to database.
func New(ctx context.Context, dbURI string, dbName string) (db DB, err error) {
	// Ref: https://docs.microsoft.com/en-us/azure/cosmos-db/create-mongodb-go
	// Using the SetDirect(true) configuration is important, without which you will get the following connectivity error:
	// "unable to connect connection(cdb-ms-prod-<azure-region>-cm1.documents.azure.com:10255[-4]) connection is closed"
	clientOptions := options.Client().ApplyURI(dbURI).SetDirect(true)
	client, err := mongo.NewClient(clientOptions)
	err = client.Connect(ctx)
	if err != nil {
		return
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return
	}

	database := client.Database(dbName)

	db.DbURI = dbURI
	db.Client = client
	db.Db = database

	return
}

// Close - Close db connection. Performs database disconnect.
func (db *DB) Close(ctx context.Context) {
	db.Client.Disconnect(ctx)
}

// DocExists - Returns true if doc exist. False if not.
func DocExists(collection *mongo.Collection, filter interface{}) (exist bool, err error) {
	findOptions := options.Find()
	findOptions.SetLimit(1)

	c := context.TODO()
	cur, err := collection.Find(c, filter, findOptions)
	if err != nil {
		return
	}
	defer cur.Close(c)

	exist = cur.Next(c) // Will return true if there's doc, false otherwise.
	err = nil
	return
}
