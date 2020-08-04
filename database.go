package qmgo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type Database struct {
	database *mongo.Database
}

// NewDatabase creates new database connection
func NewDatabase(ctx context.Context, conf *Config) (cli *Database, err error) {
	client, err := client(ctx, conf)
	if err != nil {
		fmt.Println("new client fail", err)
		return
	}
	db := client.Database(conf.Database)

	cli = &Database{
		database: db,
	}
	return
}

// Collection gets collection from database
func (d *Database) Collection(name string) *Collection {
	var cp *mongo.Collection

	cp = d.database.Collection(name)

	return &Collection{
		collection: cp,
	}
}

// GetDatabaseName returns the name of database
func (d *Database) GetDatabaseName() string {
	return d.database.Name()
}

// DropDatabase drops database
func (d *Database) DropDatabase(ctx context.Context) {
	d.database.Drop(ctx)
}
