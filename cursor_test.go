package mongox

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCursor(t *testing.T) {
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

	// 查询结果集只有 1 条记录，cursor 只能 Next 一次，第二次 Next 即返回 false
	filter1 := bson.M{
		"name": "Alice",
	}
	projection1 := bson.M{
		"name": 0,
	}

	cursor, err := coll.Find(filter1).Select(projection1).Sort("age").Limit(2).Skip(1).Cursor()
	ast.NoError(err)

	val := cursor.Next(&res)
	ast.Equal(true, val)
	ast.Equal(id2, res.Id)

	val = cursor.Next(&res)
	ast.Equal(false, val)

	cursor.Close()

	// 未查询到匹配记录, cursor 第一次执行 Next 就返回 false
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

	cursor.Close()

	// 查询结果集有 1 条记录，cursor Close后，Next 返回 false
	cursor, err = coll.Find(filter1).Select(projection1).Sort("age").Limit(2).Skip(1).Cursor()
	ast.NoError(err)
	ast.NotNil(cursor)

	cursor.Close()

	ast.Equal(false, cursor.Next(&res))
	ast.NoError(cursor.Err())
}
