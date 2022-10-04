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
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	officialOpts "go.mongodb.org/mongo-driver/mongo/options"

	"github.com/qiniu/qmgo/operator"
	"github.com/qiniu/qmgo/options"
)

func TestCollection_EnsureIndex(t *testing.T) {
	ast := require.New(t)
	cli := initClient("test")
	defer cli.Close(context.Background())
	defer cli.DropCollection(context.Background())

	cli.ensureIndex(context.Background(), nil)
	indexOpts := officialOpts.Index()
	indexOpts.SetUnique(true)
	cli.ensureIndex(context.Background(), []options.IndexModel{{Key: []string{"id1"}, IndexOptions: indexOpts}})
	cli.ensureIndex(context.Background(), []options.IndexModel{{Key: []string{"id2", "id3"}}})
	cli.ensureIndex(context.Background(), []options.IndexModel{{Key: []string{"id4", "-id5"}}})

	// same index，error
	ast.Error(cli.ensureIndex(context.Background(), []options.IndexModel{{Key: []string{"id1"}}}))

	// check if unique indexs is working
	var err error
	doc := bson.M{
		"id1": 1,
	}
	_, err = cli.InsertOne(context.Background(), doc)
	ast.NoError(err)

	coll, err := cli.CloneCollection()
	ast.NoError(err)
	_, err = coll.InsertOne(context.Background(), doc)
	ast.Equal(true, IsDup(err))
}

func TestCollection_EnsureIndexes(t *testing.T) {
	ast := require.New(t)
	cli := initClient("test")
	defer cli.Close(context.Background())
	defer cli.DropCollection(context.Background())

	unique := []string{"id1"}
	common := []string{"id2,id3", "id4,-id5"}
	cli.EnsureIndexes(context.Background(), unique, common)

	// same index，error
	ast.Error(cli.EnsureIndexes(context.Background(), nil, unique))

	// check if unique indexs is working
	var err error
	doc := bson.M{
		"id1": 1,
	}

	_, err = cli.InsertOne(context.Background(), doc)
	ast.NoError(err)
	_, err = cli.InsertOne(context.Background(), doc)
	ast.Equal(true, IsDup(err))
}

func TestCollection_CreateIndexes(t *testing.T) {
	ast := require.New(t)
	cli := initClient("test")
	defer cli.Close(context.Background())
	defer cli.DropCollection(context.Background())

	var expireS int32 = 100
	unique := []string{"id1"}
	indexOpts := officialOpts.Index()
	indexOpts.SetUnique(true).SetExpireAfterSeconds(expireS)
	ast.NoError(cli.CreateOneIndex(context.Background(), options.IndexModel{Key: unique, IndexOptions: indexOpts}))

	ast.NoError(cli.CreateIndexes(context.Background(), []options.IndexModel{{Key: []string{"id2", "id3"}},
		{Key: []string{"id4", "-id5"}}}))
	// same index，error
	ast.Error(cli.CreateOneIndex(context.Background(), options.IndexModel{Key: unique}))

	// check if unique indexs is working
	var err error
	doc := bson.M{
		"id1": 1,
	}

	_, err = cli.InsertOne(context.Background(), doc)
	ast.NoError(err)
	_, err = cli.InsertOne(context.Background(), doc)
	ast.Equal(true, IsDup(err))
}

func TestCollection_DropAllIndexes(t *testing.T) {
	ast := require.New(t)

	cli := initClient("test")
	defer cli.DropCollection(context.Background())

	var err error
	err = cli.DropAllIndexes(context.Background())
	ast.Error(err)

	unique := []string{"id1"}
	common := []string{"id2,id3", "id4,-id5"}
	cli.EnsureIndexes(context.Background(), unique, common)

	err = cli.DropAllIndexes(context.Background())
	ast.NoError(err)
}

