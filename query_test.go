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
	"testing"
	"time"

	"github.com/qiniu/qmgo/operator"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type QueryTestItem struct {
	Id   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
	Age  int                `bson:"age"`

	Instock []struct {
		Warehouse string `bson:"warehouse"`
		Qty       int    `bson:"qty"`
	} `bson:"instock"`
}

type QueryTestItem2 struct {
	Class string `bson:"class"`
}

func TestQuery_One(t *testing.T) {
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
	var res QueryTestItem

	filter1 := bson.M{
		"name": "Alice",
	}
	projection1 := bson.M{
		"age": 0,
	}

	err = cli.Find(context.Background(), filter1).Select(projection1).Sort("age").Limit(1).Skip(1).One(&res)
	ast.Nil(err)
	ast.Equal(id2, res.Id)
	ast.Equal("Alice", res.Name)

	res = QueryTestItem{}
	filter2 := bson.M{
		"name": "Lily",
	}

	err = cli.Find(context.Background(), filter2).One(&res)
	ast.Error(err)
	ast.Empty(res)

	// filter is bson.M{}，match all and return one
	res = QueryTestItem{}
	filter3 := bson.M{}

	err = cli.Find(context.Background(), filter3).One(&res)
	ast.NoError(err)
	ast.NotEmpty(res)

	// filter is nil，error
	res = QueryTestItem{}
	err = cli.Find(context.Background(), nil).One(&res)
	ast.Error(err)
	ast.Empty(res)

	// res is nil or can't parse
	err = cli.Find(context.Background(), filter1).One(nil)
	ast.Error(err)

	var tv int
	err = cli.Find(context.Background(), filter1).One(&tv)
	ast.Error(err)

	// res is a parseable object, but the bson tag is inconsistent with the mongodb record, no error is reported, res is the initialization state of the data structure
	var tt QueryTestItem2
	err = cli.Find(context.Background(), filter1).One(&tt)
	ast.NoError(err)
	ast.Empty(tt)
}

func TestQuery_All(t *testing.T) {
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
		bson.M{"_id": id1, "name": "Alice", "age": 18},
		bson.M{"_id": id2, "name": "Alice", "age": 19},
		bson.M{"_id": id3, "name": "Lucas", "age": 20},
		bson.M{"_id": id4, "name": "Lucas", "age": 21},
	}
	_, _ = cli.InsertMany(context.Background(), docs)

	var err error
	var res []QueryTestItem

	filter1 := bson.M{
		"name": "Alice",
	}
	projection1 := bson.M{
		"name": 0,
	}

	err = cli.Find(context.Background(), filter1).Select(projection1).Sort("age").Limit(2).Skip(1).All(&res)
	ast.NoError(err)
	ast.Equal(1, len(res))

	res = make([]QueryTestItem, 0)
	filter2 := bson.M{
		"name": "Lily",
	}

	err = cli.Find(context.Background(), filter2).All(&res)
	ast.NoError(err)
	ast.Empty(res)

	// filter is bson.M{}, which means to match all, will return all records in the collection
	res = make([]QueryTestItem, 0)
	filter3 := bson.M{}

	err = cli.Find(context.Background(), filter3).All(&res)
	ast.NoError(err)
	ast.Equal(4, len(res))

	res = make([]QueryTestItem, 0)
	err = cli.Find(context.Background(), nil).All(&res)
	ast.Error(err)
	ast.Empty(res)

	err = cli.Find(context.Background(), filter1).All(nil)
	ast.Error(err)

	var tv int
	err = cli.Find(context.Background(), filter1).All(&tv)
	ast.Error(err)
	// res is a parseable object, but the bson tag is inconsistent with the mongodb record, and no error is reported
	// The corresponding value will be mapped according to the bson tag of the res data structure, and the tag without the value will be the default value of the corresponding type
	// The length of res is the number of records filtered by the filter condition
	var tt []QueryTestItem2
	err = cli.Find(context.Background(), filter1).All(&tt)
	ast.NoError(err)
	ast.Equal(2, len(tt))
}

