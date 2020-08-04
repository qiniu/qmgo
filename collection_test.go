package mongox

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestCollection_EnsureIndex(t *testing.T) {
	ast := require.New(t)

	var cli *Client
	var coll *Collection

	cli = initClient("test")
	coll = cli.GetCollection(context.TODO())
	coll.DropCollection()

	// 正常设置索引
	coll.EnsureIndex(nil, false)
	coll.EnsureIndex([]string{"id1"}, true)
	coll.EnsureIndex([]string{"id2,id3"}, false)
	coll.EnsureIndex([]string{"id4,-id5"}, false)

	// 设置已存在的索引，且索引内容不一样，会发生panic
	ast.Panics(func() { coll.EnsureIndex([]string{"id1"}, false) })

	// 检测唯一索引是否设置正常
	var err error
	doc := bson.M{
		"id1": 1,
	}

	_, err = coll.InsertOne(doc)
	ast.NoError(err)
	_, err = coll.InsertOne(doc)
	ast.Equal(true, IsDup(err))
}

func TestCollection_EnsureIndexes(t *testing.T) {
	ast := require.New(t)

	var cli *Client
	var coll *Collection

	cli = initClient("test")
	coll = cli.GetCollection(context.TODO())
	coll.DropCollection()

	// 正常设置索引
	unique := []string{"id1"}
	common := []string{"id2,id3", "id4,-id5"}
	coll.EnsureIndexes(unique, common)

	// 设置已存在的索引，且索引内容不一样，会发生panic
	ast.Panics(func() { coll.EnsureIndexes(nil, unique) })

	// 检测唯一索引是否设置正常
	var err error
	doc := bson.M{
		"id1": 1,
	}

	_, err = coll.InsertOne(doc)
	ast.NoError(err)
	_, err = coll.InsertOne(doc)
	ast.Equal(true, IsDup(err))
}

func TestCollection_InsertOne(t *testing.T) {
	ast := require.New(t)

	var cli *Client
	var coll *Collection

	cli = initClient("test")
	coll = cli.GetCollection(context.TODO())
	coll.DropCollection()
	coll.EnsureIndexes([]string{"name"}, nil)

	var err error
	var res *mongo.InsertOneResult
	doc := bson.M{"_id": primitive.NewObjectID(), "name": "Alice"}

	res, err = coll.InsertOne(doc)
	ast.NoError(err)
	ast.NotEmpty(res)
	ast.Equal(doc["_id"], res.InsertedID)

	res, err = coll.InsertOne(doc)
	ast.Equal(true, IsDup(err))
	ast.Empty(res)
}

func TestCollection_InsertMany(t *testing.T) {
	ast := require.New(t)

	var cli *Client
	var coll *Collection

	cli = initClient("test")
	coll = cli.GetCollection(context.TODO())
	coll.DropCollection()
	coll.EnsureIndexes([]string{"name"}, nil)

	var err error
	var res *mongo.InsertManyResult

	docs := []bson.D{
		{{Key: "_id", Value: primitive.NewObjectID()}, {Key: "name", Value: "Alice"}},
		{{Key: "_id", Value: primitive.NewObjectID()}, {Key: "name", Value: "Lucas"}},
	}
	res, err = coll.InsertMany(docs)
	ast.NoError(err)
	ast.NotEmpty(res)
	ast.Equal(2, len(res.InsertedIDs))

	docs2 := []bson.M{
		{"_id": primitive.NewObjectID(), "name": "Alice"},
		{"_id": primitive.NewObjectID(), "name": "Lucas"},
	}
	res, err = coll.InsertMany(docs2)
	ast.Equal(true, IsDup(err))
	ast.Equal(0, len(res.InsertedIDs))

	docs3 := make(map[string]string, 10)
	docs3["name"] = "Lily"
	res, err = coll.InsertMany(docs3)
	ast.EqualError(err, "docs type do not []bson.M or []bson.D")
	ast.Empty(res)

	docs4 := []bson.M{}
	res, err = coll.InsertMany(docs4)
	ast.Error(err)
	ast.Empty(res)
}

