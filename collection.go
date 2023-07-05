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
	"fmt"
	"reflect"
	"strings"

	"github.com/qiniu/qmgo/middleware"
	"github.com/qiniu/qmgo/operator"
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collection is a handle to a MongoDB collection
type Collection struct {
	collection *mongo.Collection

	registry *bsoncodec.Registry
}

// Find find by condition filter，return QueryI
func (c *Collection) Find(ctx context.Context, filter interface{}, opts ...opts.FindOptions) QueryI {

	return &Query{
		ctx:        ctx,
		collection: c.collection,
		filter:     filter,
		opts:       opts,
		registry:   c.registry,
	}
}

// InsertOne insert one document into the collection
// If InsertHook in opts is set, hook works on it, otherwise hook try the doc as hook
// Reference: https://docs.mongodb.com/manual/reference/command/insert/
func (c *Collection) InsertOne(ctx context.Context, doc interface{}, opts ...opts.InsertOneOptions) (result *InsertOneResult, err error) {
	h := doc
	insertOneOpts := options.InsertOne()
	if len(opts) > 0 {
		if opts[0].InsertOneOptions != nil {
			insertOneOpts = opts[0].InsertOneOptions
		}
		if opts[0].InsertHook != nil {
			h = opts[0].InsertHook
		}
	}
	if err = middleware.Do(ctx, doc, operator.BeforeInsert, h); err != nil {
		return
	}
	res, err := c.collection.InsertOne(ctx, doc, insertOneOpts)
	if res != nil {
		result = &InsertOneResult{InsertedID: res.InsertedID}
	}
	if err != nil {
		return
	}
	if err = middleware.Do(ctx, doc, operator.AfterInsert, h); err != nil {
		return
	}
	return
}

// InsertMany executes an insert command to insert multiple documents into the collection.
// If InsertHook in opts is set, hook works on it, otherwise hook try the doc as hook
// Reference: https://docs.mongodb.com/manual/reference/command/insert/
func (c *Collection) InsertMany(ctx context.Context, docs interface{}, opts ...opts.InsertManyOptions) (result *InsertManyResult, err error) {
	h := docs
	insertManyOpts := options.InsertMany()
	if len(opts) > 0 {
		if opts[0].InsertManyOptions != nil {
			insertManyOpts = opts[0].InsertManyOptions
		}
		if opts[0].InsertHook != nil {
			h = opts[0].InsertHook
		}
	}
	if err = middleware.Do(ctx, docs, operator.BeforeInsert, h); err != nil {
		return
	}
	sDocs := interfaceToSliceInterface(docs)
	if sDocs == nil {
		return nil, ErrNotValidSliceToInsert
	}

	res, err := c.collection.InsertMany(ctx, sDocs, insertManyOpts)
	if res != nil {
		result = &InsertManyResult{InsertedIDs: res.InsertedIDs}
	}
	if err != nil {
		return
	}
	if err = middleware.Do(ctx, docs, operator.AfterInsert, h); err != nil {
		return
	}
	return
}

// interfaceToSliceInterface convert interface to slice interface
func interfaceToSliceInterface(docs interface{}) []interface{} {
	if reflect.Slice != reflect.TypeOf(docs).Kind() {
		return nil
	}
	s := reflect.ValueOf(docs)
	if s.Len() == 0 {
		return nil
	}
	var sDocs []interface{}
	for i := 0; i < s.Len(); i++ {
		sDocs = append(sDocs, s.Index(i).Interface())
	}
	return sDocs
}

