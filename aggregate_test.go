package qmgo

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAggregate(t *testing.T) {
	ast := require.New(t)
	cli := initClient("test")
	defer cli.Close(context.Background())
	defer cli.DropCollection(context.Background())
	cli.EnsureIndexes(context.Background(), nil, []string{"name"})

	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	id3 := primitive.NewObjectID()
	id4 := primitive.NewObjectID()
	id5 := primitive.NewObjectID()
	docs := []interface{}{
		QueryTestItem{Id: id1, Name: "Alice", Age: 10},
		QueryTestItem{Id: id2, Name: "Alice", Age: 12},
		QueryTestItem{Id: id3, Name: "Lucas", Age: 33},
		QueryTestItem{Id: id4, Name: "Lucas", Age: 22},
		QueryTestItem{Id: id5, Name: "Lucas", Age: 44},
	}
	cli.InsertMany(context.Background(), docs)
	matchStage := bson.D{{"$match", []bson.E{{"age", bson.D{{"$gt", 11}}}}}}
	groupStage := bson.D{{"$group", bson.D{{"_id", "$name"}, {"total", bson.D{{"$sum", "$age"}}}}}}
	var showsWithInfo []bson.M
	// aggregate ALL()
	err := cli.Aggregate(context.Background(), Pipeline{matchStage, groupStage}).All(&showsWithInfo)
	ast.NoError(err)
	ast.Equal(2, len(showsWithInfo))
	for _, v := range showsWithInfo {
		if "Alice" == v["_id"] {
			ast.Equal(int32(12), v["total"])
			continue
		}
		if "Lucas" == v["_id"] {
			ast.Equal(int32(99), v["total"])
			continue
		}
		ast.Error(errors.New("error"), "impossible")
	}
	// Iter()
	iter := cli.Aggregate(context.Background(), Pipeline{matchStage, groupStage})
	ast.NotNil(iter)
	err = iter.All(&showsWithInfo)
	ast.NoError(err)
	for _, v := range showsWithInfo {
		if "Alice" == v["_id"] {
			ast.Equal(int32(12), v["total"])
			continue
		}
		if "Lucas" == v["_id"] {
			ast.Equal(int32(99), v["total"])
			continue
		}
		ast.Error(errors.New("error"), "impossible")
	}
	// One()
	var oneInfo bson.M

	iter = cli.Aggregate(context.Background(), Pipeline{matchStage, groupStage})
	ast.NotNil(iter)

	err = iter.One(&oneInfo)
	ast.NoError(err)
	ast.Equal(true, oneInfo["_id"] == "Alice" || oneInfo["_id"] == "Lucas")

	// iter
	iter = cli.Aggregate(context.Background(), Pipeline{matchStage, groupStage})
	ast.NotNil(iter)

	i := iter.Iter()

	ct := i.Next(&oneInfo)
	ast.Equal(true, oneInfo["_id"] == "Alice" || oneInfo["_id"] == "Lucas")
	ast.Equal(true, ct)
	ct = i.Next(&oneInfo)
	ast.Equal(true, oneInfo["_id"] == "Alice" || oneInfo["_id"] == "Lucas")
	ast.Equal(true, ct)
	ct = i.Next(&oneInfo)
	ast.Equal(false, ct)

	// err
	ast.Error(cli.Aggregate(context.Background(), 1).All(&showsWithInfo))
	ast.Error(cli.Aggregate(context.Background(), 1).One(&showsWithInfo))
	ast.Error(cli.Aggregate(context.Background(), 1).Iter().Err())
	matchStage = bson.D{{"$match", []bson.E{{"age", bson.D{{"$gt", 100}}}}}}
	groupStage = bson.D{{"$group", bson.D{{"_id", "$name"}, {"total", bson.D{{"$sum", "$age"}}}}}}
	ast.Error(cli.Aggregate(context.Background(), Pipeline{matchStage, groupStage}).One(&showsWithInfo))

}
