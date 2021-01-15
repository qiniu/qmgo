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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BulkResult is the result type returned by Bulk.Run operation.
type BulkResult struct {
	// The number of documents inserted.
	InsertedCount int64

	// The number of documents matched by filters in update and replace operations.
	MatchedCount int64

	// The number of documents modified by update and replace operations.
	ModifiedCount int64

	// The number of documents deleted.
	DeletedCount int64

	// The number of documents upserted by update and replace operations.
	UpsertedCount int64

	// A map of operation index to the _id of each upserted document.
	UpsertedIDs map[int64]interface{}
}

// Bulk is context for batching operations to be sent to database in a single
// bulk write.
//
// Bulk is not safe for concurrent use.
//
// Notes:
//
// Individual operations inside a bulk do not trigger middlewares or hooks
// at present.
//
// Different from original mgo, the qmgo implementation of Bulk does not emulate
// bulk operations individually on old versions of MongoDB servers that do not
// natively support bulk operations.
//
// Only operations supported by the official driver are exposed, that is why
// InsertMany is missing from the methods.
type Bulk struct {
	coll *Collection

	queue   []mongo.WriteModel
	ordered *bool
}

// Bulk returns a new context for preparing bulk execution of operations.
func (c *Collection) Bulk() *Bulk {
	return &Bulk{
		coll:    c,
		queue:   nil,
		ordered: nil,
	}
}

// SetOrdered marks the bulk as ordered or unordered.
//
// If ordered, writes does not continue after one individual write fails.
// Default is ordered.
func (b *Bulk) SetOrdered(ordered bool) {
	b.ordered = &ordered
}

// InsertOne queues an InsertOne operation for bulk execution.
func (b *Bulk) InsertOne(doc interface{}) {
	wm := &mongo.InsertOneModel{
		Document: doc,
	}
	b.queue = append(b.queue, wm)
}

// Remove queues a Remove operation for bulk execution.
func (b *Bulk) Remove(filter interface{}) {
	wm := &mongo.DeleteOneModel{
		Filter: filter,
	}
	b.queue = append(b.queue, wm)
}

// RemoveId queues a RemoveId operation for bulk execution.
func (b *Bulk) RemoveId(id interface{}) {
	b.Remove(bson.M{"_id": id})
}

// RemoveAll queues a RemoveAll operation for bulk execution.
func (b *Bulk) RemoveAll(filter interface{}) {
	wm := &mongo.DeleteManyModel{
		Filter: filter,
	}
	b.queue = append(b.queue, wm)
}

// Upsert queues an Upsert operation for bulk execution.
func (b *Bulk) Upsert(filter interface{}, replacement interface{}) {
	wm := &mongo.UpdateOneModel{
		Filter: filter,
		Update: replacement,
	}
	wm.SetUpsert(true)
	b.queue = append(b.queue, wm)
}

// UpsertId queues an UpsertId operation for bulk execution.
func (b *Bulk) UpsertId(id interface{}, replacement interface{}) {
	b.Upsert(bson.M{"_id": id}, replacement)
}

// UpdateOne queues an UpdateOne operation for bulk execution.
func (b *Bulk) UpdateOne(filter interface{}, update interface{}) {
	wm := &mongo.UpdateOneModel{
		Filter: filter,
		Update: update,
	}
	b.queue = append(b.queue, wm)
}

// UpdateId queues an UpdateId operation for bulk execution.
func (b *Bulk) UpdateId(id interface{}, update interface{}) {
	b.UpdateOne(bson.M{"_id": id}, update)
}

// UpdateAll queues an UpdateAll operation for bulk execution.
func (b *Bulk) UpdateAll(filter interface{}, update interface{}) {
	wm := &mongo.UpdateManyModel{
		Filter: filter,
		Update: update,
	}
	b.queue = append(b.queue, wm)
}

// Run executes the collected operations in a single bulk operation.
//
// A successful call resets the Bulk. If an error is returned, the internal
// queue of operations is unchanged, containing both successful and failed
// operations.
func (b *Bulk) Run(ctx context.Context) (*BulkResult, error) {
	opts := options.BulkWriteOptions{
		Ordered: b.ordered,
	}
	result, err := b.coll.collection.BulkWrite(ctx, b.queue, &opts)
	if err != nil {
		// In original mgo, queue is not reset in case of error.
		return nil, err
	}

	// Empty the queue for possible reuse, as per mgo's behavior.
	b.queue = nil

	return &BulkResult{
		InsertedCount: result.InsertedCount,
		MatchedCount:  result.MatchedCount,
		ModifiedCount: result.ModifiedCount,
		DeletedCount:  result.DeletedCount,
		UpsertedCount: result.UpsertedCount,
		UpsertedIDs:   result.UpsertedIDs,
	}, nil
}
