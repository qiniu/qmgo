package mongox

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
	var res QueryTestItem

	// 成功查询到一条记录
	filter1 := bson.M{
		"name": "Alice",
	}
	projection1 := bson.M{
		"age": 0,
	}

	err = coll.Find(filter1).Select(projection1).Sort("age").Limit(1).Skip(1).One(&res)
	ast.Nil(err)
	ast.Equal(id2, res.Id)
	ast.Equal("Alice", res.Name)

	// 未查询到匹配记录
	res = QueryTestItem{}
	filter2 := bson.M{
		"name": "Lily",
	}

	err = coll.Find(filter2).One(&res)
	ast.Error(err)
	ast.Empty(res)

	// filter 为 bson.M{}，表示匹配全部，会只返回一条记录
	res = QueryTestItem{}
	filter3 := bson.M{}

	err = coll.Find(filter3).One(&res)
	ast.NoError(err)
	ast.NotEmpty(res)

	// filter 为 nil，报错
	res = QueryTestItem{}
	err = coll.Find(nil).One(&res)
	ast.Error(err)
	ast.Empty(res)

	// res 为 nil 或者 无法解析类型
	err = coll.Find(filter1).One(nil)
	ast.Error(err)

	var tv int
	err = coll.Find(filter1).One(&tv)
	ast.Error(err)

	// res 为可解析对象，但 bson tag 与 mongodb 记录不一致，不报错，res 为该数据结构初始化状态
	var tt QueryTestItem2
	err = coll.Find(filter1).One(&tt)
	ast.NoError(err)
	ast.Empty(tt)
}

func TestQuery_All(t *testing.T) {
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
	docs := []bson.M{
		{"_id": id1, "name": "Alice", "age": 18},
		{"_id": id2, "name": "Alice", "age": 19},
		{"_id": id3, "name": "Lucas", "age": 20},
		{"_id": id4, "name": "Lucas", "age": 21},
	}
	_, _ = coll.InsertMany(docs)

	var err error
	var res []QueryTestItem

	// 成功查询到多条记录
	filter1 := bson.M{
		"name": "Alice",
	}
	projection1 := bson.M{
		"name": 0,
	}

	err = coll.Find(filter1).Select(projection1).Sort("age").Limit(2).Skip(1).All(&res)
	ast.NoError(err)
	ast.Equal(1, len(res))

	// 未查询到匹配记录, res 为空
	res = make([]QueryTestItem, 0)
	filter2 := bson.M{
		"name": "Lily",
	}

	err = coll.Find(filter2).All(&res)
	ast.NoError(err)
	ast.Empty(res)

	// filter 为 bson.M{}，表示匹配全部，会返回集合里所有的记录
	res = make([]QueryTestItem, 0)
	filter3 := bson.M{}

	err = coll.Find(filter3).All(&res)
	ast.NoError(err)
	ast.Equal(4, len(res))

	// filter 为 nil，报错
	res = make([]QueryTestItem, 0)
	err = coll.Find(nil).All(&res)
	ast.Error(err)
	ast.Empty(res)

	// res 为 nil，报错
	err = coll.Find(filter1).All(nil)
	ast.Error(err)

	// res 无法解析类型，会导致panic
	var tv int
	ast.Panics(func() {
		coll.Find(filter1).All(&tv)
	})

	// res 为可解析对象，但 bson tag 与 mongodb 记录不一致，不报错
	// 会根据 res 数据结构的 bson tag 来映射对应的值，没有的 tag 则值为对应类型的默认值
	// res 长度则为 filter 条件筛选出来的记录数
	var tt []QueryTestItem2
	err = coll.Find(filter1).All(&tt)
	ast.NoError(err)
	ast.Equal(2, len(tt))
}

func TestQuery_Count(t *testing.T) {
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
	docs := []bson.M{
		{"_id": id1, "name": "Alice", "age": 18},
		{"_id": id2, "name": "Alice", "age": 19},
		{"_id": id3, "name": "Lucas", "age": 20},
		{"_id": id4, "name": "Lucas", "age": 21},
	}
	_, _ = coll.InsertMany(docs)

	var err error
	var cnt int64

	// filter 能过滤到记录
	filter1 := bson.M{
		"name": "Alice",
	}

	cnt, err = coll.Find(filter1).Limit(2).Skip(1).Count()
	ast.NoError(err)
	ast.Equal(int64(1), cnt)

	// 未查询到匹配记录, cnt 为 0
	filter2 := bson.M{
		"name": "Lily",
	}

	cnt, err = coll.Find(filter2).Count()
	ast.NoError(err)
	ast.Zero(cnt)

	// filter 为 bson.M{}，表示匹配全部，会统计集合里所有的记录
	filter3 := bson.M{}

	cnt, err = coll.Find(filter3).Count()
	ast.NoError(err)
	ast.Equal(int64(4), cnt)

	// filter 为 nil，报错
	cnt, err = coll.Find(nil).Count()
	ast.Error(err)
	ast.Zero(cnt)
}

