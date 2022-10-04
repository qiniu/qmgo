package qmgo

import (
	"context"

	"testing"

	"github.com/qiniu/qmgo/operator"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestBulk(t *testing.T) {
	ast := require.New(t)
	cli := initClient("test")
	defer cli.Close(context.Background())
	defer cli.DropCollection(context.Background())

	id := primitive.NewObjectID()
	lucas := UserInfo{Id: primitive.NewObjectID(), Name: "Lucas", Age: 12}
	alias := UserInfo{Id: id, Name: "Alias", Age: 21}
	jess := UserInfo{Id: primitive.NewObjectID(), Name: "Jess", Age: 22}
	joe := UserInfo{Id: primitive.NewObjectID(), Name: "Joe", Age: 22}
	ethanId := primitive.NewObjectID()
	ethan := UserInfo{Id: ethanId, Name: "Ethan", Age: 8}

	result, err := cli.Bulk().
		InsertOne(lucas).InsertOne(alias).InsertOne(jess).
		UpdateOne(bson.M{"name": "Jess"}, bson.M{operator.Set: bson.M{"age": 23}}).UpdateId(id, bson.M{operator.Set: bson.M{"age": 23}}).
		UpdateAll(bson.M{"age": 23}, bson.M{operator.Set: bson.M{"age": 18}}).
		Upsert(bson.M{"age": 17}, joe).UpsertId(ethanId, ethan).
		Remove(bson.M{"name": "Joe"}).RemoveId(ethanId).RemoveAll(bson.M{"age": 18}).
		Run(context.Background())
	ast.NoError(err)
	ast.Equal(int64(3), result.InsertedCount)
	ast.Equal(int64(4), result.ModifiedCount)
	ast.Equal(int64(4), result.DeletedCount)
	ast.Equal(int64(2), result.UpsertedCount)
	ast.Equal(2, len(result.UpsertedIDs))
	ast.Equal(int64(4), result.MatchedCount)

}

func TestBulkUpsertOne(t *testing.T) {
	ast := require.New(t)
	cli := initClient("test")
	defer cli.Close(context.Background())
	defer cli.DropCollection(context.Background())

	result, err := cli.Bulk().
		UpsertOne(bson.M{"name": "Jess"}, bson.M{operator.Set: bson.M{"age": 20}, operator.SetOnInsert: bson.M{"weight": 40}}).
		UpsertOne(bson.M{"name": "Jess"}, bson.M{operator.Set: bson.M{"age": 30}, operator.SetOnInsert: bson.M{"weight": 40}}).
		Run(context.Background())

	ast.NoError(err)
	ast.Equal(int64(0), result.InsertedCount)
	ast.Equal(int64(1), result.ModifiedCount)
	ast.Equal(int64(0), result.DeletedCount)
	ast.Equal(int64(1), result.UpsertedCount)
	ast.Equal(1, len(result.UpsertedIDs))
	ast.Equal(int64(1), result.MatchedCount)
}