func TestQuery_Count(t *testing.T) {
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
		bson.M{"_id": id1, "name": "Alice", "age": 18},
		bson.M{"_id": id2, "name": "Alice", "age": 19},
		bson.M{"_id": id3, "name": "Lucas", "age": 20},
		bson.M{"_id": id4, "name": "Lucas", "age": 21},
	}
	_, _ = cli.InsertMany(context.Background(), docs)

	var err error
	var cnt int64

	filter1 := bson.M{
		"name": "Alice",
	}

	cnt, err = cli.Find(context.Background(), filter1).Limit(2).Skip(1).Count()
	ast.NoError(err)
	ast.Equal(int64(1), cnt)

	filter2 := bson.M{
		"name": "Lily",
	}

	cnt, err = cli.Find(context.Background(), filter2).Count()
	ast.NoError(err)
	ast.Zero(cnt)

	filter3 := bson.M{}

	cnt, err = cli.Find(context.Background(), filter3).Count()
	ast.NoError(err)
	ast.Equal(int64(4), cnt)

	cnt, err = cli.Find(context.Background(), nil).Count()
	ast.Error(err)
	ast.Zero(cnt)
}

func TestQuery_Skip(t *testing.T) {
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
		bson.M{"_id": id1, "name": "Alice", "age": 18},
		bson.M{"_id": id2, "name": "Alice", "age": 19},
		bson.M{"_id": id3, "name": "Lucas", "age": 20},
		bson.M{"_id": id4, "name": "Lucas", "age": 21},
	}
	_, _ = cli.InsertMany(context.Background(), docs)

	var err error
	var res []QueryTestItem

	// filter can match records, skip 1 record, and return the remaining records
	filter1 := bson.M{
		"name": "Alice",
	}

	err = cli.Find(context.Background(), filter1).Skip(1).All(&res)
	ast.NoError(err)
	ast.Equal(1, len(res))

	// filter can match the records, the number of skips is greater than the total number of existing records, res returns empty
	res = make([]QueryTestItem, 0)

	err = cli.Find(context.Background(), filter1).Skip(3).All(&res)
	ast.NoError(err)
	ast.Empty(res)

	res = make([]QueryTestItem, 0)

	err = cli.Find(context.Background(), filter1).Skip(-3).All(&res)
	ast.Error(err)
	ast.Empty(res)
}

func TestQuery_Limit(t *testing.T) {
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
		bson.M{"_id": id1, "name": "Alice", "age": 18},
		bson.M{"_id": id2, "name": "Alice", "age": 19},
		bson.M{"_id": id3, "name": "Lucas", "age": 20},
		bson.M{"_id": id4, "name": "Lucas", "age": 21},
	}
	_, _ = cli.InsertMany(context.Background(), docs)

	var err error
	var res []QueryTestItem

	filter1 := bson.M{
		"name": "Alice",
	}

	err = cli.Find(context.Background(), filter1).Limit(1).All(&res)
	ast.NoError(err)
	ast.Equal(1, len(res))

	res = make([]QueryTestItem, 0)

	err = cli.Find(context.Background(), filter1).Limit(3).All(&res)
	ast.NoError(err)
	ast.Equal(2, len(res))

	res = make([]QueryTestItem, 0)
	var cursor CursorI

	cursor = cli.Find(context.Background(), filter1).Limit(-2).Cursor()
	ast.NoError(cursor.Err())
	ast.NotNil(cursor)
}