func TestCollection_Upsert(t *testing.T) {
	ast := require.New(t)

	var cli *Client
	var coll *Collection

	cli = initClient("test")
	coll = cli.GetCollection(context.TODO())
	coll.DropCollection()
	coll.EnsureIndexes([]string{"name"}, nil)

	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	docs := []bson.D{
		{{Key: "_id", Value: id1}, {Key: "name", Value: "Alice"}},
		{{Key: "_id", Value: id2}, {Key: "name", Value: "Lucas"}},
	}
	_, _ = coll.InsertMany(docs)

	var err error
	var res *mongo.UpdateResult

	// 已存在记录执行操作，替换记录
	filter1 := bson.M{
		"name": "Alice",
	}
	replacement1 := bson.M{
		"name": "Alice1",
		"age":  18,
	}
	res, err = coll.Upsert(filter1, replacement1)
	ast.NoError(err)
	ast.NotEmpty(res)
	ast.Equal(int64(1), res.MatchedCount)
	ast.Equal(int64(1), res.ModifiedCount)
	ast.Equal(int64(0), res.UpsertedCount)
	ast.Equal(nil, res.UpsertedID)

	// 记录不存在，插入一条新纪录
	filter2 := bson.M{
		"name": "Lily",
	}
	replacement2 := bson.M{
		"name": "Lily",
		"age":  20,
	}
	res, err = coll.Upsert(filter2, replacement2)
	ast.NoError(err)
	ast.NotEmpty(res)
	ast.Equal(int64(0), res.MatchedCount)
	ast.Equal(int64(0), res.ModifiedCount)
	ast.Equal(int64(1), res.UpsertedCount)
	ast.NotNil(res.UpsertedID)

	// filter 为nil 或者不符合 BSON Document 规范
	replacement3 := bson.M{
		"name": "Geek",
		"age":  21,
	}
	res, err = coll.Upsert(nil, replacement3)
	ast.Error(err)
	ast.Empty(res)

	res, err = coll.Upsert(1, replacement3)
	ast.Error(err)
	ast.Empty(res)

	// replacement 为nil 或者不符合 BSON Document 规范
	filter4 := bson.M{
		"name": "Geek",
	}
	res, err = coll.Upsert(filter4, nil)
	ast.Error(err)
	ast.Empty(res)

	res, err = coll.Upsert(filter4, 1)
	ast.Error(err)
	ast.Empty(res)
}

func TestCollection_UpdateOne(t *testing.T) {
	ast := require.New(t)

	var cli *Client
	var coll *Collection

	cli = initClient("test")
	coll = cli.GetCollection(context.TODO())
	coll.DropCollection()
	coll.EnsureIndexes([]string{"name"}, nil)

	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	docs := []bson.D{
		{{Key: "_id", Value: id1}, {Key: "name", Value: "Alice"}},
		{{Key: "_id", Value: id2}, {Key: "name", Value: "Lucas"}},
	}
	_, _ = coll.InsertMany(docs)

	var err error
	// 已存在记录执行操作，更新记录
	filter1 := bson.M{
		"name": "Alice",
	}
	update1 := bson.M{
		"$set": bson.M{
			"name": "Alice1",
			"age":  18,
		},
	}
	err = coll.UpdateOne(filter1, update1)
	ast.NoError(err)

	// 记录不存在，报错
	filter2 := bson.M{
		"name": "Lily",
	}
	update2 := bson.M{
		"$set": bson.M{
			"name": "Lily",
			"age":  20,
		},
	}
	err = coll.UpdateOne(filter2, update2)
	ast.Equal(err, NoSuchRecordErr)

	// filter 为nil 或者不符合 BSON Document 规范
	update3 := bson.M{
		"name": "Geek",
		"age":  21,
	}
	err = coll.UpdateOne(nil, update3)
	ast.Error(err)

	err = coll.UpdateOne(1, update3)
	ast.Error(err)

	// update 为nil 或者不符合 BSON Document 规范
	filter4 := bson.M{
		"name": "Geek",
	}
	err = coll.UpdateOne(filter4, nil)
	ast.Error(err)

	err = coll.UpdateOne(filter4, 1)
	ast.Error(err)
}

func TestCollection_UpdateAll(t *testing.T) {
	ast := require.New(t)

	var cli *Client
	var coll *Collection

	cli = initClient("test")
	coll = cli.GetCollection(context.TODO())
	coll.DropCollection()
	coll.EnsureIndexes(nil, []string{"name"})

	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	id3 := primitive.NewObjectID()
	docs := []bson.D{
		{{Key: "_id", Value: id1}, {Key: "name", Value: "Alice"}, {Key: "age", Value: 18}},
		{{Key: "_id", Value: id2}, {Key: "name", Value: "Alice"}, {Key: "age", Value: 19}},
		{{Key: "_id", Value: id3}, {Key: "name", Value: "Lucas"}, {Key: "age", Value: 20}},
	}
	_, _ = coll.InsertMany(docs)

	var err error
	var res *mongo.UpdateResult
	// 已存在记录执行操作，更新记录
	filter1 := bson.M{
		"name": "Alice",
	}
	update1 := bson.M{
		"$set": bson.M{
			"age": 33,
		},
	}
	res, err = coll.UpdateAll(filter1, update1)
	ast.NoError(err)
	ast.NotEmpty(res)
	ast.Equal(int64(2), res.MatchedCount)
	ast.Equal(int64(2), res.ModifiedCount)
	ast.Equal(int64(0), res.UpsertedCount)
	ast.Equal(nil, res.UpsertedID)

	// 记录不存在，err 为 nil，res 里 MatchedCount 为 0
	filter2 := bson.M{
		"name": "Lily",
	}
	update2 := bson.M{
		"$set": bson.M{
			"age": 22,
		},
	}
	res, err = coll.UpdateAll(filter2, update2)
	ast.Nil(err)
	ast.NotNil(res)
	ast.Equal(int64(0), res.MatchedCount)

	// filter 为nil 或者不符合 BSON Document 规范
	update3 := bson.M{
		"name": "Geek",
		"age":  21,
	}
	res, err = coll.UpdateAll(nil, update3)
	ast.Error(err)
	ast.Nil(res)

	res, err = coll.UpdateAll(1, update3)
	ast.Error(err)
	ast.Nil(res)

	// update 为nil 或者不符合 BSON Document 规范
	filter4 := bson.M{
		"name": "Geek",
	}
	res, err = coll.UpdateAll(filter4, nil)
	ast.Error(err)
	ast.Nil(res)

	res, err = coll.UpdateAll(filter4, 1)
	ast.Error(err)
	ast.Nil(res)
}

