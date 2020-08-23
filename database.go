package qmgo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

// Database is a handle to a MongoDB database
type Database struct {
	database *mongo.Database
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
func (d *Database) DropDatabase(ctx context.Context) error {
	return d.database.Drop(ctx)
}
