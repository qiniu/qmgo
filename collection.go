package qmgo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// Collection is a handle to a MongoDB collection
type Collection struct {
	collection *mongo.Collection
}

// Find find by condition filter，return QueryI
func (c *Collection) Find(ctx context.Context, filter interface{}) QueryI {
	return &Query{
		ctx:        ctx,
		collection: c.collection,
		filter:     filter,
	}
}

// InsertOne insert one document into the collection
// Reference: https://docs.mongodb.com/manual/reference/command/insert/
func (c *Collection) InsertOne(ctx context.Context, doc interface{}) (result *InsertOneResult, err error) {
	res, err := c.collection.InsertOne(ctx, doc)
	if res != nil {
		result = &InsertOneResult{InsertedID: res.InsertedID}
	}
	return
}

// InsertMany executes an insert command to insert multiple documents into the collection.
// e.g. docs := []interface{}{myDocsInstance1, myDocsInstance2}
// TODO need a function which translate slice to []interface
// Reference: https://docs.mongodb.com/manual/reference/command/insert/
func (c *Collection) InsertMany(ctx context.Context, docs []interface{}) (result *InsertManyResult, err error) {
	res, err := c.collection.InsertMany(ctx, docs)
	if res != nil {
		result = &InsertManyResult{InsertedIDs: res.InsertedIDs}
	}
	return
}

// Upsert updates one documents if filter match, inserts one document if filter is not match
// Reference: https://docs.mongodb.com/manual/reference/operator/update/
func (c *Collection) Upsert(ctx context.Context, filter interface{}, replacement interface{}) (result *UpdateResult, err error) {
	opts := options.Replace().SetUpsert(true)
	res, err := c.collection.ReplaceOne(ctx, filter, replacement, opts)
	if res != nil {
		result = translateUpdateResult(res)
	}
	return
}

// UpdateOne executes an update command to update at most one document in the collection.
// Reference: https://docs.mongodb.com/manual/reference/operator/update/
func (c *Collection) UpdateOne(ctx context.Context, filter interface{}, update interface{}) error {
	var err error
	var res *mongo.UpdateResult

	if res, err = c.collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		err = ErrNoSuchDocuments
	}

	return err
}

// UpdateAll executes an update command to update documents in the collection.
// The matchedCount is 0 in UpdateResult if no document updated
// Reference: https://docs.mongodb.com/manual/reference/operator/update/
func (c *Collection) UpdateAll(ctx context.Context, filter interface{}, update interface{}) (result *UpdateResult, err error) {

	res, err := c.collection.UpdateMany(ctx, filter, update)
	if res != nil {
		result = translateUpdateResult(res)
	}
	return
}

// Remove executes a delete command to delete at most one document from the collection.
// if filter is bson.M{}，DeleteOne will delete one document in collection
// Reference: https://docs.mongodb.com/manual/reference/command/delete/
func (c *Collection) Remove(ctx context.Context, filter interface{}) (err error) {

	res, err := c.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		err = ErrNoSuchDocuments
	}

	return err
}

// RemoveId executes a delete command to delete at most one document from the collection.
func (c *Collection) RemoveId(ctx context.Context, id string) (err error) {
	res, err := c.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		err = ErrNoSuchDocuments
	}

	return err
}

// DeleteAll executes a delete command to delete documents from the collection.
// If filter is bson.M{}，all ducuments in Collection will be deleted
// Reference: https://docs.mongodb.com/manual/reference/command/delete/
func (c *Collection) DeleteAll(ctx context.Context, filter interface{}) (result *DeleteResult, err error) {

	res, err := c.collection.DeleteMany(ctx, filter)
	if res != nil {
		result = &DeleteResult{DeletedCount: res.DeletedCount}
	}
	return
}

// Aggregate executes an aggregate command against the collection and returns a AggregateI to get resulting documents.
func (c *Collection) Aggregate(ctx context.Context, pipeline interface{}) AggregateI {
	return &Aggregate{
		ctx:        ctx,
		collection: c.collection,
		pipeline:   pipeline,
	}
}

// ensureIndex create multiple indexes on the collection and returns the names of
// Example：indexes = []string{"idx1", "-idx2", "idx3,idx4"}
// Three indexes will be created, index idx1 with ascending order, index idx2 with descending order, idex3 and idex4 are Compound ascending sort index
// Reference: https://docs.mongodb.com/manual/reference/command/createIndexes/
func (c *Collection) ensureIndex(ctx context.Context, indexes []string, isUnique bool) {
	var indexModels []mongo.IndexModel

	// 组建[]mongo.IndexModel
	for _, idx := range indexes {
		var model mongo.IndexModel
		var keysDoc bsonx.Doc

		colIndexArr := strings.Split(idx, ",")
		for _, field := range colIndexArr {
			key, n := SplitSortField(field)

			keysDoc = keysDoc.Append(key, bsonx.Int32(n))
		}

		model = mongo.IndexModel{
			Keys:    keysDoc,
			Options: options.Index().SetUnique(isUnique),
		}

		indexModels = append(indexModels, model)
	}

	if len(indexModels) == 0 {
		return
	}

	var err error
	var res []string
	res, err = c.collection.Indexes().CreateMany(ctx, indexModels)

	if err != nil || len(res) == 0 {
		s := fmt.Sprint("<MongoDB.C>: ", c.collection.Name(), " Index: ", indexes, " error: ", err, "res: ", res)
		panic(s)
	}
	return
}

// EnsureIndexes creates unique and non-unique indexes in collection
func (c *Collection) EnsureIndexes(ctx context.Context, uniques []string, indexes []string) {
	// 创建唯一索引
	if len(uniques) != 0 {
		c.ensureIndex(ctx, uniques, true)
	}

	// 创建普通索引
	if len(indexes) != 0 {
		c.ensureIndex(ctx, indexes, false)
	}

	return
}

// DropIndexes drop indexes in collection, indexes that be dropped should be in line with inputting indexes
func (c *Collection) DropIndexes(ctx context.Context, indexes []string) error {

	var err error
	for _, index := range indexes {
		_, err = c.collection.Indexes().DropOne(ctx, generateDroppedIndex(index))
		if err != nil {
			return err
		}
	}
	return err
}

// generate indexes that store in mongo which may consist more than one index(like "index1,index2" is stored as "index1_1_index2_1")
func generateDroppedIndex(index string) string {
	var res string
	s := strings.Split(index, ",")
	for _, e := range s {
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