func TestCollection_DeleteOne(t *testing.T) {
	ast := require.New(t)

	var cli *Client
	var coll *Collection

	cli = initClient("test")
	coll = cli.GetCollection(context.TODO())
	coll.DropCollection()
	coll.EnsureIndexes(nil, []string{"name"})

	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	id3 := primitive.NewObjectID()
	docs := []bson.D{
		{{Key: "_id", Value: id1}, {Key: "name", Value: "Alice"}, {Key: "age", Value: 18}},
		{{Key: "_id", Value: id2}, {Key: "name", Value: "Alice"}, {Key: "age", Value: 19}},
		{{Key: "_id", Value: id3}, {Key: "name", Value: "Lucas"}, {Key: "age", Value: 20}},
	}
	_, _ = coll.InsertMany(docs)

	var err error
	// 删除一条 name = "Alice" 的记录，删除成功后，集合里还存在一条 name = "Alice"的记录
	filter1 := bson.M{
		"name": "Alice",
	}
	err = coll.DeleteOne(filter1)
	ast.NoError(err)

	cnt, err := coll.Find(filter1).Count()
	ast.NoError(err)
	ast.Equal(int64(1), cnt)

	// 删除 filter 无记录，报错
	filter2 := bson.M{
		"name": "Lily",
	}
	err = coll.DeleteOne(filter2)
	ast.Equal(err, NoSuchRecordErr)

	// filter 为 bson.M{}，会导致某一条记录被删除
	filter3 := bson.M{}
	preCnt, err := coll.Find(filter3).Count()
	ast.Equal(int64(2), preCnt)

	err = coll.DeleteOne(filter3)
	ast.NoError(err)

	afterCnt, err := coll.Find(filter3).Count()
	ast.Equal(preCnt-1, afterCnt)

	// filter 为 nil 或者不符合 BSON Document 规范
	err = coll.DeleteOne(nil)
	ast.Error(err)

	err = coll.DeleteOne(1)
	ast.Error(err)
}

func TestCollection_DeleteAll(t *testing.T) {
	ast := require.New(t)

	var cli *Client
	var coll *Collection

	cli = initClient("test")
	coll = cli.GetCollection(context.TODO())
	coll.DropCollection()
	coll.EnsureIndexes(nil, []string{"name"})

	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	id3 := primitive.NewObjectID()
	id4 := primitive.NewObjectID()
	docs := []bson.D{
		{{Key: "_id", Value: id1}, {Key: "name", Value: "Alice"}, {Key: "age", Value: 18}},
		{{Key: "_id", Value: id2}, {Key: "name", Value: "Alice"}, {Key: "age", Value: 19}},
		{{Key: "_id", Value: id3}, {Key: "name", Value: "Lucas"}, {Key: "age", Value: 20}},
		{{Key: "_id", Value: id4}, {Key: "name", Value: "Rocket"}, {Key: "age", Value: 23}},
	}
	_, _ = coll.InsertMany(docs)

	var err error
	var res *mongo.DeleteResult
	// 删除 name = "Alice" 的记录，删除成功后，集合里name = "Alice"的记录为 0
	filter1 := bson.M{
		"name": "Alice",
	}
	res, err = coll.DeleteAll(filter1)
	ast.NoError(err)
	ast.NotNil(res)
	ast.Equal(int64(2), res.DeletedCount)

	cnt, err := coll.Find(filter1).Count()
	ast.NoError(err)
	ast.Equal(int64(0), cnt)

	// 删除 filter 无记录，res 里 DeletedCount 为 0
	filter2 := bson.M{
		"name": "Lily",
	}
	res, err = coll.DeleteAll(filter2)
	ast.NoError(err)
	ast.NotNil(res)
	ast.Equal(int64(0), res.DeletedCount)

	// filter 为 bson.M{}，会导致集合里所有记录被删除
	filter3 := bson.M{}
	preCnt, err := coll.Find(filter3).Count()
	ast.NoError(err)
	ast.Equal(int64(2), preCnt)

	res, err = coll.DeleteAll(filter3)
	ast.NoError(err)
	ast.NotNil(res)
	ast.Equal(preCnt, res.DeletedCount)

	afterCnt, err := coll.Find(filter3).Count()
	ast.NoError(err)
	ast.Equal(int64(0), afterCnt)

	// filter 为 nil 或者不符合 BSON Document 规范
	res, err = coll.DeleteAll(nil)
	ast.Error(err)
	ast.Nil(res)

	res, err = coll.DeleteAll(1)
	ast.Error(err)
	ast.Nil(res)
}
