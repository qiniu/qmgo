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
	"testing"

	"github.com/qiniu/qmgo/operator"
	"github.com/qiniu/qmgo/options"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	opts "go.mongodb.org/mongo-driver/mongo/options"
)

const (
	URI      = "mongodb://localhost:27017"
	DATABASE = "class"
	COLL     = "user"
)

type UserInfo struct {
	Id     primitive.ObjectID `bson:"_id"`
	Name   string             `bson:"name"`
	Age    uint16             `bson:"age"`
	Weight uint32             `bson:"weight"`
}

var userInfo = UserInfo{
	Id:     NewObjectID(),
	Name:   "xm",
	Age:    7,
	Weight: 40,
}

var userInfos = []UserInfo{
	{Id: NewObjectID(), Name: "a1", Age: 6, Weight: 20},
	{Id: NewObjectID(), Name: "b2", Age: 6, Weight: 25},
	{Id: NewObjectID(), Name: "c3", Age: 6, Weight: 30},
	{Id: NewObjectID(), Name: "d4", Age: 6, Weight: 35},
	{Id: NewObjectID(), Name: "a1", Age: 7, Weight: 40},
	{Id: NewObjectID(), Name: "a1", Age: 8, Weight: 45},
}

var poolMonitor = &event.PoolMonitor{
	Event: func(evt *event.PoolEvent) {
		switch evt.Type {
		case event.GetSucceeded:
		case event.ConnectionReturned:
		}
	},
}

func TestQmgo(t *testing.T) {
	ast := require.New(t)
	ctx := context.Background()

	// create connect
	opt := opts.Client().SetAppName("example")
	cli, err := Open(ctx, &Config{Uri: URI, Database: DATABASE, Coll: COLL}, options.ClientOptions{ClientOptions: opt})

	ast.Nil(err)
	defer func() {
		if err = cli.Close(ctx); err != nil {
			panic(err)
		}
	}()
	defer cli.DropDatabase(ctx)

	cli.EnsureIndexes(ctx, []string{}, []string{"age", "name,weight"})
	// insert one document
	_, err = cli.InsertOne(ctx, userInfo)
	ast.Nil(err)

	// find one document
	one := UserInfo{}
	err = cli.Find(ctx, bson.M{"name": userInfo.Name}).One(&one)
	ast.Nil(err)
	ast.Equal(userInfo, one)

	// multiple insert
	_, err = cli.Collection.InsertMany(ctx, userInfos)
	ast.Nil(err)

	// find all 、sort and limit
	batch := []UserInfo{}
	cli.Find(ctx, bson.M{"age": 6}).Sort("weight").Limit(7).All(&batch)
	ast.Equal(4, len(batch))

	count, err := cli.Find(ctx, bson.M{"age": 6}).Count()
	ast.NoError(err)
	ast.Equal(int64(4), count)

	// aggregate
	matchStage := bson.D{{operator.Match, []bson.E{{"weight", bson.D{{operator.Gt, 30}}}}}}
	groupStage := bson.D{{operator.Group, bson.D{{"_id", "$name"}, {"total", bson.D{{operator.Sum, "$age"}}}}}}
	var showsWithInfo []bson.M
	err = cli.Aggregate(context.Background(), Pipeline{matchStage, groupStage}).All(&showsWithInfo)
	ast.Equal(3, len(showsWithInfo))
	for _, v := range showsWithInfo {
		if "a1" == v["_id"] {
			ast.Equal(int32(15), v["total"])
			continue
		}
		if "d4" == v["_id"] {
			ast.Equal(int32(6), v["total"])
			continue
		}
		ast.Error(errors.New("error"), "impossible")
	}
	// Update one
	err = cli.UpdateOne(ctx, bson.M{"name": "d4"}, bson.M{"$set": bson.M{"age": 17}})
	ast.NoError(err)
	cli.Find(ctx, bson.M{"age": 17}).One(&one)
	ast.Equal("d4", one.Name)
	// UpdateAll
	result, err := cli.UpdateAll(ctx, bson.M{"age": 6}, bson.M{"$set": bson.M{"age": 10}})
	ast.NoError(err)
	count, err = cli.Find(ctx, bson.M{"age": 10}).Count()
	ast.NoError(err)
	ast.Equal(result.ModifiedCount, count)
	// select
	one = UserInfo{}
	err = cli.Find(ctx, bson.M{"age": 10}).Select(bson.M{"age": 1}).One(&one)
	ast.NoError(err)
	ast.Equal(10, int(one.Age))
	ast.Equal("", one.Name)
	// remove
	err = cli.Remove(ctx, bson.M{"age": 7})
	ast.Nil(err)
}