func TestCollection_DropIndex(t *testing.T) {
	ast := require.New(t)

	cli := initClient("test")
	defer cli.DropCollection(context.Background())

	indexOpts := officialOpts.Index()
	indexOpts.SetUnique(true)
	cli.ensureIndex(context.Background(), []options.IndexModel{{Key: []string{"index1"}, IndexOptions: indexOpts}})

	// same index，error
	ast.Error(cli.ensureIndex(context.Background(), []options.IndexModel{{Key: []string{"index1"}}}))

	err := cli.DropIndex(context.Background(), []string{"index1"})
	ast.NoError(err)
	ast.NoError(cli.ensureIndex(context.Background(), []options.IndexModel{{Key: []string{"index1"}}}))

	indexOpts = officialOpts.Index()
	indexOpts.SetUnique(true)
	cli.ensureIndex(context.Background(), []options.IndexModel{{Key: []string{"-index1"}, IndexOptions: indexOpts}})
	// same index，error
	ast.Error(cli.ensureIndex(context.Background(), []options.IndexModel{{Key: []string{"-index1"}}}))

	err = cli.DropIndex(context.Background(), []string{"-index1"})
	ast.NoError(err)
	ast.NoError(cli.ensureIndex(context.Background(), []options.IndexModel{{Key: []string{"-index1"}}}))

	err = cli.DropIndex(context.Background(), []string{""})
	ast.Error(err)

	err = cli.DropIndex(context.Background(), []string{"index2"})
	ast.Error(err)

	indexOpts = officialOpts.Index()
	indexOpts.SetUnique(true)
	cli.ensureIndex(context.Background(), []options.IndexModel{{Key: []string{"index2", "-index1"}, IndexOptions: indexOpts}})
	ast.Error(cli.ensureIndex(context.Background(), []options.IndexModel{{Key: []string{"index2", "-index1"}}}))
	err = cli.DropIndex(context.Background(), []string{"index2", "-index1"})
	ast.NoError(err)
	ast.NoError(cli.ensureIndex(context.Background(), []options.IndexModel{{Key: []string{"index2", "-index1"}}}))

	err = cli.DropIndex(context.Background(), []string{"-"})
	ast.Error(err)
}

func TestCollection_Insert(t *testing.T) {
	ast := require.New(t)

	cli := initClient("test")

	defer cli.Close(context.Background())
	defer cli.DropCollection(context.Background())

	cli.EnsureIndexes(context.Background(), []string{"name"}, nil)

	var err error
	doc := bson.M{"_id": primitive.NewObjectID(), "name": "Alice"}

	opts := options.InsertOneOptions{}
	opts.InsertOneOptions = officialOpts.InsertOne().SetBypassDocumentValidation(true)
	res, err := cli.InsertOne(context.Background(), doc, opts)
	ast.NoError(err)
	ast.NotEmpty(res)
	ast.Equal(doc["_id"], res.InsertedID)

	res, err = cli.InsertOne(context.Background(), doc)
	ast.Equal(true, IsDup(err))
	ast.Empty(res)
}

func TestCollection_InsertMany(t *testing.T) {
	ast := require.New(t)
	cli := initClient("test")
	defer cli.Close(context.Background())
	defer cli.DropCollection(context.Background())
	cli.EnsureIndexes(context.Background(), []string{"name"}, nil)

	var err error
	newDocs := []UserInfo{{Id: NewObjectID(), Name: "Alice", Age: 10}, {Id: NewObjectID(), Name: "Lucas", Age: 11}}
	res, err := cli.InsertMany(context.Background(), newDocs)
	ast.NoError(err)
	ast.NotEmpty(res)
	ast.Equal(2, len(res.InsertedIDs))

	newPDocs := []*UserInfo{{Id: NewObjectID(), Name: "Alice03", Age: 10}, {Id: NewObjectID(), Name: "Lucas03", Age: 11}}
	res, err = cli.InsertMany(context.Background(), newPDocs)
	ast.NoError(err)
	ast.NotEmpty(res)
	ast.Equal(2, len(res.InsertedIDs))

	docs2 := []UserInfo{
		{Name: "Alice"},
		{Name: "Lucas"},
	}
	opts := options.InsertManyOptions{}
	opts.InsertManyOptions = officialOpts.InsertMany().SetBypassDocumentValidation(true)
	res, err = cli.InsertMany(context.Background(), docs2, opts)
	ast.Equal(true, IsDup(err))
	ast.Equal(0, len(res.InsertedIDs))

	docs4 := []UserInfo{}
	res, err = cli.InsertMany(context.Background(), []interface{}{docs4})
	ast.Error(err)
	ast.Empty(res)

}