func TestQuery_Sort(t *testing.T) {
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
		bson.M{"_id": id1, "name": "Alice", "age": 18},
		bson.M{"_id": id2, "name": "Alice", "age": 19},
		bson.M{"_id": id3, "name": "Lucas", "age": 18},
		bson.M{"_id": id4, "name": "Lucas", "age": 19},
	}
	_, _ = cli.InsertMany(context.Background(), docs)

	var err error
	var res []QueryTestItem

	// Sort a single field in ascending order
	filter1 := bson.M{
		"name": "Alice",
	}

	err = cli.Find(context.Background(), filter1).Sort("age").All(&res)
	ast.NoError(err)
	ast.Equal(2, len(res))
	ast.Equal(id1, res[0].Id)
	ast.Equal(id2, res[1].Id)

	// Sort a single field in descending order
	err = cli.Find(context.Background(), filter1).Sort("-age").All(&res)
	ast.NoError(err)
	ast.Equal(2, len(res))
	ast.Equal(id2, res[0].Id)
	ast.Equal(id1, res[1].Id)

	// Sort a single field in descending order, and sort the other field in ascending order
	err = cli.Find(context.Background(), bson.M{}).Sort("-age", "+name").All(&res)
	ast.NoError(err)
	ast.Equal(4, len(res))
	ast.Equal(id2, res[0].Id)
	ast.Equal(id4, res[1].Id)
	ast.Equal(id1, res[2].Id)
	ast.Equal(id3, res[3].Id)

	// fields is ""，panic
	res = make([]QueryTestItem, 0)
	ast.Panics(func() {
		cli.Find(context.Background(), filter1).Sort("").All(&res)
	})

	// fields is empty, does not panic or error (#128)
	err = cli.Find(context.Background(), bson.M{}).Sort().All(&res)
	ast.NoError(err)
	ast.Equal(4, len(res))

}

func TestQuery_Distinct(t *testing.T) {
	ast := require.New(t)
	cli := initClient("test")
	defer cli.Close(context.Background())
	defer cli.DropCollection(context.Background())
	cli.EnsureIndexes(context.Background(), nil, []string{"name"})

	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	id3 := primitive.NewObjectID()
	id4 := primitive.NewObjectID()
	id5 := primitive.NewObjectID()
	id6 := primitive.NewObjectID()
	docs := []interface{}{
		bson.M{"_id": id1, "name": "Alice", "age": 18},
		bson.M{"_id": id2, "name": "Alice", "age": 19},
		bson.M{"_id": id3, "name": "Lucas", "age": 20},
		bson.M{"_id": id4, "name": "Lucas", "age": 21},
		bson.M{"_id": id5, "name": "Kitty", "age": 23, "detail": bson.M{"errInfo": "timeout", "extra": "i/o"}},
		bson.M{"_id": id6, "name": "Kitty", "age": "23", "detail": bson.M{"errInfo": "timeout", "extra": "i/o"}},
	}
	_, _ = cli.InsertMany(context.Background(), docs)

	var err error

	filter1 := bson.M{
		"name": "Lily",
	}
	var res1 []int32

	err = cli.Find(context.Background(), filter1).Distinct("age", &res1)
	ast.NoError(err)
	ast.Equal(0, len(res1))

	filter2 := bson.M{
		"name": "Alice",
	}
	var res2 []int32

	err = cli.Find(context.Background(), filter2).Distinct("age", &res2)
	ast.NoError(err)
	ast.Equal(2, len(res2))

	var res3 []int32

	err = cli.Find(context.Background(), filter2).Distinct("age", res3)
	ast.EqualError(err, ErrQueryNotSlicePointer.Error())

	var res4 int

	err = cli.Find(context.Background(), filter2).Distinct("age", &res4)
	ast.EqualError(err, ErrQueryNotSliceType.Error())

	var res5 []string

	err = cli.Find(context.Background(), filter2).Distinct("age", &res5)
	ast.EqualError(err, ErrQueryResultTypeInconsistent.Error())

	// different behavior with different version of mongod, v4.4.0 return err and v4.0.19 return nil
	//var res6 []int32
	//err = cli.Find(context.Background(), filter2).Distinct("", &res6)
	//ast.Error(err) // (Location40352) FieldPath cannot be constructed with empty string
	//ast.Equal(0, len(res6))

	var res7 []int32
	filter3 := 1

	err = cli.Find(context.Background(), filter3).Distinct("age", &res7)
	ast.Error(err)
	ast.Equal(0, len(res7))

	var res8 interface{}

	res8 = []string{}
	err = cli.Find(context.Background(), filter2).Distinct("age", &res8)
	ast.NoError(err)
	ast.NotNil(res8)

	res9, ok := res8.(primitive.A)
	ast.Equal(true, ok)
	ast.Len(res9, 2)

	filter4 := bson.M{}
	var res10 []int32
	err = cli.Find(context.Background(), filter4).Distinct("detail", &res10)
	ast.EqualError(err, ErrQueryResultTypeInconsistent.Error())

	type tmpStruct struct {
		ErrInfo string `bson:"errInfo"`
		Extra   string `bson:"extra"`
	}
	var res11 []tmpStruct
	err = cli.Find(context.Background(), filter4).Distinct("detail", &res11)
	ast.NoError(err)

	type tmpErrStruct struct {
		ErrInfo string    `bson:"errInfo"`
		Extra   time.Time `bson:"extra"`
	}
	var res12 []tmpErrStruct
	err = cli.Find(context.Background(), filter4).Distinct("detail", &res12)
	ast.EqualError(err, ErrQueryResultTypeInconsistent.Error())

	var res13 []int32
	err = cli.Find(context.Background(), filter4).Distinct("age", &res13)
	ast.EqualError(err, ErrQueryResultTypeInconsistent.Error())

	var res14 []interface{}
	err = cli.Find(context.Background(), filter4).Distinct("age", &res14)
	ast.NoError(err)
	ast.Len(res14, 6)
	for _, v := range res14 {
		switch v.(type) {
		case int32:
			fmt.Printf("int32 :%d\n", v)
		case string:
			fmt.Printf("string :%s\n", v)
		default:
			fmt.Printf("defalut err: %v %T\n", v, v)
		}
	}
}

