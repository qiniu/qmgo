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
	"errors"
	"fmt"
	"testing"

	opts "github.com/qiniu/qmgo/options"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

	fn := func(sCtx context.Context) (interface{}, error) {
		if _, err := cli.InsertOne(sCtx, bson.D{{"abc", int32(1)}}); err != nil {
			return nil, err
		}
		if _, err := cli.InsertOne(sCtx, bson.D{{"xyz", int32(999)}}); err != nil {
			return nil, err
		}
		return nil, nil
	}
	tops := options.Transaction()
	op := &opts.TransactionOptions{TransactionOptions: tops}
	_, err := cli.DoTransaction(ctx, fn, op)
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
	sOpts := options.Session().SetSnapshot(false)
	o := &opts.SessionOptions{sOpts}
	s, err := cli.Session(o)
	ast.NoError(err)
	ctx := context.Background()
	defer s.EndSession(ctx)

	callback := func(sCtx context.Context) (interface{}, error) {
		if _, err := cli.InsertOne(sCtx, bson.D{{"abc", int32(1)}}); err != nil {
			return nil, err
		}
		if _, err := cli.InsertOne(sCtx, bson.D{{"xyz", int32(999)}}); err != nil {
			return nil, err
		}
		err = s.AbortTransaction(sCtx)

		return nil, nil
	}

	_, err = s.StartTransaction(ctx, callback)
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