func TestCollection_Upsert(t *testing.T) {
	ast := require.New(t)
	cli := initClient("test")
	defer cli.Close(context.Background())
	defer cli.DropCollection(context.Background())
	cli.EnsureIndexes(context.Background(), []string{"name"}, nil)

	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	docs := []interface{}{
		bson.D{{Key: "_id", Value: id1}, {Key: "name", Value: "Alice"}},
		bson.D{{Key: "_id", Value: id2}, {Key: "name", Value: "Lucas"}},
	}
	_, err := cli.InsertMany(context.Background(), docs)

	ast.NoError(err)
	// replace already exist
	filter1 := bson.M{
		"name": "Alice",
	}
	replacement1 := bson.M{
		"name": "Alice1",
		"age":  18,
	}
	opts := options.UpsertOptions{}
	opts.ReplaceOptions = officialOpts.Replace()
	res, err := cli.Upsert(context.Background(), filter1, replacement1, opts)
	ast.NoError(err)
	ast.NotEmpty(res)
	ast.Equal(int64(1), res.MatchedCount)
	ast.Equal(int64(1), res.ModifiedCount)
	ast.Equal(int64(0), res.UpsertedCount)
	ast.Equal(nil, res.UpsertedID)

	// not exist
	filter2 := bson.M{
		"name": "Lily",
	}
	replacement2 := bson.M{
		"name": "Lily",
		"age":  20,
	}
	res, err = cli.Upsert(context.Background(), filter2, replacement2)
	ast.NoError(err)
	ast.NotEmpty(res)
	ast.Equal(int64(0), res.MatchedCount)
	ast.Equal(int64(0), res.ModifiedCount)
	ast.Equal(int64(1), res.UpsertedCount)
	ast.NotNil(res.UpsertedID)

	// filter is nil or wrong BSON Document format
	replacement3 := bson.M{
		"name": "Geek",
		"age":  21,
	}
	res, err = cli.Upsert(context.Background(), nil, replacement3)
	ast.Error(err)
	ast.Empty(res)

	res, err = cli.Upsert(context.Background(), 1, replacement3)
	ast.Error(err)
	ast.Empty(res)

	// replacement is nil or wrong BSON Document format
	filter4 := bson.M{
		"name": "Geek",
	}
	res, err = cli.Upsert(context.Background(), filter4, nil)
	ast.Error(err)
	ast.Empty(res)

	res, err = cli.Upsert(context.Background(), filter4, 1)
	ast.Error(err)
	ast.Empty(res)
}

