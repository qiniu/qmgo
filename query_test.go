package qmgo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QueryTestItem struct {
	Id   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
	Age  int                `bson:"age"`
}

type QueryTestItem2 struct {
	Class string `bson:"class"`
}

func TestQuery_One(t *testing.T) {
	ast := require.New(t)

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
}

func TestQuery_Distinct(t *testing.T) {
	ast := require.New(t)

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

	//var res6 []int32
	//
	//err = cli.Find(context.Background(), filter2).Distinct("", &res6)
	//ast.NoError(err)
	//ast.Equal(0, len(res6))

	var res7 []int32
	filter3 := 1

	err = cli.Find(context.Background(), filter3).Distinct("age", &res7)
	ast.Error(err)
	ast.Equal(0, len(res7))

	var res8 interface{}

	res8 = []int32{}
	err = cli.Find(context.Background(), filter2).Distinct("age", &res8)
	ast.NoError(err)
	ast.NotNil(res8)

	res9, ok := res8.([]int32)
	ast.Equal(true, ok)
	ast.Len(res9, 2)
}

func TestQuery_Select(t *testing.T) {
	ast := require.New(t)

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
