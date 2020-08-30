package qmgo

import (
	"context"
	"errors"
	"github.com/qiniu/qmgo/operator"
	"github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"

	"github.com/qiniu/qmgo/field"
	"github.com/stretchr/testify/require"
)

type UserHook struct {
	field.DefaultField `bson:",inline"`

	Name         string    `bson:"name"`
	Age          int       `bson:"age"`
	CreateTimeAt time.Time `bson:"createTimeAt"`
	UpdateTimeAt int64     `bson:"updateTimeAt"`
}

func (u *UserHook) BeforeInsert() error {
	if u.Name == "jz" || u.Name == "xm" {
		u.Age = 17
	}
	return nil
}

var afterInsertCount = 0

func (u *UserHook) AfterInsert() error {
	afterInsertCount++
	return nil
}

func (u *UserHook) CustomFields() field.CustomFieldsBuilder {
	return field.NewCustom().SetCreateAt("CreateTimeAt").SetUpdateAt("UpdateTimeAt")
}

type MyQueryHook struct {
	beforeCount int
	afterCount  int
}

func (q *MyQueryHook) BeforeQuery() error {
	q.beforeCount++
	return nil
}

func (q *MyQueryHook) AfterQuery() error {
	q.afterCount++
	return nil
}

// insert & insertMany可以直接用文档做为hook，为了保持和其他接口一致，还是用opts来传递hook，只是hook可以直接赋值为文档
// - 当然，也可以自己实现hook并通过opts注册
// - 能做自定义和default的field修改
func TestInsertHook(t *testing.T) {
	ast := require.New(t)
	cli := initClient("test")
	ctx := context.Background()
	defer cli.Close(ctx)
	defer cli.DropCollection(ctx)

	afterInsertCount = 0
	u := &UserHook{Name: "jz", Age: 7}
	_, err := cli.InsertOne(context.Background(), u, options.InsertOneOptions{
		InsertHook: u,
	})
	ast.NoError(err)

	uc := bson.M{"name": "jz"}
	ur := &UserHook{}
	uk := &MyQueryHook{}
	err = cli.Find(ctx, uc, options.FindOptions{
		QueryHook: uk,
	}).One(ur)
	ast.NoError(err)

	ast.Equal(17, ur.Age)
	// default fields
	ast.NotEqual(time.Time{}, ur.CreateAt)
	ast.NotEqual(time.Time{}, ur.UpdateAt)
	// custom fields
	ast.NotEqual(time.Time{}, ur.CreateTimeAt)
	ast.NotEqual(time.Time{}, ur.UpdateTimeAt)

	ast.Equal(1, afterInsertCount)
	ast.Equal(1, uk.beforeCount)
	ast.Equal(1, uk.afterCount)
}

// query因为不传入用户定义的文档结构体，所以需要
// - 自己实现hook并通过opts注册
// - 暂时不能在hook里修改文档里的东西
func TestInsertManyHook(t *testing.T) {
	ast := require.New(t)
	cli := initClient("test")
	ctx := context.Background()
	defer cli.Close(ctx)
	defer cli.DropCollection(ctx)

	afterInsertCount = 0
	u1 := &UserHook{Name: "jz", Age: 7}
	u2 := &UserHook{Name: "xm", Age: 7}
	us := []interface{}{u1, u2}
	_, err := cli.InsertMany(ctx, us, options.InsertManyOptions{
		InsertHook: us,
	})
	ast.NoError(err)

	uc := bson.M{"name": "jz"}
	ur := []UserHook{}
	qh := &MyQueryHook{}
	err = cli.Find(ctx, uc, options.FindOptions{
		QueryHook: qh,
	}).All(&ur)
	ast.NoError(err)

	ast.Equal(17, ur[0].Age)
	// default fields
	ast.NotEqual(time.Time{}, ur[0].CreateAt)
	ast.NotEqual(time.Time{}, ur[0].UpdateAt)
	// custom fields
	ast.NotEqual(time.Time{}, ur[0].CreateTimeAt)
	ast.NotEqual(time.Time{}, ur[0].UpdateTimeAt)

	ast.Equal(2, afterInsertCount)
	ast.Equal(1, qh.afterCount)
	ast.Equal(1, qh.beforeCount)

}