func TestCollection_UpsertId(t *testing.T) {
	ast := require.New(t)
	cli := initClient("test")
	defer cli.Close(context.Background())
	defer cli.DropCollection(context.Background())
	cli.EnsureIndexes(context.Background(), []string{"name"}, nil)

	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	docs := []interface{}{
		bson.D{{Key: "_id", Value: id1}, {Key: "name", Value: "Alice"}},
		bson.D{{Key: "_id", Value: id2}, {Key: "name", Value: "Lucas"}},
	}
	_, _ = cli.InsertMany(context.Background(), docs)

	var err error

	// replace already exist
	replacement1 := bson.M{
		"name": "Alice1",
		"age":  18,
	}
	res, err := cli.UpsertId(context.Background(), id1, replacement1)
	ast.NoError(err)
	ast.NotEmpty(res)
	ast.Equal(int64(1), res.MatchedCount)
	ast.Equal(int64(1), res.ModifiedCount)
	ast.Equal(int64(0), res.UpsertedCount)
	ast.Equal(nil, res.UpsertedID)

	// not exist filter id
	replacement2 := bson.M{
		"name": "Lily",
		"age":  20,
	}
	id3 := primitive.NewObjectID()
	opts := options.UpsertOptions{}
	opts.ReplaceOptions = officialOpts.Replace()
	res, err = cli.UpsertId(context.Background(), id3, replacement2, opts)
	ast.NoError(err)
	ast.NotEmpty(res)
	ast.Equal(int64(0), res.MatchedCount)
	ast.Equal(int64(0), res.ModifiedCount)
	ast.Equal(int64(1), res.UpsertedCount)
	ast.Equal(id3, res.UpsertedID) // id3 will insert into the inserted document

	// filter with id different from id in document, error
	id4 := primitive.NewObjectID()
	replacement3 := bson.M{
		"_id":  id4,
		"name": "Joe",
		"age":  20,
	}
	id5 := primitive.NewObjectID()
	res, err = cli.UpsertId(context.Background(), id5, replacement3)
	ast.Error(err)

	// filter is nil
	replacement4 := bson.M{
		"name": "Geek",
		"age":  21,
	}
	res, err = cli.UpsertId(context.Background(), nil, replacement4)
	ast.NoError(err)
	ast.NotEmpty(res)
	ast.Equal(int64(0), res.MatchedCount)
	ast.Equal(int64(0), res.ModifiedCount)
	ast.Equal(int64(1), res.UpsertedCount)
	ast.Nil(res.UpsertedID)

	// replacement is nil or wrong BSON Document format
	res, err = cli.UpsertId(context.Background(), id1, nil)
	ast.Error(err)
	ast.Empty(res)

	res, err = cli.UpsertId(context.Background(), id1, 1)
	ast.Error(err)
	ast.Empty(res)
}

func TestCollection_Update(t *testing.T) {
	ast := require.New(t)
	cli := initClient("test")
	defer cli.Close(context.Background())
	defer cli.DropCollection(context.Background())
	cli.EnsureIndexes(context.Background(), []string{"name"}, nil)

	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	docs := []interface{}{
		bson.D{{Key: "_id", Value: id1}, {Key: "name", Value: "Alice"}},
		bson.D{{Key: "_id", Value: id2}, {Key: "name", Value: "Lucas"}},
	}
	_, _ = cli.InsertMany(context.Background(), docs)

	var err error
	// update already exist record
	filter1 := bson.M{
		"name": "Alice",
	}
	update1 := bson.M{
		operator.Set: bson.M{
			"name": "Alice1",
			"age":  18,
		},
	}
	opts := options.UpdateOptions{}
	opts.UpdateOptions = officialOpts.Update().SetBypassDocumentValidation(false)
	err = cli.UpdateOne(context.Background(), filter1, update1, opts)
	ast.NoError(err)

	// error when not exist
	filter2 := bson.M{
		"name": "Lily",
	}
	update2 := bson.M{
		operator.Set: bson.M{
			"name": "Lily",
			"age":  20,
		},
	}
	err = cli.UpdateOne(context.Background(), filter2, update2)
	ast.Equal(err, ErrNoSuchDocuments)

	opt := officialOpts.Update().SetUpsert(true)
	opts = options.UpdateOptions{UpdateOptions: opt}
	err = cli.UpdateOne(context.Background(), filter2, update2, opts)
	ast.NoError(err)

	// filter is nil or wrong BSON Document format
	update3 := bson.M{
		"name": "Geek",
		"age":  21,
	}
	err = cli.UpdateOne(context.Background(), nil, update3)
	ast.Error(err)

	err = cli.UpdateOne(context.Background(), 1, update3)
	ast.Error(err)

	// update is nil or wrong BSON Document format
	filter4 := bson.M{
		"name": "Geek",
	}
	err = cli.UpdateOne(context.Background(), filter4, nil)
	ast.Error(err)

	err = cli.UpdateOne(context.Background(), filter4, 1)
	ast.Error(err)
}