func TestQuery_Skip(t *testing.T) {
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
	docs := []bson.M{
		{"_id": id1, "name": "Alice", "age": 18},
		{"_id": id2, "name": "Alice", "age": 19},
		{"_id": id3, "name": "Lucas", "age": 20},
		{"_id": id4, "name": "Lucas", "age": 21},
	}
	_, _ = coll.InsertMany(docs)

	var err error
	var res []QueryTestItem

	// filter 可匹配到记录，跳过 1 条记录，返回剩余的记录
	filter1 := bson.M{
		"name": "Alice",
	}

	err = coll.Find(filter1).Skip(1).All(&res)
	ast.NoError(err)
	ast.Equal(1, len(res))

	// filter 可匹配到记录，跳过数大于已有记录总数，res 返回空
	res = make([]QueryTestItem, 0)

	err = coll.Find(filter1).Skip(3).All(&res)
	ast.NoError(err)
	ast.Empty(res)

	// skip 传入负数，报错
	res = make([]QueryTestItem, 0)

	err = coll.Find(filter1).Skip(-3).All(&res)
	ast.Error(err)
	ast.Empty(res)
}

func TestQuery_Limit(t *testing.T) {
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
	docs := []bson.M{
		{"_id": id1, "name": "Alice", "age": 18},
		{"_id": id2, "name": "Alice", "age": 19},
		{"_id": id3, "name": "Lucas", "age": 20},
		{"_id": id4, "name": "Lucas", "age": 21},
	}
	_, _ = coll.InsertMany(docs)

	var err error
	var res []QueryTestItem

	// filter 可匹配到多条记录，limit置1，返回1条记录
	filter1 := bson.M{
		"name": "Alice",
	}

	err = coll.Find(filter1).Limit(1).All(&res)
	ast.NoError(err)
	ast.Equal(1, len(res))

	// filter 可匹配到记录，limit数大于已有记录总数，res 最大已有记录数
	res = make([]QueryTestItem, 0)

	err = coll.Find(filter1).Limit(3).All(&res)
	ast.NoError(err)
	ast.Equal(2, len(res))

	// todo limit 传入负数，limit 类似整数，但是会关闭光标
	res = make([]QueryTestItem, 0)
	var cursor CursorI

	cursor, err = coll.Find(filter1).Limit(-2).Cursor()
	ast.NoError(err)
	ast.NotNil(cursor)
}

func TestQuery_Sort(t *testing.T) {
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
	docs := []bson.M{
		{"_id": id1, "name": "Alice", "age": 18},
		{"_id": id2, "name": "Alice", "age": 19},
		{"_id": id3, "name": "Lucas", "age": 18},
		{"_id": id4, "name": "Lucas", "age": 19},
	}
	_, _ = coll.InsertMany(docs)

	var err error
	var res []QueryTestItem

	// 对单字段升序排列
	filter1 := bson.M{
		"name": "Alice",
	}

	err = coll.Find(filter1).Sort("age").All(&res)
	ast.NoError(err)
	ast.Equal(2, len(res))
	ast.Equal(id1, res[0].Id)
	ast.Equal(id2, res[1].Id)

	// 对单字段降序排列
	err = coll.Find(filter1).Sort("-age").All(&res)
	ast.NoError(err)
	ast.Equal(2, len(res))
	ast.Equal(id2, res[0].Id)
	ast.Equal(id1, res[1].Id)

	// 对某个单字段降序排列, 另外一个字段按升序排列
	err = coll.Find(bson.M{}).Sort("-age", "+name").All(&res)
	ast.NoError(err)
	ast.Equal(4, len(res))
	ast.Equal(id2, res[0].Id)
	ast.Equal(id4, res[1].Id)
	ast.Equal(id1, res[2].Id)
	ast.Equal(id3, res[3].Id)

	// fields 为 ""，panic
	res = make([]QueryTestItem, 0)
	ast.Panics(func() {
		coll.Find(filter1).Sort("").All(&res)
	})
}

