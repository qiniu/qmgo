package mongox

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func initClient(col string) *Client {
	cfg := Config{
		Uri:         "mongodb://localhost:27017",
		DB:          "mongoxtest",
		Coll:        col,
		Timeout:     5000,
		MaxPollSize: 100,
	}

	cli, err := Open(&cfg)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return cli
}

func TestClient(t *testing.T) {
	ast := require.New(t)

	// uri 错误
	cfg := Config{
		Uri:     "://127.0.0.1",
		Timeout: 1000,
	}

	var err error
	_, err = Open(&cfg)
	ast.NotNil(err)

	// Open 成功
	var cli *Client
	cfg = Config{
		Uri:         "mongodb://localhost:27017",
		DB:          "mongoxtest",
		Coll:        "testopen",
		Timeout:     5000,
		MaxPollSize: 100,
	}

	cli, err = Open(&cfg)
	ast.NoError(err)
	ast.Equal(cli.GetDatabaseName(), "mongoxtest")
	ast.Equal(cli.GetCollectionName(), "testopen")

	err = cli.Ping(5)
	ast.NoError(err)

	coll := cli.GetCollection(context.TODO())
	res, err := coll.InsertOne(bson.D{{Key: "x", Value: 1}})
	ast.NoError(err)
	ast.NotNil(res)

	coll.DropCollection()

	// 关闭 Client
	cli.Close(context.TODO())
	_, err = coll.InsertOne(bson.D{{Key: "x", Value: 1}})
	ast.EqualError(err, "client is disconnected")

	err = cli.Ping(5)
	ast.Error(err)
}

func TestClient_GetCollection(t *testing.T) {
	ast := require.New(t)

	var cli *Client
	var coll *Collection

	cli = initClient("test")
	coll = cli.GetCollection(context.TODO())
	coll.DropCollection()

	// 参数 ctx 为 nil
	collNew := cli.GetCollection(nil)
	ast.NotNil(collNew.Ctx)

	// 参数 ctx 不为 nil，检测 context 是否被替换, 并测试 collNew 是否正确
	key := "mongox_key"
	val := "mongox_value"
	ctxNew := context.WithValue(context.Background(), key, val)
	collNew = cli.GetCollection(ctxNew)

	v := collNew.Ctx.Value(key)
	ast.NotNil(v)
	ast.Equal(v, val)

	res, err := collNew.InsertOne(bson.M{"_id": primitive.NewObjectID(), "name": "Alice"})
	ast.NoError(err)
	ast.NotEmpty(res)
}