func TestCollection_UpdateId(t *testing.T) {
	ast := require.New(t)
	cli := initClient("test")
	defer cli.Close(context.Background())
	defer cli.DropCollection(context.Background())
	cli.EnsureIndexes(context.Background(), []string{"name"}, nil)

	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	docs := []interface{}{
		bson.D{{Key: "_id", Value: id1}, {Key: "name", Value: "Alice"}},
		bson.D{{Key: "_id", Value: id2}, {Key: "name", Value: "Lucas"}},
	}
	_, _ = cli.InsertMany(context.Background(), docs)

	var err error
	// update already exist record
	update1 := bson.M{
		operator.Set: bson.M{
			"name": "Alice1",
			"age":  18,
		},
	}
	opts := options.UpdateOptions{}
	opts.UpdateOptions = officialOpts.Update().SetBypassDocumentValidation(false)
	err = cli.UpdateId(context.Background(), id1, update1, opts)
	ast.NoError(err)

	// id is nil or not exist
	update3 := bson.M{
		"name": "Geek",
		"age":  21,
	}
	err = cli.UpdateId(context.Background(), nil, update3)
	ast.Error(err)

	err = cli.UpdateId(context.Background(), 1, update3)
	ast.Error(err)

	err = cli.UpdateId(context.Background(), "not_exist_id", nil)
	ast.Error(err)

	err = cli.UpdateId(context.Background(), "not_exist_id", 1)
	ast.Error(err)
}

func TestCollection_UpdateAll(t *testing.T) {
	ast := require.New(t)
	cli := initClient("test")
	defer cli.Close(context.Background())
	defer cli.DropCollection(context.Background())
	cli.EnsureIndexes(context.Background(), nil, []string{"name"})

	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	id3 := primitive.NewObjectID()
	docs := []interface{}{
		bson.D{{Key: "_id", Value: id1}, {Key: "name", Value: "Alice"}, {Key: "age", Value: 18}},
		bson.D{{Key: "_id", Value: id2}, {Key: "name", Value: "Alice"}, {Key: "age", Value: 19}},
		bson.D{{Key: "_id", Value: id3}, {Key: "name", Value: "Lucas"}, {Key: "age", Value: 20}},
	}
	_, _ = cli.InsertMany(context.Background(), docs)

	var err error
	// update already exist record
	filter1 := bson.M{
		"name": "Alice",
	}
	update1 := bson.M{
		operator.Set: bson.M{
			"age": 33,
		},
	}
	opts := options.UpdateOptions{}
	opts.UpdateOptions = officialOpts.Update().SetBypassDocumentValidation(false)
	res, err := cli.UpdateAll(context.Background(), filter1, update1, opts)
	ast.NoError(err)
	ast.NotEmpty(res)
	ast.Equal(int64(2), res.MatchedCount)
	ast.Equal(int64(2), res.ModifiedCount)
	ast.Equal(int64(0), res.UpsertedCount)
	ast.Equal(nil, res.UpsertedID)

	// if record is not exist，err is nil， MatchedCount in res is 0
	filter2 := bson.M{
		"name": "Lily",
	}
	update2 := bson.M{
		operator.Set: bson.M{
			"age": 22,
		},
	}
	res, err = cli.UpdateAll(context.Background(), filter2, update2)
	ast.Nil(err)
	ast.NotNil(res)
	ast.Equal(int64(0), res.MatchedCount)

	// filter is nil or wrong BSON Document format
	update3 := bson.M{
		"name": "Geek",
		"age":  21,
	}
	res, err = cli.UpdateAll(context.Background(), nil, update3)
	ast.Error(err)
	ast.Nil(res)

	res, err = cli.UpdateAll(context.Background(), 1, update3)
	ast.Error(err)
	ast.Nil(res)

	// update is nil or wrong BSON Document format
	filter4 := bson.M{
		"name": "Geek",
	}
	res, err = cli.UpdateAll(context.Background(), filter4, nil)
	ast.Error(err)
	ast.Nil(res)

	res, err = cli.UpdateAll(context.Background(), filter4, 1)
	ast.Error(err)
	ast.Nil(res)
}

