package qmgo

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/description"
	"testing"
	"time"
)

func initTransactionClient(coll string) *QmgoClient {
	cfg := Config{
		Uri:      "mongodb://localhost:27017",
		Database: "transaction",
		Coll:     coll,
	}
	var cTimeout int64 = 0
	var sTimeout int64 = 500000
	var maxPoolSize uint64 = 30000
	var minPoolSize uint64 = 0
	cfg.ConnectTimeoutMS = &cTimeout
	cfg.SocketTimeoutMS = &sTimeout
	cfg.MaxPoolSize = &maxPoolSize
	cfg.MinPoolSize = &minPoolSize
	cfg.ReadPreference = &ReadPref{Mode: readpref.PrimaryMode}
	qClient, err := Open(context.Background(), &cfg)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	qClient.InsertOne(context.Background(), bson.M{"name": "before_transaction"})
	return qClient

}
func TestClient_DoTransaction(t *testing.T) {
	ast := require.New(t)
	ctx := context.Background()
	cli := initTransactionClient("test")
	defer cli.DropDatabase(ctx)

	if !okRunTransaction() {
		t.Skip("can't run transaction")
	}

	fn := func(sCtx context.Context) (interface{}, error) {
		if _, err := cli.InsertOne(sCtx, bson.D{{"abc", int32(1)}}); err != nil {
			return nil, err
		}
		if _, err := cli.InsertOne(sCtx, bson.D{{"xyz", int32(999)}}); err != nil {
			return nil, err
		}
		return nil, nil
	}
	_, err := cli.DoTransaction(ctx, fn)
	ast.NoError(err)
	r := bson.M{}
	cli.Find(ctx, bson.M{"abc": 1}).One(&r)
	ast.Equal(r["abc"], int32(1))

	cli.Find(ctx, bson.M{"xyz": 999}).One(&r)
	ast.Equal(r["xyz"], int32(999))
}

func TestSession_AbortTransaction(t *testing.T) {
	ast := require.New(t)
	cli := initTransactionClient("test")

	defer cli.DropCollection(context.Background())
	s, err := cli.Session()
	ast.NoError(err)
	ctx := context.Background()
	defer s.EndSession(ctx)

	if !okRunTransaction() {
		t.Skip("can't run transaction")
	}

	callback := func(sCtx context.Context) (interface{}, error) {
		if _, err := cli.InsertOne(sCtx, bson.D{{"abc", int32(1)}}); err != nil {
			return nil, err
		}
		if _, err := cli.InsertOne(sCtx, bson.D{{"xyz", int32(999)}}); err != nil {
			return nil, err
		}
		time.Sleep(5 * time.Second)
		return nil, nil
	}
	go func() {
		time.Sleep(3 * time.Second)
		// abort the already worked operation, can't abort the later operation
		// it seems a mongodb-go-driver bug
		err = s.AbortTransaction(ctx)
	}()
	_, err = s.StartTransaction(ctx, callback)
	ast.NoError(err)

	ast.NoError(err)
	r := bson.M{}
	err = cli.Find(ctx, bson.M{"abc": 1}).One(&r)
	ast.Error(err)
	// abort the already worked operation, can't abort the later operation
	// it seems a mongodb-go-driver bug
	err = cli.Find(ctx, bson.M{"xyz": 999}).One(&r)
	ast.Error(err)
}

func TestSession_Cancel(t *testing.T) {
	ast := require.New(t)
	cli := initTransactionClient("test")

	defer cli.DropCollection(context.Background())
	s, err := cli.Session()
	ast.NoError(err)
	ctx := context.Background()
	defer s.EndSession(ctx)
	if !okRunTransaction() {
		t.Skip("can't run transaction")
	}
	callback := func(sCtx context.Context) (interface{}, error) {
		if _, err := cli.InsertOne(sCtx, bson.D{{"abc", int32(1)}}); err != nil {
			return nil, err
		}
		if _, err := cli.InsertOne(sCtx, bson.D{{"xyz", int32(999)}}); err != nil {
			return nil, err
		}
		return nil, errors.New("cancel operations")
	}
	_, err = s.StartTransaction(ctx, callback)
	ast.Error(err)
	r := bson.M{}
	err = cli.Find(ctx, bson.M{"abc": 1}).One(&r)
	ast.True(IsErrNoDocuments(err))
	err = cli.Find(ctx, bson.M{"xyz": 999}).One(&r)
	ast.True(IsErrNoDocuments(err))
}

func TestSession_RetryTransAction(t *testing.T) {
	ast := require.New(t)
	cli := initTransactionClient("test")
	defer cli.DropCollection(context.Background())
	s, err := cli.Session()
	ast.NoError(err)
	ctx := context.Background()
	defer s.EndSession(ctx)
	if !okRunTransaction() {
		t.Skip("can't run transaction")
	}
	count := 0
	callback := func(sCtx context.Context) (interface{}, error) {
		if _, err := cli.InsertOne(sCtx, bson.D{{"abc", int32(1)}}); err != nil {
			return nil, err
		}
		if _, err := cli.InsertOne(sCtx, bson.D{{"xyz", int32(999)}}); err != nil {
			return nil, err
		}
		if count == 0 {
			count++
			return nil, ErrTransactionRetry
		}
		return nil, nil
	}
	_, err = s.StartTransaction(ctx, callback)
	ast.NoError(err)
	r := bson.M{}
	cli.Find(ctx, bson.M{"abc": 1}).One(&r)
	ast.Equal(r["abc"], int32(1))
	cli.Find(ctx, bson.M{"xyz": 999}).One(&r)
	ast.Equal(r["xyz"], int32(999))
	ast.Equal(count, 1)
}

func okRunTransaction() bool {
	vr, err := CompareVersions("4.0", cli.ServerVersion())
	if err != nil {
		return false
	}
	if vr > 0 {
		return false
	}
	topo, err := cli.topology()
	if topo == description.Single {
		return false
	}
	return true
}