type MyUpdateHook struct {
	beforeUpdateCount int
	afterUpdateCount  int
}

func (u *MyUpdateHook) BeforeUpdate() error {
	u.beforeUpdateCount++
	return nil
}

func (u *MyUpdateHook) AfterUpdate() error {
	u.afterUpdateCount++
	return nil
}

// update因为不传入用户定义的文档结构体，所以需要
// - 自己实现hook并通过opts注册
// - 暂时不能在hook里修改文档里的东西
func TestUpdateHook(t *testing.T) {
	ast := require.New(t)
	cli := initClient("test")
	ctx := context.Background()
	defer cli.Close(ctx)
	defer cli.DropCollection(ctx)

	u := UserHook{Name: "jz", Age: 7}
	uh := &MyUpdateHook{}
	_, err := cli.InsertOne(context.Background(), u)
	ast.NoError(err)

	err = cli.UpdateOne(ctx, bson.M{"name": "jz"}, bson.M{operator.Set: bson.M{"age": 27}}, options.UpdateOptions{
		UpdateHook: uh,
	})
	ast.NoError(err)
	ast.Equal(1, uh.beforeUpdateCount)
	ast.Equal(1, uh.afterUpdateCount)

	cli.UpdateAll(ctx, bson.M{"name": "jz"}, bson.M{operator.Set: bson.M{"age": 27}}, options.UpdateOptions{
		UpdateHook: uh,
	})
	ast.NoError(err)
	ast.Equal(2, uh.beforeUpdateCount)
	ast.Equal(2, uh.afterUpdateCount)
}

type MyRemoveHook struct {
	beforeCount int
	afterCount  int
}

func (m *MyRemoveHook) BeforeRemove() error {
	m.beforeCount++
	return nil
}

func (m *MyRemoveHook) AfterRemove() error {
	m.afterCount++
	return nil
}

// remove因为不传入用户定义的文档结构体，所以需要
// - 自己实现hook并通过opts注册
// - 暂时不能在hook里修改文档里的东西
func TestRemoveHook(t *testing.T) {
	ast := require.New(t)
	cli := initClient("test")
	ctx := context.Background()
	defer cli.Close(ctx)
	defer cli.DropCollection(ctx)

	u := []interface{}{UserHook{Name: "jz", Age: 7}, UserHook{Name: "xm", Age: 7},
		UserHook{Name: "wxy", Age: 7}, UserHook{Name: "zp", Age: 7}}
	rlt, err := cli.InsertMany(context.Background(), u)
	ast.NoError(err)

	rh := &MyRemoveHook{}
	err = cli.RemoveId(ctx, rlt.InsertedIDs[0].(primitive.ObjectID).String(), options.RemoveOptions{
		RemoveHook: rh,
	})
	ast.NoError(err)
	ast.Equal(1, rh.afterCount)
	ast.Equal(1, rh.beforeCount)

	rh = &MyRemoveHook{}
	err = cli.Remove(ctx, bson.M{"age": 17}, options.RemoveOptions{
		RemoveHook: rh,
	})
	ast.NoError(err)
	ast.Equal(1, rh.afterCount)
	ast.Equal(1, rh.beforeCount)

	rh = &MyRemoveHook{}
	_, err = cli.RemoveAll(ctx, bson.M{"age": "7"}, options.RemoveOptions{
		RemoveHook: rh,
	})
	ast.NoError(err)
	ast.Equal(1, rh.afterCount)
	ast.Equal(1, rh.beforeCount)

}

type MyErrorHook struct {
	beforeQCount int
	afterQCount  int
	beforeRCount int
	afterRCount  int
	beforeUCount int
	afterUCount  int
	beforeICount int
	afterICount  int
}