// Upsert updates one documents if filter match, inserts one document if filter is not match, Error when the filter is invalid
// The replacement parameter must be a document that will be used to replace the selected document. It cannot be nil
// and cannot contain any update operators
// Reference: https://docs.mongodb.com/manual/reference/operator/update/
// If replacement has "_id" field and the document is existed, please initial it with existing id(even with Qmgo default field feature).
// Otherwise, "the (immutable) field '_id' altered" error happens.
func (c *Collection) Upsert(ctx context.Context, filter interface{}, replacement interface{}, opts ...opts.UpsertOptions) (result *UpdateResult, err error) {
	h := replacement
	officialOpts := options.Replace().SetUpsert(true)

	if len(opts) > 0 {
		if opts[0].ReplaceOptions != nil {
			opts[0].ReplaceOptions.SetUpsert(true)
			officialOpts = opts[0].ReplaceOptions
		}
		if opts[0].UpsertHook != nil {
			h = opts[0].UpsertHook
		}
	}
	if err = middleware.Do(ctx, replacement, operator.BeforeUpsert, h); err != nil {
		return
	}

	res, err := c.collection.ReplaceOne(ctx, filter, replacement, officialOpts)

	if res != nil {
		result = translateUpdateResult(res)
	}
	if err != nil {
		return
	}
	if err = middleware.Do(ctx, replacement, operator.AfterUpsert, h); err != nil {
		return
	}
	return
}

// UpsertId updates one documents if id match, inserts one document if id is not match and the id will inject into the document
// The replacement parameter must be a document that will be used to replace the selected document. It cannot be nil
// and cannot contain any update operators
// Reference: https://docs.mongodb.com/manual/reference/operator/update/
func (c *Collection) UpsertId(ctx context.Context, id interface{}, replacement interface{}, opts ...opts.UpsertOptions) (result *UpdateResult, err error) {
	h := replacement
	officialOpts := options.Replace().SetUpsert(true)

	if len(opts) > 0 {
		if opts[0].ReplaceOptions != nil {
			opts[0].ReplaceOptions.SetUpsert(true)
			officialOpts = opts[0].ReplaceOptions
		}
		if opts[0].UpsertHook != nil {
			h = opts[0].UpsertHook
		}
	}
	if err = middleware.Do(ctx, replacement, operator.BeforeUpsert, h); err != nil {
		return
	}
	res, err := c.collection.ReplaceOne(ctx, bson.M{"_id": id}, replacement, officialOpts)
	if res != nil {
		result = translateUpdateResult(res)
	}
	if err != nil {
		return
	}
	if err = middleware.Do(ctx, replacement, operator.AfterUpsert, h); err != nil {
		return
	}
	return
}

// UpdateOne executes an update command to update at most one document in the collection.
// Reference: https://docs.mongodb.com/manual/reference/operator/update/
func (c *Collection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...opts.UpdateOptions) (err error) {
	updateOpts := options.Update()

	if len(opts) > 0 {
		if opts[0].UpdateOptions != nil {
			updateOpts = opts[0].UpdateOptions
		}
		if opts[0].UpdateHook != nil {
			if err = middleware.Do(ctx, opts[0].UpdateHook, operator.BeforeUpdate); err != nil {
				return
			}
		}
	}

	res, err := c.collection.UpdateOne(ctx, filter, update, updateOpts)
	if res != nil && res.MatchedCount == 0 {
		// UpdateOne support upsert function
		if updateOpts.Upsert == nil || !*updateOpts.Upsert {
			err = ErrNoSuchDocuments
		}
	}
	if err != nil {
		return err
	}
	if len(opts) > 0 && opts[0].UpdateHook != nil {
		if err = middleware.Do(ctx, opts[0].UpdateHook, operator.AfterUpdate); err != nil {
			return
		}
	}
	return err
}

// UpdateId executes an update command to update at most one document in the collection.
// Reference: https://docs.mongodb.com/manual/reference/operator/update/
func (c *Collection) UpdateId(ctx context.Context, id interface{}, update interface{}, opts ...opts.UpdateOptions) (err error) {
	updateOpts := options.Update()

	if len(opts) > 0 {
		if opts[0].UpdateOptions != nil {
			updateOpts = opts[0].UpdateOptions
		}
		if opts[0].UpdateHook != nil {
			if err = middleware.Do(ctx, opts[0].UpdateHook, operator.BeforeUpdate); err != nil {
				return
			}
		}
	}

	res, err := c.collection.UpdateOne(ctx, bson.M{"_id": id}, update, updateOpts)
	if res != nil && res.MatchedCount == 0 {
		err = ErrNoSuchDocuments
	}
	if err != nil {
		return err
	}
	if len(opts) > 0 && opts[0].UpdateHook != nil {
		if err = middleware.Do(ctx, opts[0].UpdateHook, operator.AfterUpdate); err != nil {
			return
		}
	}
	return err
}

