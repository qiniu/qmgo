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
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCursor(t *testing.T) {
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
	_, err := cli.InsertMany(context.Background(), docs)
	ast.NoError(err)

	var res QueryTestItem

	// if query has 1 record，cursor can run Next one time， Next time return false
	filter1 := bson.M{
		"name": "Alice",
	}
	projection1 := bson.M{
		"name": 0,
	}

	cursor := cli.Find(context.Background(), filter1).Select(projection1).Sort("age").Limit(2).Skip(1).Cursor()
	ast.NoError(cursor.Err())

	val := cursor.Next(&res)
	ast.Equal(true, val)
	ast.Equal(id2, res.Id)

	val = cursor.Next(&res)
	ast.Equal(false, val)

	cursor.Close()

	// cursor ALL
	cursor = cli.Find(context.Background(), filter1).Select(projection1).Sort("age").Limit(2).Cursor()
	ast.NoError(cursor.Err())

	var results []QueryTestItem
	cursor.All(&results)
	ast.Equal(2, len(results))
	// can't match record, cursor run Next and return false
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

	cursor.Close()

	//  1 record，after cursor close，Next return false
	cursor = cli.Find(context.Background(), filter1).Select(projection1).Sort("age").Limit(2).Skip(1).Cursor()
	ast.NoError(cursor.Err())
	ast.NotNil(cursor)

	cursor.Close()

	ast.Equal(false, cursor.Next(&res))
	ast.NoError(cursor.Err())

	// generate Cursor with err
	cursor = cli.Find(context.Background(), 1).Select(projection1).Sort("age").Limit(2).Skip(1).Cursor()
	ast.Error(cursor.Err())
	//ast.Equal(int64(0), cursor.ID())
	ast.Error(cursor.All(&res))
	ast.Error(cursor.Close())
	ast.Equal(false, cursor.Next(&res))
}
