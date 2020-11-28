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

package field

import (
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

type CustomUser struct {
	Create        time.Time
	Update        int64
	MyId          primitive.ObjectID
	MyIdString    string
	InvalidId     int
	InvalidCreate int
	InvalidUpdate float32
}

func (c *CustomUser) CustomFields() CustomFieldsBuilder {
	return NewCustom().SetUpdateAt("Create").SetCreateAt("Update").SetId("MyId")
}

func (c *CustomUser) CustomFieldsIdString() CustomFieldsBuilder {
	return NewCustom().SetId("MyIdString")
}

func TestCustomFields(t *testing.T) {
	ast := require.New(t)
	u := &CustomUser{}
	c := u.CustomFields()
	c.(*CustomFields).CustomCreateTime(u)
	c.(*CustomFields).CustomUpdateTime(u)
	c.(*CustomFields).CustomId(u)
	ast.NotEqual(0, u.Update)
	ast.NotEqual(time.Time{}, u.Create)
	ast.NotEqual(primitive.NilObjectID, u.MyId)

	// id string
	u1 := &CustomUser{}
	c1 := u.CustomFieldsIdString()
	c1.(*CustomFields).CustomId(u1)
	ast.NotEqual("", u1.MyIdString)

}

func (c *CustomUser) CustomFieldsInvalid() CustomFieldsBuilder {
	return NewCustom().SetCreateAt("InvalidCreate")
}
func (c *CustomUser) CustomFieldsInvalid2() CustomFieldsBuilder {
	return NewCustom().SetUpdateAt("InvalidUpdate")
}

func (c *CustomUser) CustomFieldsInvalid3() CustomFieldsBuilder {
	return NewCustom().SetId("InvalidId")
}

func TestCustomFieldsInvalid(t *testing.T) {
	u := &CustomUser{}
	c := u.CustomFieldsInvalid()
	c.(*CustomFields).CustomCreateTime(u)
	c.(*CustomFields).CustomUpdateTime(u)
	ast := require.New(t)
	ast.Equal(0, u.InvalidCreate)
	ast.Equal(float32(0), u.InvalidUpdate)

	u1 := &CustomUser{}
	c = u1.CustomFieldsInvalid2()
	c.(*CustomFields).CustomCreateTime(u1)
	c.(*CustomFields).CustomUpdateTime(u1)
	ast.Equal(0, u1.InvalidCreate)
	ast.Equal(float32(0), u1.InvalidUpdate)

	u2 := CustomUser{}
	c = u2.CustomFieldsInvalid()
	c.(*CustomFields).CustomCreateTime(u2)
	c.(*CustomFields).CustomUpdateTime(u2)
	ast.Equal(0, u2.InvalidCreate)
	ast.Equal(float32(0), u2.InvalidUpdate)

	u3 := CustomUser{}
	c = u3.CustomFieldsInvalid3()
	c.(*CustomFields).CustomId(u3)
	ast.Equal(0, u3.InvalidId)

	u4 := &CustomUser{}
	c = u4.CustomFieldsInvalid3()
	c.(*CustomFields).CustomId(u4)
	ast.Equal(0, u4.InvalidId)

}
