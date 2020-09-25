/*
 Copyright 2020 The Qmgo Authors.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

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
