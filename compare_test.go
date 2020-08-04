package qmgo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"labix.org/v2/mgo"
)

const (
	MGO_URI  = "localhost:27017"
	URI      = "mongodb://localhost:27017"
	DATABASE = "class"
	COLL     = "user"
)

type BsonT map[string]interface{}

type UserInfo struct {
	Name   string `bson:"name"`
	Age    uint16 `bson:"age"`
	Weight uint32 `bson:"weight"`
}

var oneUserInfo = UserInfo{
	Name:   "xm",
	Age:    7,
	Weight: 40,
}

var batchUserInfo = []UserInfo{
	{Name: "wxy", Age: 6, Weight: 20},
	{Name: "jZ", Age: 6, Weight: 25},
	{Name: "zp", Age: 6, Weight: 30},
	{Name: "yxw", Age: 6, Weight: 35},
}
var batchUserInfoI = []interface{}{
	UserInfo{Name: "wxy", Age: 6, Weight: 20},
	UserInfo{Name: "jZ", Age: 6, Weight: 25},
	UserInfo{Name: "zp", Age: 6, Weight: 30},
	UserInfo{Name: "yxw", Age: 6, Weight: 35},
}

func TestMgo(t *testing.T) {
	ast := require.New(t)
	// create connection
	session, err := mgo.Dial(MGO_URI)
	ast.Nil(err)
	db := session.DB(DATABASE)
	defer db.DropDatabase()
	coll := db.C(COLL)

	// insert one document
	err = coll.Insert(oneUserInfo)
	ast.Nil(err)

	// find one document
	one := UserInfo{}
	coll.Find(BsonT{"name": oneUserInfo.Name}).One(&one)
	ast.Nil(err)
	ast.Equal(oneUserInfo, one)

	// batch insert
	for _, v := range batchUserInfo {
		err = coll.Insert(v)
		ast.Nil(err)
	}
	batch := []UserInfo{}

	// find all 、sort and limit
	coll.Find(BsonT{"age": 6}).Sort("weight").Limit(7).All(&batch)
	ast.Equal(true, reflect.DeepEqual(batchUserInfo, batch))
}

func TestQmgo(t *testing.T) {
	ast := require.New(t)
	ctx := context.Background()

	// create connect
	cli, err := Open(ctx, &Config{Uri: URI, Database: DATABASE, Coll: COLL})

	ast.Nil(err)
	defer func() {
		if err = cli.Close(ctx); err != nil {
			panic(err)
		}
	}()
	defer cli.DropDatabase(ctx)

	cli.EnsureIndexes(ctx, []string{"name"}, []string{"age", "name,weight"})
	// insert one document
	_, err = cli.InsertOne(ctx, oneUserInfo)
	ast.Nil(err)

	// find one document
	one := UserInfo{}
	err = cli.Find(ctx, BsonT{"name": oneUserInfo.Name}).One(&one)
	ast.Nil(err)
	ast.Equal(oneUserInfo, one)

	// batch insert
	_, err = cli.Collection.InsertMany(ctx, batchUserInfoI)
	ast.Nil(err)

	// find all 、sort and limit
	batch := []UserInfo{}
	cli.Find(ctx, BsonT{"age": 6}).Sort("weight").Limit(7).All(&batch)
	ast.Equal(true, reflect.DeepEqual(batchUserInfo, batch))

	//remove
	err = cli.Remove(ctx, BsonT{"age": 7})
	ast.Nil(err)
}

func TestOfficialMongoDriver(t *testing.T) {
	ast := require.New(t)
	ctx := context.Background()

	// create connect
	var opts *options.ClientOptions
	opts = new(options.ClientOptions)
	opts.ApplyURI(URI)
	c, err := mongo.Connect(ctx, opts)
	ast.Nil(err)
	db := c.Database(DATABASE)
	coll := db.Collection(COLL)
	defer db.Drop(ctx)

	// insert one document
	_, err = coll.InsertOne(ctx, oneUserInfo)
	ast.Nil(err)

	// find one document
	one := UserInfo{}
	err = coll.FindOne(ctx, BsonT{"name": oneUserInfo.Name}).Decode(&one)
	ast.Nil(err)

	// batch insert
	_, err = coll.InsertMany(ctx, batchUserInfoI)
	ast.Nil(err)

	// find all 、sort and limit
	findOptions := options.Find()
	findOptions.SetLimit(7)
	var sorts bson.D
	sorts = append(sorts, bson.E{Key: "weight", Value: 1})

	findOptions.SetSort(sorts)

	batch := []UserInfo{}
	cur, err := coll.Find(ctx, BsonT{"age": 6}, findOptions)
	ast.Nil(err)
	cur.All(ctx, &batch)
	ast.Equal(true, reflect.DeepEqual(batchUserInfo, batch))
}
