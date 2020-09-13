package qmgo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"

	"github.com/qiniu/qmgo/field"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

type UserField struct {
	field.DefaultField `bson:",inline"`

	Name string `bson:"name"`
	Age  int    `bson:"age"`

	MyId         string    `bson:"myId"`
	CreateTimeAt time.Time `bson:"createTimeAt"`
	UpdateTimeAt int64     `bson:"updateTimeAt"`
}

func (u *UserField) CustomFields() field.CustomFieldsBuilder {
	return field.NewCustom().SetCreateAt("CreateTimeAt").SetUpdateAt("UpdateTimeAt").SetId("MyId")
}

func TestInsertField(t *testing.T) {
	ast := require.New(t)
	cli := initClient("test")
	ctx := context.Background()
	defer cli.Close(ctx)
	defer cli.DropCollection(ctx)

	u := &UserField{Name: "Lucas", Age: 7}
	_, err := cli.InsertOne(context.Background(), u)
	ast.NoError(err)

	uc := bson.M{"name": "Lucas"}
	ur := &UserField{}
	err = cli.Find(ctx, uc).One(ur)
	ast.NoError(err)

	// default fields
	ast.NotEqual(time.Time{}, ur.CreateAt)
	ast.NotEqual(time.Time{}, ur.UpdateAt)
	ast.NotEqual(primitive.NilObjectID, ur.Id)
	// custom fields
	ast.NotEqual(time.Time{}, ur.CreateTimeAt)
	ast.NotEqual(int64(0), ur.UpdateTimeAt)
	ast.NotEqual("", ur.MyId)

}

func TestFieldInsertMany(t *testing.T) {
	ast := require.New(t)
	cli := initClient("test")
	ctx := context.Background()
	defer cli.Close(ctx)
	defer cli.DropCollection(ctx)

	u1 := &UserField{Name: "Lucas", Age: 7}
	u2 := &UserField{Name: "Alice", Age: 7}
	us := []*UserField{u1, u2}
	_, err := cli.InsertMany(ctx, us)
	ast.NoError(err)

	uc := bson.M{"age": 7}
	ur := []UserField{}
	err = cli.Find(ctx, uc).All(&ur)
	ast.NoError(err)

	// default fields
	ast.NotEqual(time.Time{}, ur[0].CreateAt)
	ast.NotEqual(time.Time{}, ur[0].UpdateAt)
	ast.NotEqual(primitive.NilObjectID, ur[0].Id)
	// default fields
	ast.NotEqual(time.Time{}, ur[1].CreateAt)
	ast.NotEqual(time.Time{}, ur[1].UpdateAt)
	ast.NotEqual(primitive.NilObjectID, ur[1].Id)

	// custom fields
	ast.NotEqual(time.Time{}, ur[0].CreateTimeAt)
	ast.NotEqual(int64(0), ur[0].UpdateTimeAt)
}

func TestUpdateDoc(t *testing.T) {
	ast := require.New(t)
	cli := initClient("test")
	defer cli.Close(context.Background())
	defer cli.DropCollection(context.Background())
	cli.EnsureIndexes(context.Background(), []string{"name"}, nil)

	ui := &UserField{Name: "Lucas", Age: 17}
	_, err := cli.InsertOne(context.Background(), ui)
	ast.NoError(err)

	err = cli.UpdateOne(context.Background(), bson.M{"name": "Lucas"}, bson.M{"$set": bson.M{"updateTimeAt": 0, "updateAt": time.Time{}}})
	ast.NoError(err)

	findUi := UserField{}
	err = cli.Find(context.Background(), bson.M{"name": "Lucas"}).One(&findUi)
	ast.Equal(int64(0), findUi.UpdateTimeAt)
	ast.Equal(time.Time{}, findUi.UpdateAt)

	ast.NoError(err)
	ui.Id = findUi.Id
	err = cli.UpdateWithDocument(context.Background(), bson.M{"_id": findUi.Id}, &ui)
	ast.NoError(err)
	err = cli.Find(context.Background(), bson.M{"name": "Lucas"}).One(&findUi)
	ast.NotEqual(int64(0), findUi.UpdateTimeAt)
	ast.NotEqual(time.Time{}, findUi.UpdateAt)

}