func TestQuery_Select(t *testing.T) {
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
		bson.M{"_id": id1, "name": "Alice", "age": 18},
		bson.M{"_id": id2, "name": "Alice", "age": 19},
		bson.M{"_id": id3, "name": "Lucas", "age": 20},
		bson.M{"_id": id4, "name": "Lucas", "age": 21},
	}
	_, _ = cli.InsertMany(context.Background(), docs)

	var err error
	var res QueryTestItem

	filter1 := bson.M{
		"_id": id1,
	}
	projection1 := bson.M{
		"age": 1,
	}

	err = cli.Find(context.Background(), filter1).Select(projection1).One(&res)
	ast.NoError(err)
	ast.NotNil(res)
	ast.Equal("", res.Name)
	ast.Equal(18, res.Age)
	ast.Equal(id1, res.Id)

	res = QueryTestItem{}
	projection2 := bson.M{
		"age": 0,
	}

	err = cli.Find(context.Background(), filter1).Select(projection2).One(&res)
	ast.NoError(err)
	ast.NotNil(res)
	ast.Equal("Alice", res.Name)
	ast.Equal(0, res.Age)
	ast.Equal(id1, res.Id)

	res = QueryTestItem{}
	projection3 := bson.M{
		"_id": 0,
	}

	err = cli.Find(context.Background(), filter1).Select(projection3).One(&res)
	ast.NoError(err)
	ast.NotNil(res)
	ast.Equal("Alice", res.Name)
	ast.Equal(18, res.Age)
	ast.Equal(primitive.NilObjectID, res.Id)
}

func TestQuery_Cursor(t *testing.T) {
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
		bson.D{{"_id", id1}, {"name", "Alice"}, {"age", 18}},
		bson.D{{"_id", id2}, {"name", "Alice"}, {"age", 19}},
		bson.D{{"_id", id3}, {"name", "Lucas"}, {"age", 20}},
		bson.D{{"_id", id4}, {"name", "Lucas"}, {"age", 21}},
	}
	_, _ = cli.InsertMany(context.Background(), docs)

	var res QueryTestItem

	filter1 := bson.M{
		"name": "Alice",
	}
	projection1 := bson.M{
		"name": 0,
	}

	cursor := cli.Find(context.Background(), filter1).Select(projection1).Sort("age").Limit(2).Skip(1).Cursor()
	ast.NoError(cursor.Err())
	ast.NotNil(cursor)

	val := cursor.Next(&res)
	ast.Equal(true, val)
	ast.Equal(id2, res.Id)

	val = cursor.Next(&res)
	ast.Equal(false, val)

	filter2 := bson.M{
		"name": "Lily",
	}

	cursor = cli.Find(context.Background(), filter2).Cursor()
	ast.NoError(cursor.Err())
	ast.NotNil(cursor)

	res = QueryTestItem{}
	val = cursor.Next(&res)
	ast.Equal(false, val)
	ast.Empty(res)

	filter3 := 1

	cursor = cli.Find(context.Background(), filter3).Cursor()
	ast.Error(cursor.Err())
}