// UpdateAll executes an update command to update documents in the collection.
// The matchedCount is 0 in UpdateResult if no document updated
// Reference: https://docs.mongodb.com/manual/reference/operator/update/
func (c *Collection) UpdateAll(ctx context.Context, filter interface{}, update interface{}, opts ...opts.UpdateOptions) (result *UpdateResult, err error) {
	updateOpts := options.Update()
	if len(opts) > 0 {
		if opts[0].UpdateOptions != nil {
			updateOpts = opts[0].UpdateOptions
		}
		if opts[0].UpdateHook != nil {
			if err = middleware.Do(ctx, opts[0].UpdateHook, operator.BeforeUpdate); err != nil {
				return
			}
		}
	}
	res, err := c.collection.UpdateMany(ctx, filter, update, updateOpts)
	if res != nil {
		result = translateUpdateResult(res)
	}
	if err != nil {
		return
	}
	if len(opts) > 0 && opts[0].UpdateHook != nil {
		if err = middleware.Do(ctx, opts[0].UpdateHook, operator.AfterUpdate); err != nil {
			return
		}
	}
	return
}

// ReplaceOne executes an update command to update at most one document in the collection.
// If UpdateHook in opts is set, hook works on it, otherwise hook try the doc as hook
// Expect type of the doc is the define of user's document
func (c *Collection) ReplaceOne(ctx context.Context, filter interface{}, doc interface{}, opts ...opts.ReplaceOptions) (err error) {
	h := doc
	replaceOpts := options.Replace()

	if len(opts) > 0 {
		if opts[0].ReplaceOptions != nil {
			replaceOpts = opts[0].ReplaceOptions
			replaceOpts.SetUpsert(false)
		}
		if opts[0].UpdateHook != nil {
			h = opts[0].UpdateHook
		}
	}
	if err = middleware.Do(ctx, doc, operator.BeforeReplace, h); err != nil {
		return
	}
	res, err := c.collection.ReplaceOne(ctx, filter, doc, replaceOpts)
	if res != nil && res.MatchedCount == 0 {
		err = ErrNoSuchDocuments
	}
	if err != nil {
		return err
	}
	if err = middleware.Do(ctx, doc, operator.AfterReplace, h); err != nil {
		return
	}

	return err
}

// Remove executes a delete command to delete at most one document from the collection.
// if filter is bson.M{}，DeleteOne will delete one document in collection
// Reference: https://docs.mongodb.com/manual/reference/command/delete/
func (c *Collection) Remove(ctx context.Context, filter interface{}, opts ...opts.RemoveOptions) (err error) {
	deleteOptions := options.Delete()
	if len(opts) > 0 {
		if opts[0].DeleteOptions != nil {
			deleteOptions = opts[0].DeleteOptions
		}
		if opts[0].RemoveHook != nil {
			if err = middleware.Do(ctx, opts[0].RemoveHook, operator.BeforeRemove); err != nil {
				return err
			}
		}
	}
	res, err := c.collection.DeleteOne(ctx, filter, deleteOptions)
	if res != nil && res.DeletedCount == 0 {
		err = ErrNoSuchDocuments
	}
	if err != nil {
		return err
	}
	if len(opts) > 0 && opts[0].RemoveHook != nil {
		if err = middleware.Do(ctx, opts[0].RemoveHook, operator.AfterRemove); err != nil {
			return err
		}
	}
	return err
}

