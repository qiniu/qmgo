package field

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	DefaultField `bson:",inline"`

	Name         string             `bson:"name"`
	Age          int                `bson:"age"`
	CreateTimeAt time.Time          `bson:"createTimeAt"`
	UpdateTimeAt int64              `bson:"updateTimeAt"`
	MyId         primitive.ObjectID `bson:"myId"`
}

func (u *User) CustomFields() CustomFieldsBuilder {
	return NewCustom().SetCreateAt("CreateTimeAt").SetUpdateAt("UpdateTimeAt").SetId("MyId")
}

func TestBeforeInsert(t *testing.T) {
	ast := require.New(t)

	u := &User{Name: "Lucas", Age: 7}
	err := Do(u, BeforeInsert)
	ast.NoError(err)
	// default fields
	ast.NotEqual(time.Time{}, u.CreateAt)
	ast.NotEqual(time.Time{}, u.UpdateAt)
	ast.NotEqual(primitive.NilObjectID, u.Id)
	// custom fields
	ast.NotEqual(time.Time{}, u.CreateTimeAt)
	ast.NotEqual(time.Time{}, u.UpdateTimeAt)
	ast.NotEqual(primitive.NilObjectID, u.MyId)

	u1, u2 := &User{Name: "Lucas", Age: 7}, &User{Name: "Alice", Age: 8}
	us := []*User{u1, u2}
	err = Do(us, BeforeInsert)
	ast.NoError(err)

	for _, v := range us {
		ast.NotEqual(time.Time{}, v.CreateAt)
		ast.NotEqual(time.Time{}, v.UpdateAt)
		ast.NotEqual(primitive.NilObjectID, v.Id)
	}

	u3 := User{Name: "Lucas", Age: 7}
	err = Do(u3, BeforeInsert)
	ast.NoError(err)

	// insert with valid value
	tBefore3s := time.Now().Add(-3 * time.Second)
	id := primitive.NewObjectID()
	u = &User{Name: "Lucas", Age: 7}
	u.CreateAt = tBefore3s
	u.UpdateAt = tBefore3s
	u.Id = id
	u.MyId = id
	u.CreateTimeAt = tBefore3s
	u.UpdateTimeAt = tBefore3s.Unix()

	err = Do(u, BeforeUpsert)
	ast.NoError(err)

	ast.Equal(tBefore3s, u.CreateAt)
	ast.Equal(id, u.Id)
	ast.NotEqual(tBefore3s, u.UpdateAt)

	ast.Equal(tBefore3s, u.CreateTimeAt)
	ast.Equal(id, u.MyId)
	ast.NotEqual(tBefore3s.Unix(), u.UpdateTimeAt)
}

func TestBeforeUpdate(t *testing.T) {
	ast := require.New(t)

	u := &User{Name: "Lucas", Age: 7}
	err := Do(u, BeforeUpdate)
	ast.NoError(err)
	// default field
	ast.NotEqual(time.Time{}, u.UpdateAt)

	// custom fields
	ast.NotEqual(time.Time{}, u.UpdateTimeAt)

	u1, u2 := &User{Name: "Lucas", Age: 7}, &User{Name: "Alice", Age: 8}
	us := []*User{u1, u2}
	err = Do(us, BeforeUpdate)
	ast.NoError(err)
	for _, v := range us {
		// default field
		ast.NotEqual(time.Time{}, v.UpdateAt)

		// custom fields
		ast.NotEqual(time.Time{}, v.UpdateTimeAt)
	}

	us1 := []interface{}{u1, u2}
	err = Do(us1, BeforeUpdate)
	ast.NoError(err)
	for _, v := range us {
		// default field
		ast.NotEqual(time.Time{}, v.UpdateAt)

		// custom fields
		ast.NotEqual(time.Time{}, v.UpdateTimeAt)
	}

}

func TestBeforeUpsert(t *testing.T) {
	ast := require.New(t)

	// with empty fileds
	u := &User{Name: "Lucas", Age: 7}
	err := Do(u, BeforeUpsert)
	ast.NoError(err)
	// default fields
	ast.NotEqual(time.Time{}, u.CreateAt)
	ast.NotEqual(time.Time{}, u.UpdateAt)
	ast.NotEqual(primitive.NilObjectID, u.Id)
	// custom fields
	ast.NotEqual(time.Time{}, u.CreateTimeAt)
	ast.NotEqual(time.Time{}, u.UpdateTimeAt)
	ast.NotEqual(primitive.NilObjectID, u.MyId)

	u1, u2 := &User{Name: "Lucas", Age: 7}, &User{Name: "Alice", Age: 8}
	us := []*User{u1, u2}
	err = Do(us, BeforeUpsert)
	ast.NoError(err)

	for _, v := range us {
		ast.NotEqual(time.Time{}, v.CreateAt)
		ast.NotEqual(time.Time{}, v.UpdateAt)
		ast.NotEqual(primitive.NilObjectID, v.Id)
	}

	u3 := User{Name: "Lucas", Age: 7}
	err = Do(u3, BeforeUpsert)
	ast.NoError(err)

	// upsert with valid value
	tBefore3s := time.Now().Add(-3 * time.Second)
	id := primitive.NewObjectID()
	u = &User{Name: "Lucas", Age: 7}
	u.CreateAt = tBefore3s
	u.UpdateAt = tBefore3s
	u.Id = id
	u.MyId = id
	u.CreateTimeAt = tBefore3s
	u.UpdateTimeAt = tBefore3s.Unix()

	err = Do(u, BeforeUpsert)
	ast.NoError(err)

	ast.Equal(tBefore3s, u.CreateAt)
	ast.Equal(id, u.Id)
	ast.NotEqual(tBefore3s, u.UpdateAt)

	ast.Equal(tBefore3s, u.CreateTimeAt)
	ast.Equal(id, u.MyId)
	ast.NotEqual(tBefore3s.Unix(), u.UpdateTimeAt)

}

func TestNilError(t *testing.T) {
	ast := require.New(t)

	err := Do(nil, BeforeUpsert)
	ast.NoError(err)

}