func TestCollection_Remove(t *testing.T) {
	ast := require.New(t)
	cli := initClient("test")
	defer cli.Close(context.Background())
	defer cli.DropCollection(context.Background())
	cli.EnsureIndexes(context.Background(), nil, []string{"name"})

	id1 := primitive.NewObjectID().Hex()
	id2 := primitive.NewObjectID().Hex()
	id3 := primitive.NewObjectID().Hex()
	id4 := primitive.NewObjectID().Hex()
	id5 := primitive.NewObjectID()
	docs := []interface{}{
		bson.D{{Key: "_id", Value: id1}, {Key: "name", Value: "Alice"}, {Key: "age", Value: 18}},
		bson.D{{Key: "_id", Value: id2}, {Key: "name", Value: "Alice"}, {Key: "age", Value: 19}},
		bson.D{{Key: "_id", Value: id3}, {Key: "name", Value: "Lucas"}, {Key: "age", Value: 20}},
		bson.D{{Key: "_id", Value: id4}, {Key: "name", Value: "Joe"}, {Key: "age", Value: 20}},
		bson.D{{Key: "_id", Value: id5}, {Key: "name", Value: "Ethan"}, {Key: "age", Value: 1}},
	}
	_, _ = cli.InsertMany(context.Background(), docs)

	var err error
	// remove id
	err = cli.RemoveId(context.Background(), "")
	ast.Error(err)
	err = cli.RemoveId(context.Background(), "not-exists-id")
	ast.True(IsErrNoDocuments(err))
	ast.NoError(cli.RemoveId(context.Background(), id4))
	ast.NoError(cli.RemoveId(context.Background(), id5))

	// delete record: name = "Alice" , after that, expect one name = "Alice" record
	filter1 := bson.M{
		"name": "Alice",
	}
	opts := options.RemoveOptions{}
	opts.DeleteOptions = officialOpts.Delete()
	err = cli.Remove(context.Background(), filter1, opts)
	ast.NoError(err)

	cnt, err := cli.Find(context.Background(), filter1).Count()
	ast.NoError(err)
	ast.Equal(int64(1), cnt)

	// delete not match  record , report err
	filter2 := bson.M{
		"name": "Lily",
	}
	err = cli.Remove(context.Background(), filter2)
	ast.Equal(err, ErrNoSuchDocuments)

	// filter is bson.M{}，delete one document
	filter3 := bson.M{}
	preCnt, err := cli.Find(context.Background(), filter3).Count()
	ast.NoError(err)
	ast.Equal(int64(2), preCnt)

	err = cli.Remove(context.Background(), filter3)
	ast.NoError(err)

	afterCnt, err := cli.Find(context.Background(), filter3).Count()
	ast.NoError(err)
	ast.Equal(preCnt-1, afterCnt)

	// filter is nil or wrong BSON Document format
	err = cli.Remove(context.Background(), nil)
	ast.Error(err)

	err = cli.Remove(context.Background(), 1)
	ast.Error(err)
}