// RemoveId executes a delete command to delete at most one document from the collection.
func (c *Collection) RemoveId(ctx context.Context, id interface{}, opts ...opts.RemoveOptions) (err error) {
	deleteOptions := options.Delete()
	if len(opts) > 0 {
		if opts[0].DeleteOptions != nil {
			deleteOptions = opts[0].DeleteOptions
		}
		if opts[0].RemoveHook != nil {
			if err = middleware.Do(ctx, opts[0].RemoveHook, operator.BeforeRemove); err != nil {
				return err
			}
		}
	}
	res, err := c.collection.DeleteOne(ctx, bson.M{"_id": id}, deleteOptions)
	if res != nil && res.DeletedCount == 0 {
		err = ErrNoSuchDocuments
	}
	if err != nil {
		return err
	}

	if len(opts) > 0 && opts[0].RemoveHook != nil {
		if err = middleware.Do(ctx, opts[0].RemoveHook, operator.AfterRemove); err != nil {
			return err
		}
	}
	return err
}

// RemoveAll executes a delete command to delete documents from the collection.
// If filter is bson.M{}，all ducuments in Collection will be deleted
// Reference: https://docs.mongodb.com/manual/reference/command/delete/
func (c *Collection) RemoveAll(ctx context.Context, filter interface{}, opts ...opts.RemoveOptions) (result *DeleteResult, err error) {
	deleteOptions := options.Delete()
	if len(opts) > 0 {
		if opts[0].DeleteOptions != nil {
			deleteOptions = opts[0].DeleteOptions
		}
		if opts[0].RemoveHook != nil {
			if err = middleware.Do(ctx, opts[0].RemoveHook, operator.BeforeRemove); err != nil {
				return
			}
		}
	}
	res, err := c.collection.DeleteMany(ctx, filter, deleteOptions)
	if res != nil {
		result = &DeleteResult{DeletedCount: res.DeletedCount}
	}
	if err != nil {
		return
	}
	if len(opts) > 0 && opts[0].RemoveHook != nil {
		if err = middleware.Do(ctx, opts[0].RemoveHook, operator.AfterRemove); err != nil {
			return
		}
	}
	return
}

// Aggregate executes an aggregate command against the collection and returns a AggregateI to get resulting documents.
func (c *Collection) Aggregate(ctx context.Context, pipeline interface{}, opts ...opts.AggregateOptions) AggregateI {
	return &Aggregate{
		ctx:        ctx,
		collection: c.collection,
		pipeline:   pipeline,
		options:    opts,
	}
}

// ensureIndex create multiple indexes on the collection and returns the names of
// Example：indexes = []string{"idx1", "-idx2", "idx3,idx4"}
// Three indexes will be created, index idx1 with ascending order, index idx2 with descending order, idex3 and idex4 are Compound ascending sort index
// Reference: https://docs.mongodb.com/manual/reference/command/createIndexes/
func (c *Collection) ensureIndex(ctx context.Context, indexes []opts.IndexModel) error {
	var indexModels []mongo.IndexModel
	for _, idx := range indexes {
		var model mongo.IndexModel
		var keysDoc bson.D

		for _, field := range idx.Key {
			key, n := SplitSortField(field)

			keysDoc = append(keysDoc, bson.E{Key: key, Value: n})
		}
		model = mongo.IndexModel{
			Keys:    keysDoc,
			Options: idx.IndexOptions,
		}

		indexModels = append(indexModels, model)
	}

	if len(indexModels) == 0 {
		return nil
	}

	res, err := c.collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil || len(res) == 0 {
		fmt.Println("<MongoDB.C>: ", c.collection.Name(), " Index: ", indexes, " error: ", err, "res: ", res)
		return err
	}
	return nil
}