func TestQuery_Distinct(t *testing.T) {
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
	docs := []bson.M{
		{"_id": id1, "name": "Alice", "age": 18},
		{"_id": id2, "name": "Alice", "age": 19},
		{"_id": id3, "name": "Lucas", "age": 20},
		{"_id": id4, "name": "Lucas", "age": 21},
	}
	_, _ = coll.InsertMany(docs)

	var err error

	// 未查询到匹配记录, res 为空
	filter1 := bson.M{
		"name": "Lily",
	}
	var res1 []int32

	err = coll.Find(filter1).Distinct("age", &res1)
	ast.NoError(err)
	ast.Equal(0, len(res1))

	// 成功匹配到记录
	filter2 := bson.M{
		"name": "Alice",
	}
	var res2 []int32

	err = coll.Find(filter2).Distinct("age", &res2)
	ast.NoError(err)
	ast.Equal(2, len(res2))

	// result 参数不为 slice 指针 或 不为 slice
	var res3 []int32

	err = coll.Find(filter2).Distinct("age", res3)
	ast.EqualError(err, ERR_QUERY_NOT_SLICE_POINTER.Error())

	var res4 int

	err = coll.Find(filter2).Distinct("age", &res4)
	ast.EqualError(err, ERR_QUERY_NOT_SLICE_TYPE.Error())

	// result 参数为 slice 指针，但是数据类型与 mongodb 里不一致
	var res5 []string

	err = coll.Find(filter2).Distinct("age", &res5)
	ast.EqualError(err, ERR_QUERY_RESULT_TYPE_INCONSISTEN.Error())

	// key 为空字符串，不报错，res 为空
	var res6 []int32

	err = coll.Find(filter2).Distinct("", &res6)
	ast.NoError(err)
	ast.Equal(0, len(res6))

	// filter 语法错误，引发 error
	var res7 []int32
	filter3 := 1

	err = coll.Find(filter3).Distinct("age", &res7)
	ast.Error(err)
	ast.Equal(0, len(res7))

	// result 静态类型为 interface
	var res8 interface{}

	res8 = []int32{}
	err = coll.Find(filter2).Distinct("age", &res8)
	ast.NoError(err)
	ast.NotNil(res8)

	res9, ok := res8.([]int32)
	ast.Equal(true, ok)
	ast.Len(res9, 2)
}

func TestQuery_Select(t *testing.T) {
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
	docs := []bson.M{
		{"_id": id1, "name": "Alice", "age": 18},
		{"_id": id2, "name": "Alice", "age": 19},
		{"_id": id3, "name": "Lucas", "age": 20},
		{"_id": id4, "name": "Lucas", "age": 21},
	}
	_, _ = coll.InsertMany(docs)

	var err error
	var res QueryTestItem

	// 只显示 age 字段
	filter1 := bson.M{
		"_id": id1,
	}
	projection1 := bson.M{
		"age": 1,
	}

	err = coll.Find(filter1).Select(projection1).One(&res)
	ast.NoError(err)
	ast.NotNil(res)
	ast.Equal("", res.Name)
	ast.Equal(18, res.Age)
	ast.Equal(id1, res.Id)

	// 不显示 age 字段
	res = QueryTestItem{}
	projection2 := bson.M{
		"age": 0,
	}

	err = coll.Find(filter1).Select(projection2).One(&res)
	ast.NoError(err)
	ast.NotNil(res)
	ast.Equal("Alice", res.Name)
	ast.Equal(0, res.Age)
	ast.Equal(id1, res.Id)

	// 不显示 _id 字段
	res = QueryTestItem{}
	projection3 := bson.M{
		"_id": 0,
	}

	err = coll.Find(filter1).Select(projection3).One(&res)
	ast.NoError(err)
	ast.NotNil(res)
	ast.Equal("Alice", res.Name)
	ast.Equal(18, res.Age)
	ast.Equal(primitive.NilObjectID, res.Id)
}

func TestQuery_Cursor(t *testing.T) {
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
	docs := []bson.M{
		{"_id": id1, "name": "Alice", "age": 18},
		{"_id": id2, "name": "Alice", "age": 19},
		{"_id": id3, "name": "Lucas", "age": 20},
		{"_id": id4, "name": "Lucas", "age": 21},
	}
	_, _ = coll.InsertMany(docs)

	var err error
	var res QueryTestItem

	// 查询结果集只有 1 条记录，cursor 只能 Next 一次
	filter1 := bson.M{
		"name": "Alice",
	}
	projection1 := bson.M{
		"name": 0,
	}

	cursor, err := coll.Find(filter1).Select(projection1).Sort("age").Limit(2).Skip(1).Cursor()
	ast.NoError(err)
	ast.NotNil(cursor)

	val := cursor.Next(&res)
	ast.Equal(true, val)
	ast.Equal(id2, res.Id)

	val = cursor.Next(&res)
	ast.Equal(false, val)

	// 未查询到匹配记录
	filter2 := bson.M{
		"name": "Lily",
	}

	cursor, err = coll.Find(filter2).Cursor()
	ast.NoError(err)
	ast.NotNil(cursor)

	res = QueryTestItem{}
	val = cursor.Next(&res)
	ast.Equal(false, val)
	ast.Empty(res)

	// filter 语法错误，引发 error
	filter3 := 1

	cursor, err = coll.Find(filter3).Cursor()
	ast.Error(err)
	ast.Nil(cursor)
}
