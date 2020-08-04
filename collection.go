package mongox

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

var (
	NoSuchRecordErr = errors.New("no such record")
)

// Collection
type Collection struct {
	Ctx        context.Context //注入ctx，便于未来注入trace
	Collection *mongo.Collection
}

// Find
// 查询函数，插入查询条件 filter，返回一个 QueryI 对象
func (c *Collection) Find(filter interface{}) QueryI {
	return &Query{
		ctx:        c.Ctx,
		collection: c.Collection,
		filter:     filter,
	}
}

// InsertOne
// 插入一条记录
// doc 中如果没有 _id 字段，将自动生成
// InsertOneResult 中返回插入记录的 _id
// 语法参考 https://docs.mongodb.com/manual/reference/command/insert/
func (c *Collection) InsertOne(doc interface{}) (*mongo.InsertOneResult, error) {
	var err error
	var res *mongo.InsertOneResult

	res, err = c.Collection.InsertOne(c.Ctx, doc)

	return res, err
}

// InsertMany
// 插入多条记录
// docs 中如果某条记录没有 _id 字段，将自动生成
// InsertManyResult 中返回插入记录的 _id slice
// 语法参考 https://docs.mongodb.com/manual/reference/command/insert/
func (c *Collection) InsertMany(docs interface{}) (*mongo.InsertManyResult, error) {
	var err error
	var res *mongo.InsertManyResult

	var list []interface{}
	switch v := docs.(type) {
	case []bson.M:
		for _, x := range v {
			list = append(list, x)
		}
	case []bson.D:
		for _, x := range v {
			list = append(list, x)
		}
	default:
		err = errors.New("docs type do not []bson.M or []bson.D")
		return res, err
	}

	res, err = c.Collection.InsertMany(c.Ctx, list)

	return res, err
}

// Upsert
// 如果 filter 可以筛选出一条记录，则使用 replacement 替换原来的记录
// 如果 filter 未筛选出记录，则将 replacement 插入
// replacement 语法参考 https://docs.mongodb.com/manual/reference/operator/update/
func (c *Collection) Upsert(filter interface{}, replacement interface{}) (*mongo.UpdateResult, error) {
	var err error
	var res *mongo.UpdateResult

	opts := options.Replace().SetUpsert(true)
	res, err = c.Collection.ReplaceOne(c.Ctx, filter, replacement, opts)

	return res, err
}

// UpdateOne
// 更新 filter 条件筛选出来的一条记录，如果筛选失败，返回错误
// update 语法参考 https://docs.mongodb.com/manual/reference/operator/update/
func (c *Collection) UpdateOne(filter interface{}, update interface{}) error {
	var err error
	var res *mongo.UpdateResult

	if res, err = c.Collection.UpdateOne(c.Ctx, filter, update); err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		err = NoSuchRecordErr
	}

	return err
}

// UpdateAll
// 更新 filter 条件筛选出来的所有记录
// 如果筛选失败，UpdateResult 字段中 MatchedCount 为 0
// update 语法参考 https://docs.mongodb.com/manual/reference/operator/update/
func (c *Collection) UpdateAll(filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	var err error
	var res *mongo.UpdateResult

	if res, err = c.Collection.UpdateMany(c.Ctx, filter, update); err != nil {
		return nil, err
	}

	return res, err
}

// DeleteOne
// 根据 filter 筛选一条记录删除
// 特别注意，如果 filter 为 bson.M{}，会删除 Collection 中的一条记录
// 如果筛选失败，返回错误
// 语法参考 https://docs.mongodb.com/manual/reference/command/delete/
func (c *Collection) DeleteOne(filter interface{}) error {
	var err error
	var res *mongo.DeleteResult

	if res, err = c.Collection.DeleteOne(c.Ctx, filter); err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		err = NoSuchRecordErr
	}

	return err
}

// DeleteAll
// 删除根据 filter 筛选出的所有记录
// 特别注意，如果 filter 为 bson.M{}，会导致 Collection 中所有记录删除
// 如果筛选失败，DeleteResult 字段中 DeletedCount 为 0
// 语法参考 https://docs.mongodb.com/manual/reference/command/delete/
func (c *Collection) DeleteAll(filter interface{}) (*mongo.DeleteResult, error) {
	var err error
	var result *mongo.DeleteResult

	result, err = c.Collection.DeleteMany(c.Ctx, filter)

	return result, err
}

// EnsureIndex
// 创建一组同类索引
// isUnique, true表示唯一索引，false表示普通索引
// 语法参考 https://docs.mongodb.com/manual/reference/command/createIndexes/
// indexes 例：[]string{"idx1", "-idx2", "idx3,idx4"}
// 上例建立了3个索引，idx1 字段按升序排列的索引，idx2 字段按降序排列的索引，以及 idx3和idx4 都按升序排序的联合索引
func (c *Collection) EnsureIndex(indexes []string, isUnique bool) {
	var indexModels []mongo.IndexModel

	// 组建[]mongo.IndexModel
	for _, idx := range indexes {
		var model mongo.IndexModel
		var keysDoc bsonx.Doc

		colIndexArr := strings.Split(idx, ",")
		for _, field := range colIndexArr {
			key, n := SplitSymbol(field)

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
	res, err = c.Collection.Indexes().CreateMany(c.Ctx, indexModels)

	if err != nil || len(res) == 0 {
		s := fmt.Sprint("<MongoDB.C>: ", c.Collection.Name(), " Index: ", indexes, " error: ", err, "res: ", res)
		panic(s)
	}
	return
}

// EnsureIndexes
// 同时创建唯一索引和普通索引，uniques表示唯一索引，indexes表示普通索引
func (c *Collection) EnsureIndexes(uniques []string, indexes []string) {
	// 创建唯一索引
	if len(uniques) != 0 {
		c.EnsureIndex(uniques, true)
	}

	// 创建普通索引
	if len(indexes) != 0 {
		c.EnsureIndex(indexes, false)
	}

	return
}

// DropCollection
// 当collection不存在mongodb server时，函数执行也是安全的
func (c *Collection) DropCollection() error {
	return c.Collection.Drop(c.Ctx)
}