// EnsureIndexes Deprecated
// Recommend to use CreateIndexes / CreateOneIndex for more function)
// EnsureIndexes creates unique and non-unique indexes in collection
// the combination of indexes is different from CreateIndexes:
// if uniques/indexes is []string{"name"}, means create index "name"
// if uniques/indexes is []string{"name,-age","uid"} means create Compound indexes: name and -age, then create one index: uid
func (c *Collection) EnsureIndexes(ctx context.Context, uniques []string, indexes []string) (err error) {
	var uniqueModel []opts.IndexModel
	var indexesModel []opts.IndexModel
	for _, v := range uniques {
		vv := strings.Split(v, ",")
		indexOpts := options.Index()
		indexOpts.SetUnique(true)
		model := opts.IndexModel{Key: vv, IndexOptions: indexOpts}
		uniqueModel = append(uniqueModel, model)
	}
	if err = c.CreateIndexes(ctx, uniqueModel); err != nil {
		return
	}

	for _, v := range indexes {
		vv := strings.Split(v, ",")
		model := opts.IndexModel{Key: vv}
		indexesModel = append(indexesModel, model)
	}
	if err = c.CreateIndexes(ctx, indexesModel); err != nil {
		return
	}
	return
}

// CreateIndexes creates multiple indexes in collection
// If the Key in opts.IndexModel is []string{"name"}, means create index: name
// If the Key in opts.IndexModel is []string{"name","-age"} means create Compound indexes: name and -age
func (c *Collection) CreateIndexes(ctx context.Context, indexes []opts.IndexModel) (err error) {
	err = c.ensureIndex(ctx, indexes)
	return
}

// CreateOneIndex creates one index
// If the Key in opts.IndexModel is []string{"name"}, means create index name
// If the Key in opts.IndexModel is []string{"name","-age"} means create Compound index: name and -age
func (c *Collection) CreateOneIndex(ctx context.Context, index opts.IndexModel) error {
	return c.ensureIndex(ctx, []opts.IndexModel{index})

}

// DropAllIndexes drop all indexes on the collection except the index on the _id field
// if there is only _id field index on the collection, the function call will report an error
func (c *Collection) DropAllIndexes(ctx context.Context) (err error) {
	_, err = c.collection.Indexes().DropAll(ctx)
	return err
}

// DropIndex drop indexes in collection, indexes that be dropped should be in line with inputting indexes
// The indexes is []string{"name"} means drop index: name
// The indexes is []string{"name","-age"} means drop Compound indexes: name and -age
func (c *Collection) DropIndex(ctx context.Context, indexes []string) error {
	_, err := c.collection.Indexes().DropOne(ctx, generateDroppedIndex(indexes))
	if err != nil {
		return err
	}
	return err
}

// generate indexes that store in mongo which may consist more than one index(like []string{"index1","index2"} is stored as "index1_1_index2_1")
func generateDroppedIndex(index []string) string {
	var res string
	for _, e := range index {
		key, sort := SplitSortField(e)
		n := key + "_" + fmt.Sprint(sort)
		if len(res) == 0 {
			res = n
		} else {
			res += "_" + n
		}
	}
	return res
}

// DropCollection drops collection
// it's safe even collection is not exists
func (c *Collection) DropCollection(ctx context.Context) error {
	return c.collection.Drop(ctx)
}

// CloneCollection creates a copy of the Collection
func (c *Collection) CloneCollection() (*mongo.Collection, error) {
	return c.collection.Clone()
}

// GetCollectionName returns the name of collection
func (c *Collection) GetCollectionName() string {
	return c.collection.Name()
}

// Watch returns a change stream for all changes on the corresponding collection. See
// https://docs.mongodb.com/manual/changeStreams/ for more information about change streams.
func (c *Collection) Watch(ctx context.Context, pipeline interface{}, opts ...*opts.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	changeStreamOption := options.ChangeStream()
	if len(opts) > 0 && opts[0].ChangeStreamOptions != nil {
		changeStreamOption = opts[0].ChangeStreamOptions
	}
	return c.collection.Watch(ctx, pipeline, changeStreamOption)
}

// translateUpdateResult translates mongo update result to qmgo define UpdateResult
func translateUpdateResult(res *mongo.UpdateResult) (result *UpdateResult) {
	result = &UpdateResult{
		MatchedCount:  res.MatchedCount,
		ModifiedCount: res.ModifiedCount,
		UpsertedCount: res.UpsertedCount,
		UpsertedID:    res.UpsertedID,
	}
	return
}
