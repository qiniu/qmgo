package mongox

import "go.mongodb.org/mongo-driver/mongo"

// CollectionI
// 集合操作接口
type CollectionI interface {
	Find(filter interface{}) QueryI
	InsertOne(doc interface{}) (*mongo.InsertOneResult, error)
	InsertMany(docs ...interface{}) (*mongo.InsertManyResult, error)
	Upsert(filter interface{}, replacement interface{}) (*mongo.UpdateResult, error)
	UpdateOne(filter interface{}, update interface{}) error
	UpdateAll(filter interface{}, update interface{}) (*mongo.UpdateResult, error)
	DeleteOne(filter interface{}) error
	DeleteAll(selector interface{}) (*mongo.DeleteResult, error)
	EnsureIndex(indexes []string, isUnique bool)
	EnsureIndexes(uniques []string, indexes []string)
}

// CursorI
// 游标接口
type CursorI interface {
	Next(result interface{}) bool
	Close() error
	Err() error
}

// QueryI
// 查询接口
type QueryI interface {
	Sort(fields ...string) QueryI
	Select(selector interface{}) QueryI
	Skip(n int64) QueryI
	Limit(n int64) QueryI
	One(result interface{}) error
	All(result interface{}) error
	Count() (n int64, err error)
	Distinct(key string, result interface{}) error
	Cursor() (CursorI, error)
}
