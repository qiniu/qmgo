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

import "go.mongodb.org/mongo-driver/mongo/options"

// CollectionI
// 集合操作接口
//type CollectionI interface {
//	Find(filter interface{}) QueryI
//	InsertOne(doc interface{}) (*mongo.InsertOneResult, error)
//	InsertMany(docs ...interface{}) (*mongo.InsertManyResult, error)
//	Upsert(filter interface{}, replacement interface{}) (*mongo.UpdateResult, error)
//	UpdateOne(filter interface{}, update interface{}) error
//	UpdateAll(filter interface{}, update interface{}) (*mongo.UpdateResult, error)
//	DeleteOne(filter interface{}) error
//	RemoveAll(selector interface{}) (*mongo.DeleteResult, error)
//	EnsureIndex(indexes []string, isUnique bool)
//	EnsureIndexes(uniques []string, indexes []string)
//}

// Change holds fields for running a findAndModify command via the Query.Apply method.
type Change struct {
	Update    interface{} // update/replace document
	Replace   bool        // Whether to replace the document rather than updating
	Remove    bool        // Whether to remove the document found rather than updating
	Upsert    bool        // Whether to insert in case the document isn't found, take effect when Remove is false
	ReturnNew bool        // Should the modified document be returned rather than the old one, take effect when Remove is false
}

// CursorI Cursor interface
type CursorI interface {
	Next(result interface{}) bool
	Close() error
	Err() error
	All(results interface{}) error
	//ID() int64
}

// QueryI Query interface
type QueryI interface {
	Collation(collation *options.Collation) QueryI
	SetArrayFilters(*options.ArrayFilters) QueryI
	Sort(fields ...string) QueryI
	Select(selector interface{}) QueryI
	Skip(n int64) QueryI
	BatchSize(n int64) QueryI
	NoCursorTimeout(n bool) QueryI
	Limit(n int64) QueryI
	One(result interface{}) error
	All(result interface{}) error
	Count() (n int64, err error)
	EstimatedCount() (n int64, err error)
	Distinct(key string, result interface{}) error
	Cursor() CursorI
	Apply(change Change, result interface{}) error
	Hint(hint interface{}) QueryI
}

// AggregateI define the interface of aggregate
type AggregateI interface {
	All(results interface{}) error
	One(result interface{}) error
	Iter() CursorI // Deprecated, please use Cursor instead
	Cursor() CursorI
}