func TestQuery_Hint(t *testing.T) {
	ast := require.New(t)
	cli := initClient("test")
	defer cli.Close(context.Background())
	defer cli.DropCollection(context.Background())
	cli.EnsureIndexes(context.Background(), nil, []string{"name", "age"})

	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	id3 := primitive.NewObjectID()
	id4 := primitive.NewObjectID()
	docs := []interface{}{
		bson.M{"_id": id1, "name": "Alice", "age": 18},
		bson.M{"_id": id2, "name": "Alice", "age": 19},
		bson.M{"_id": id3, "name": "Lucas", "age": 20},
		bson.M{"_id": id4, "name": "Lucas", "age": 21},
	}
	_, _ = cli.InsertMany(context.Background(), docs)

	var err error
	var res []QueryTestItem

	filter1 := bson.M{
		"name": "Alice",
		"age":  18,
	}

	// index name as hint
	err = cli.Find(context.Background(), filter1).Hint("age_1").All(&res)
	ast.NoError(err)
	ast.Equal(1, len(res))

	// index name as hint
	var resOne QueryTestItem
	err = cli.Find(context.Background(), filter1).Hint("name_1").One(&resOne)
	ast.NoError(err)

	// not index name as hint
	err = cli.Find(context.Background(), filter1).Hint("age").All(&res)
	ast.Error(err)

	// nil hint
	err = cli.Find(context.Background(), filter1).Hint(nil).All(&res)
	ast.NoError(err)

}