func TestCollection_RemoveAll(t *testing.T) {
	ast := require.New(t)
	cli := initClient("test")
	defer cli.Close(context.Background())
	defer cli.DropCollection(context.Background())
	cli.EnsureIndexes(context.Background(), nil, []string{"name"})

	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	id3 := primitive.NewObjectID()
	id4 := primitive.NewObjectID()
	docs := []interface{}{
		bson.D{{Key: "_id", Value: id1}, {Key: "name", Value: "Alice"}, {Key: "age", Value: 18}},
		bson.D{{Key: "_id", Value: id2}, {Key: "name", Value: "Alice"}, {Key: "age", Value: 19}},
		bson.D{{Key: "_id", Value: id3}, {Key: "name", Value: "Lucas"}, {Key: "age", Value: 20}},
		bson.D{{Key: "_id", Value: id4}, {Key: "name", Value: "Rocket"}, {Key: "age", Value: 23}},
	}
	_, _ = cli.InsertMany(context.Background(), docs)

	var err error
	// delete record: name = "Alice" ,after that, expect - record : name = "Alice"
	filter1 := bson.M{
		"name": "Alice",
	}
	opts := options.RemoveOptions{}
	opts.DeleteOptions = officialOpts.Delete()
	res, err := cli.RemoveAll(context.Background(), filter1, opts)
	ast.NoError(err)
	ast.NotNil(res)
	ast.Equal(int64(2), res.DeletedCount)

	cnt, err := cli.Find(context.Background(), filter1).Count()
	ast.NoError(err)
	ast.Equal(int64(0), cnt)

	// delete with not match filter， DeletedCount in res is 0
	filter2 := bson.M{
		"name": "Lily",
	}
	res, err = cli.RemoveAll(context.Background(), filter2)
	ast.NoError(err)
	ast.NotNil(res)
	ast.Equal(int64(0), res.DeletedCount)

	// filter is bson.M{}，delete all docs
	filter3 := bson.M{}
	preCnt, err := cli.Find(context.Background(), filter3).Count()
	ast.NoError(err)
	ast.Equal(int64(2), preCnt)

	res, err = cli.RemoveAll(context.Background(), filter3)
	ast.NoError(err)
	ast.NotNil(res)
	ast.Equal(preCnt, res.DeletedCount)

	afterCnt, err := cli.Find(context.Background(), filter3).Count()
	ast.NoError(err)
	ast.Equal(int64(0), afterCnt)

	// filter is nil or wrong BSON Document format
	res, err = cli.RemoveAll(context.Background(), nil)
	ast.Error(err)
	ast.Nil(res)

	res, err = cli.RemoveAll(context.Background(), 1)
	ast.Error(err)
	ast.Nil(res)
}
func TestSliceInsert(t *testing.T) {
	newDocs := []UserInfo{{Name: "Alice", Age: 10}, {Name: "Lucas", Age: 11}}
	di := interface{}(newDocs)
	dis := interfaceToSliceInterface(di)
	ast := require.New(t)
	ast.Len(dis, 2)

	newDocs_1 := []interface{}{UserInfo{Name: "Alice", Age: 10}, UserInfo{Name: "Lucas", Age: 11}}
	di = interface{}(newDocs_1)
	dis = interfaceToSliceInterface(di)
	ast.Len(dis, 2)

	newDocs_2 := UserInfo{Name: "Alice", Age: 10}
	di = interface{}(newDocs_2)
	dis = interfaceToSliceInterface(di)
	ast.Nil(dis)

	newDocs_3 := []UserInfo{}
	di = interface{}(newDocs_3)
	dis = interfaceToSliceInterface(di)
	ast = require.New(t)
	ast.Nil(dis)
}

func TestCollection_ReplaceOne(t *testing.T) {
	ast := require.New(t)
	cli := initClient("test")
	defer cli.Close(context.Background())
	defer cli.DropCollection(context.Background())
	cli.EnsureIndexes(context.Background(), []string{"name"}, nil)

	id := primitive.NewObjectID()
	ui := UserInfo{Id: id, Name: "Lucas", Age: 17}
	_, err := cli.InsertOne(context.Background(), ui)
	ast.NoError(err)
	ui.Id = id
	ui.Age = 27
	err = cli.ReplaceOne(context.Background(), bson.M{"_id": id}, &ui)
	ast.NoError(err)

	findUi := UserInfo{}
	err = cli.Find(context.Background(), bson.M{"name": "Lucas"}).One(&findUi)
	ast.NoError(err)
	ast.Equal(ui.Age, findUi.Age)

	opts := options.ReplaceOptions{}
	opts.ReplaceOptions = officialOpts.Replace()
	err = cli.ReplaceOne(context.Background(), bson.M{"_id": "notexist"}, &ui, opts)
	ast.Equal(ErrNoSuchDocuments, err)

	err = cli.ReplaceOne(context.Background(), bson.M{"_id": "notexist"}, nil)
	ast.Error(err)
}

func TestChangeStream(t *testing.T) {
	ast := require.New(t)
	cli := initClient("test")
	defer cli.Close(context.Background())
	defer cli.DropCollection(context.Background())

	opts := &options.ChangeStreamOptions{officialOpts.ChangeStream()}
	c, e := cli.Watch(context.Background(), mongo.Pipeline{}, opts)
	ast.NoError(e)
	defer c.Close(context.Background())

	doneChane := make(chan struct{})
	go func() {
		ok := c.Next(context.Background())
		ast.True(ok)
		doneChane <- struct{}{}
	}()

	id := primitive.NewObjectID()
	ui := UserInfo{Id: id, Name: "Lucas", Age: 17}
	_, err := cli.InsertOne(context.Background(), ui)
	ast.NoError(err)
	<-doneChane

}
