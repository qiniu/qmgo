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
func (b *Bulk) SetOrdered(ordered bool) *Bulk {
	b.ordered = &ordered
	return b
}

// InsertOne queues an InsertOne operation for bulk execution.
func (b *Bulk) InsertOne(doc interface{}) *Bulk {
	wm := mongo.NewInsertOneModel().SetDocument(doc)
	b.queue = append(b.queue, wm)
	return b
}

// Remove queues a Remove operation for bulk execution.
func (b *Bulk) Remove(filter interface{}) *Bulk {
	wm := mongo.NewDeleteOneModel().SetFilter(filter)
	b.queue = append(b.queue, wm)
	return b
}

// RemoveId queues a RemoveId operation for bulk execution.
func (b *Bulk) RemoveId(id interface{}) *Bulk {
	b.Remove(bson.M{"_id": id})
	return b
}

// RemoveAll queues a RemoveAll operation for bulk execution.
func (b *Bulk) RemoveAll(filter interface{}) *Bulk {
	wm := mongo.NewDeleteManyModel().SetFilter(filter)
	b.queue = append(b.queue, wm)
	return b
}

// Upsert queues an Upsert operation for bulk execution.
// The replacement should be document without operator
func (b *Bulk) Upsert(filter interface{}, replacement interface{}) *Bulk {
	wm := mongo.NewReplaceOneModel().SetFilter(filter).SetReplacement(replacement).SetUpsert(true)
	b.queue = append(b.queue, wm)
	return b
}

// UpsertOne queues an UpsertOne operation for bulk execution.
// The update should contain operator
func (b *Bulk) UpsertOne(filter interface{}, update interface{}) *Bulk {
	wm := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
	b.queue = append(b.queue, wm)
	return b
}

// UpsertId queues an UpsertId operation for bulk execution.
// The replacement should be document without operator
func (b *Bulk) UpsertId(id interface{}, replacement interface{}) *Bulk {
	b.Upsert(bson.M{"_id": id}, replacement)
	return b
}

// UpdateOne queues an UpdateOne operation for bulk execution.
// The update should contain operator
func (b *Bulk) UpdateOne(filter interface{}, update interface{}) *Bulk {
	wm := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update)
	b.queue = append(b.queue, wm)
	return b
}

// UpdateId queues an UpdateId operation for bulk execution.
// The update should contain operator
func (b *Bulk) UpdateId(id interface{}, update interface{}) *Bulk {
	b.UpdateOne(bson.M{"_id": id}, update)
	return b
}

// UpdateAll queues an UpdateAll operation for bulk execution.
// The update should contain operator
func (b *Bulk) UpdateAll(filter interface{}, update interface{}) *Bulk {
	wm := mongo.NewUpdateManyModel().SetFilter(filter).SetUpdate(update)
	b.queue = append(b.queue, wm)
	return b
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