func TestQuery_Apply(t *testing.T) {
	ast := require.New(t)
	cli := initClient("test")
	defer cli.Close(context.Background())
	defer cli.DropCollection(context.Background())
	cli.EnsureIndexes(context.Background(), nil, []string{"name"})

	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	id3 := primitive.NewObjectID()
	docs := []interface{}{
		bson.M{"_id": id1, "name": "Alice", "age": 18},
		bson.M{"_id": id2, "name": "Alice", "age": 19},
		bson.M{"_id": id3, "name": "Lucas", "age": 20, "instock": []bson.M{
			{"warehouse": "B", "qty": 15},
			{"warehouse": "C", "qty": 35},
			{"warehouse": "E", "qty": 15},
			{"warehouse": "F", "qty": 45},
		}}}

	_, _ = cli.InsertMany(context.Background(), docs)

	var err error
	res1 := QueryTestItem{}
	filter1 := bson.M{
		"name": "Tom",
	}
	change1 := Change{}

	err = cli.Find(context.Background(), filter1).Apply(change1, &res1)
	ast.EqualError(err, mongo.ErrNilDocument.Error())

	change1.Update = bson.M{
		operator.Set: bson.M{
			"name": "Tom",
			"age":  18,
		},
	}
	err = cli.Find(context.Background(), filter1).Apply(change1, &res1)
	ast.EqualError(err, mongo.ErrNoDocuments.Error())

	change1.ReturnNew = true
	err = cli.Find(context.Background(), filter1).Apply(change1, &res1)
	ast.EqualError(err, mongo.ErrNoDocuments.Error())

	change1.ReturnNew = false
	change1.Upsert = true
	err = cli.Find(context.Background(), filter1).Apply(change1, &res1)
	ast.NoError(err)
	ast.Equal("", res1.Name)
	ast.Equal(0, res1.Age)

	change1.Update = bson.M{
		operator.Set: bson.M{
			"name": "Tom",
			"age":  19,
		},
	}
	change1.ReturnNew = true
	change1.Upsert = true
	err = cli.Find(context.Background(), filter1).Apply(change1, &res1)
	ast.NoError(err)
	ast.Equal("Tom", res1.Name)
	ast.Equal(19, res1.Age)

	res2 := QueryTestItem{}
	filter2 := bson.M{
		"name": "Alice",
	}
	change2 := Change{
		ReturnNew: true,
		Update: bson.M{
			operator.Set: bson.M{
				"name": "Alice",
				"age":  22,
			},
		},
	}
	projection2 := bson.M{
		"age": 1,
	}
	err = cli.Find(context.Background(), filter2).Sort("age").Select(projection2).Apply(change2, &res2)
	ast.NoError(err)
	ast.Equal("", res2.Name)
	ast.Equal(22, res2.Age)

	res3 := QueryTestItem{}
	filter3 := bson.M{
		"name": "Bob",
	}
	change3 := Change{
		Remove: true,
	}
	err = cli.Find(context.Background(), filter3).Apply(change3, &res3)
	ast.EqualError(err, mongo.ErrNoDocuments.Error())

	res3 = QueryTestItem{}
	filter3 = bson.M{
		"name": "Alice",
	}
	projection3 := bson.M{
		"age": 1,
	}
	err = cli.Find(context.Background(), filter3).Sort("age").Select(projection3).Apply(change3, &res3)
	ast.NoError(err)
	ast.Equal("", res3.Name)
	ast.Equal(19, res3.Age)

	res4 := QueryTestItem{}
	filter4 := bson.M{
		"name": "Bob",
	}
	change4 := Change{
		Replace: true,
		Update: bson.M{
			operator.Set: bson.M{
				"name": "Bob",
				"age":  23,
			},
		},
	}
	err = cli.Find(context.Background(), filter4).Apply(change4, &res4)
	ast.EqualError(err, ErrReplacementContainUpdateOperators.Error())

	change4.Update = bson.M{"name": "Bob", "age": 23}
	err = cli.Find(context.Background(), filter4).Apply(change4, &res4)
	ast.EqualError(err, mongo.ErrNoDocuments.Error())

	change4.ReturnNew = true
	err = cli.Find(context.Background(), filter4).Apply(change4, &res4)
	ast.EqualError(err, mongo.ErrNoDocuments.Error())

	change4.Upsert = true
	change4.ReturnNew = true
	err = cli.Find(context.Background(), filter4).Apply(change4, &res4)
	ast.NoError(err)
	ast.Equal("Bob", res4.Name)
	ast.Equal(23, res4.Age)

	change4 = Change{
		Replace:   true,
		Update:    bson.M{"name": "Bob", "age": 25},
		Upsert:    true,
		ReturnNew: false,
	}
	projection4 := bson.M{
		"age":  1,
		"name": 1,
	}
	err = cli.Find(context.Background(), filter4).Sort("age").Select(projection4).Apply(change4, &res4)
	ast.NoError(err)
	ast.Equal("Bob", res4.Name)
	ast.Equal(23, res4.Age)

	res4 = QueryTestItem{}
	filter4 = bson.M{
		"name": "James",
	}
	change4 = Change{
		Replace:   true,
		Update:    bson.M{"name": "James", "age": 26},
		Upsert:    true,
		ReturnNew: false,
	}
	err = cli.Find(context.Background(), filter4).Apply(change4, &res4)
	ast.NoError(err)
	ast.Equal("", res4.Name)
	ast.Equal(0, res4.Age)

	var res5 = QueryTestItem{}
	filter5 := bson.M{"name": "Lucas"}
	change5 := Change{
		Update:    bson.M{"$set": bson.M{"instock.$[elem].qty": 100}},
		ReturnNew: true,
	}
	err = cli.Find(context.Background(), filter5).SetArrayFilters(&options.ArrayFilters{Filters: []interface{}{
		bson.M{"elem.warehouse": bson.M{"$in": []string{"C", "F"}}},
	}}).Apply(change5, &res5)
	ast.NoError(err)

	for _, item := range res5.Instock {
		switch item.Warehouse {
		case "C", "F":
			ast.Equal(100, item.Qty)
		case "B", "E":
			ast.Equal(15, item.Qty)
		}
	}
}

func TestQuery_BatchSize(t *testing.T) {
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
		bson.M{"_id": id1, "name": "Alice", "age": 18},
		bson.M{"_id": id2, "name": "Alice", "age": 19},
		bson.M{"_id": id3, "name": "Lucas", "age": 20},
		bson.M{"_id": id4, "name": "Lucas", "age": 21},
	}
	_, _ = cli.InsertMany(context.Background(), docs)
	var res []QueryTestItem

	err := cli.Find(context.Background(), bson.M{"name": "Alice"}).BatchSize(1).All(&res)
	ast.NoError(err)
	ast.Len(res, 2)

}
