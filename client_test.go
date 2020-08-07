package qmgo

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func initClient(col string) *QmgoClient {
	cfg := Config{
		Uri:      "mongodb://localhost:27017",
		Database: "mongoxtest",
		Coll:     col,
	}
	var cTimeout int64 = 0
	var sTimeout int64 = 500000
	var maxPoolSize uint64 = 3000
	cfg.ConnectTimeoutMS = &cTimeout
	cfg.SocketTimeoutMS = &sTimeout
	cfg.MaxPoolSize = &maxPoolSize
	cli, err := Open(context.Background(), &cfg)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return cli
}

func TestQmgoClient(t *testing.T) {
	ast := require.New(t)
	var timeout int64 = 50

	// uri 错误
	cfg := Config{
		Uri:              "://127.0.0.1",
		ConnectTimeoutMS: &timeout,
	}

	var err error
	_, err = Open(context.Background(), &cfg)
	ast.NotNil(err)

	// Open 成功
	var maxPoolSize uint64 = 100

	cfg = Config{
		Uri:              "mongodb://localhost:27017",
		Database:         "mongoxtest",
		Coll:             "testopen",
		ConnectTimeoutMS: &timeout,
		MaxPoolSize:      &maxPoolSize,
	}

	cli, err := Open(context.Background(), &cfg)
	ast.NoError(err)
	ast.Equal(cli.GetDatabaseName(), "mongoxtest")
	ast.Equal(cli.GetCollectionName(), "testopen")

	err = cli.Ping(5)
	ast.NoError(err)

	res, err := cli.InsertOne(context.Background(), bson.D{{Key: "x", Value: 1}})
	ast.NoError(err)
	ast.NotNil(res)

	cli.DropCollection(context.Background())

	// close Client
	cli.Close(context.TODO())
	_, err = cli.InsertOne(context.Background(), bson.D{{Key: "x", Value: 1}})
	ast.EqualError(err, "client is disconnected")

	err = cli.Ping(5)
	ast.Error(err)
}

func TestIsDup(t *testing.T) {
	ast := require.New(t)

	ast.Equal(false, IsDup(errors.New("invaliderror")))
	ast.Equal(true, IsDup(errors.New("E11000")))
}

func TestClient(t *testing.T) {
	ast := require.New(t)

	var maxPoolSize uint64 = 100
	var timeout int64 = 50

	cfg := &Config{
		Uri:              "mongodb://localhost:27017",
		ConnectTimeoutMS: &timeout,
		MaxPoolSize:      &maxPoolSize,
	}

	c, err := NewClient(context.Background(), cfg)
	ast.Equal(nil, err)
	coll := c.Database("mongoxtest").Collection("testopen")

	res, err := coll.InsertOne(context.Background(), bson.D{{Key: "x", Value: 1}})
	ast.NoError(err)
	ast.NotNil(res)
	coll.DropCollection(context.Background())
}