func (m *MyErrorHook) BeforeRemove() error {
	if m.beforeRCount == 0 {
		m.beforeRCount++
		return errors.New("error")
	}
	m.beforeRCount++
	return nil
}

func (m *MyErrorHook) AfterRemove() error {
	m.afterRCount++
	return errors.New("error")
}

func (m *MyErrorHook) BeforeQuery() error {
	if m.beforeQCount == 0 {
		m.beforeQCount++
		return errors.New("error")
	}
	m.beforeQCount++

	return nil
}

func (m *MyErrorHook) AfterQuery() error {
	m.afterQCount++
	return errors.New("error")
}

func (m *MyErrorHook) BeforeInsert() error {
	if m.beforeICount == 0 {
		m.beforeICount++
		return errors.New("error")
	}
	m.beforeICount++

	return nil
}

func (m *MyErrorHook) AfterInsert() error {
	m.afterICount++
	return errors.New("error")
}

func (m *MyErrorHook) BeforeUpdate() error {
	if m.beforeUCount == 0 {
		m.beforeUCount++
		return errors.New("error")
	}
	m.beforeUCount++
	return nil
}

func (m *MyErrorHook) AfterUpdate() error {
	m.afterUCount++
	return errors.New("error")
}

func TestHookErr(t *testing.T) {
	ast := require.New(t)
	cli := initClient("test")
	ctx := context.Background()
	defer cli.Close(ctx)
	defer cli.DropCollection(ctx)

	u := &UserHook{Name: "jz", Age: 7}
	myHook := &MyErrorHook{}
	_, err := cli.InsertOne(context.Background(), u, options.InsertOneOptions{
		InsertHook: myHook,
	})
	ast.Error(err)
	ast.Equal(1, myHook.beforeICount)
	ast.Equal(0, myHook.afterICount)

	_, err = cli.InsertOne(context.Background(), u, options.InsertOneOptions{
		InsertHook: myHook,
	})
	ast.Error(err)
	ast.Equal(2, myHook.beforeICount)
	ast.Equal(1, myHook.afterICount)

	err = cli.UpdateOne(ctx, bson.M{"name": "jz"}, bson.M{operator.Set: bson.M{"age": 27}}, options.UpdateOptions{
		UpdateHook: myHook,
	})
	ast.Error(err)
	ast.Equal(1, myHook.beforeUCount)
	ast.Equal(0, myHook.afterUCount)

	err = cli.UpdateOne(ctx, bson.M{"name": "jz"}, bson.M{operator.Set: bson.M{"age": 27}}, options.UpdateOptions{
		UpdateHook: myHook,
	})
	ast.Error(err)
	ast.Equal(2, myHook.beforeUCount)
	ast.Equal(1, myHook.afterUCount)

	err = cli.Find(ctx, bson.M{"age": 27}, options.FindOptions{
		QueryHook: myHook,
	}).One(u)
	ast.Error(err)
	ast.Equal(1, myHook.beforeQCount)
	ast.Equal(0, myHook.afterQCount)

	err = cli.Find(ctx, bson.M{"age": 27}, options.FindOptions{
		QueryHook: myHook,
	}).One(u)
	ast.Error(err)
	ast.Equal(2, myHook.beforeQCount)
	ast.Equal(1, myHook.afterQCount)

	err = cli.Remove(ctx, bson.M{"age": 27}, options.RemoveOptions{
		RemoveHook: myHook,
	})
	ast.Error(err)
	ast.Equal(1, myHook.beforeRCount)
	ast.Equal(0, myHook.afterRCount)

	err = cli.Remove(ctx, bson.M{"age": 27}, options.RemoveOptions{
		RemoveHook: myHook,
	})
	ast.Error(err)
	ast.Equal(2, myHook.beforeRCount)
	ast.Equal(1, myHook.afterRCount)

}
