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

package hook

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

type User struct {
	Name string `bson:"name"`
	Age  int    `bson:"age"`
	// for test
	AfterCalled bool `bson:"afterCalled"`
}

func (u *User) BeforeInsert() error {
	if u.Name == "Lucas" || u.Name == "Alices" {
		u.Age = 17
	}
	return nil
}

func (u *User) AfterInsert() error {
	u.AfterCalled = true
	return nil
}

func TestInsertHook(t *testing.T) {
	ast := require.New(t)

	u := &User{Name: "Lucas", Age: 7}
	err := Do(u, BeforeInsert)
	ast.NoError(err)
	ast.Equal(17, u.Age)

	err = Do(u, AfterInsert)
	ast.NoError(err)
	ast.True(u.AfterCalled)

	u1, u2 := &User{Name: "Lucas", Age: 7}, &User{Name: "Alices", Age: 8}
	us := []interface{}{u1, u2}
	err = Do(us, BeforeInsert)
	ast.NoError(err)
	err = Do(us, AfterInsert)
	ast.NoError(err)
	for _, v := range us {
		vv := v.(*User)
		if vv.Name == "Lucas" {
			ast.Equal(17, vv.Age)
		}
		if vv.Name == "Alices" {
			ast.Equal(17, vv.Age)
		}
		ast.True(vv.AfterCalled)
	}

	u3 := User{Name: "Lucas", Age: 7}
	err = Do(u3, BeforeInsert)
	ast.NoError(err)
}

func (u *User) BeforeUpdate() error {
	if u.Name == "Lucas" || u.Name == "Alices" {
		u.Age = 17
	}
	return nil
}

func (u *User) AfterUpdate() error {
	u.AfterCalled = true
	return nil
}
func TestUpdateHook(t *testing.T) {
	ast := require.New(t)

	u := &User{Name: "Lucas", Age: 7}
	err := Do(u, BeforeUpdate)
	ast.NoError(err)
	ast.Equal(17, u.Age)

	err = Do(u, AfterUpdate)
	ast.NoError(err)
	ast.True(u.AfterCalled)

	u1, u2 := &User{Name: "Lucas", Age: 7}, &User{Name: "Alices", Age: 8}
	us := []interface{}{u1, u2}
	err = Do(us, BeforeUpdate)
	ast.NoError(err)
	err = Do(us, AfterUpdate)
	ast.NoError(err)
	for _, v := range us {
		vv := v.(*User)
		if vv.Name == "Lucas" {
			ast.Equal(17, vv.Age)
		}
		if vv.Name == "Alices" {
			ast.Equal(17, vv.Age)
		}
		ast.True(vv.AfterCalled)
	}

}

func (u *User) BeforeQuery() error {
	if u.Name == "Lucas" || u.Name == "Alices" {
		u.Age = 17
	}
	fmt.Println("into before query")
	sliceBeforeQueryCount++
	return nil
}

func (u *User) AfterQuery() error {
	u.AfterCalled = true
	return nil
}
func TestQueryHook(t *testing.T) {
	ast := require.New(t)

	u := &User{Name: "Lucas", Age: 7}
	err := Do(u, BeforeQuery)
	ast.NoError(err)
	ast.Equal(17, u.Age)

	err = Do(u, AfterQuery)
	ast.NoError(err)
	ast.True(u.AfterCalled)

	u1, u2 := &User{Name: "Lucas", Age: 7}, &User{Name: "Alices", Age: 8}
	us := []interface{}{u1, u2}
	err = Do(us, BeforeQuery)
	ast.NoError(err)
	err = Do(us, AfterQuery)
	ast.NoError(err)
	for _, v := range us {
		vv := v.(*User)
		if vv.Name == "Lucas" {
			ast.Equal(17, vv.Age)
		}
		if vv.Name == "Alices" {
			ast.Equal(17, vv.Age)
		}
		ast.True(vv.AfterCalled)
	}

	uss := []*User{&User{Name: "Lucas"}, &User{Name: "Alices"}}
	Do(&uss, BeforeQuery)

}
func (u *User) BeforeRemove() error {
	if u.Name == "Lucas" || u.Name == "Alices" {
		u.Age = 17
	}
	return nil
}

func (u *User) AfterRemove() error {
	u.AfterCalled = true
	return nil
}

func TestRemoveHook(t *testing.T) {
	ast := require.New(t)

	u := &User{Name: "Lucas", Age: 7}
	err := Do(u, BeforeRemove)
	ast.NoError(err)
	ast.Equal(17, u.Age)

	err = Do(u, AfterRemove)
	ast.NoError(err)
	ast.True(u.AfterCalled)

	u1, u2 := &User{Name: "Lucas", Age: 7}, &User{Name: "Alices", Age: 8}
	us := []interface{}{u1, u2}
	err = Do(us, BeforeRemove)
	ast.NoError(err)
	err = Do(us, AfterRemove)
	ast.NoError(err)
	for _, v := range us {
		vv := v.(*User)
		if vv.Name == "Lucas" {
			ast.Equal(17, vv.Age)
		}
		if vv.Name == "Alices" {
			ast.Equal(17, vv.Age)
		}
		ast.True(vv.AfterCalled)
	}

}
func (u *User) BeforeUpsert() error {
	if u.Name == "Lucas" || u.Name == "Alices" {
		u.Age = 17
	}
	return nil
}

func (u *User) AfterUpsert() error {
	u.AfterCalled = true
	return nil
}

func TestUpsertHook(t *testing.T) {
	ast := require.New(t)

	u := &User{Name: "Lucas", Age: 7}
	err := Do(u, BeforeUpsert)
	ast.NoError(err)
	ast.Equal(17, u.Age)

	err = Do(u, AfterUpsert)
	ast.NoError(err)
	ast.True(u.AfterCalled)

	u1, u2 := &User{Name: "Lucas", Age: 7}, &User{Name: "Alices", Age: 8}
	us := []interface{}{u1, u2}
	err = Do(us, BeforeUpsert)
	ast.NoError(err)
	err = Do(us, AfterUpsert)
	ast.NoError(err)
	for _, v := range us {
		vv := v.(*User)
		if vv.Name == "Lucas" {
			ast.Equal(17, vv.Age)
		}
		if vv.Name == "Alices" {
			ast.Equal(17, vv.Age)
		}
		ast.True(vv.AfterCalled)
	}

	u3 := User{Name: "Lucas", Age: 7}
	err = Do(u3, BeforeInsert)
	ast.NoError(err)
}

type UserError struct {
	Name string
	Age  int

	// for test
	mock.Mock `bson:"-"`
}

func (u *UserError) BeforeInsert() error {
	return nil
}

func (u *UserError) AfterInsert() error {
	args := u.Called()
	return args.Error(0)
}
func TestSliceError(t *testing.T) {
	ast := require.New(t)

	u1, u2 := &UserError{Name: "Lucas", Age: 7}, &UserError{Name: "Alices", Age: 8}
	us := []interface{}{u1, u2}

	u1.On("AfterInsert").Return(nil)
	u2.On("AfterInsert").Return(errors.New("called"))
	err := Do(us, AfterInsert)
	ast.Equal("called", err.Error())

}

type UserNoHook struct {
	Name string
	Age  int

	// for test
	AfterCalled bool
}

func TestUserNoHook(t *testing.T) {
	ast := require.New(t)

	u := &UserNoHook{Name: "Lucas", Age: 7}
	err := Do(u, BeforeInsert)
	ast.NoError(err)
	ast.Equal(7, u.Age)

	err = Do(u, AfterInsert)
	ast.NoError(err)

	u1, u2 := &UserNoHook{Name: "Lucas", Age: 7}, &UserNoHook{Name: "Alices", Age: 8}
	us := []interface{}{u1, u2}
	err = Do(us, BeforeInsert)
	ast.NoError(err)
	err = Do(us, AfterInsert)
	ast.NoError(err)
	for _, v := range us {
		vv := v.(*UserNoHook)
		if vv.Name == "Lucas" {
			ast.Equal(7, vv.Age)
		}
		if vv.Name == "Alices" {
			ast.Equal(8, vv.Age)
		}
		ast.False(vv.AfterCalled)
	}

	err = Do(u, BeforeUpdate)
	ast.NoError(err)
	err = Do(u, AfterUpdate)
	ast.NoError(err)
	err = Do(us, BeforeUpdate)
	ast.NoError(err)
	err = Do(us, AfterUpdate)
	ast.NoError(err)

	err = Do(u, BeforeQuery)
	ast.NoError(err)
	err = Do(u, AfterQuery)
	ast.NoError(err)
	err = Do(us, BeforeQuery)
	ast.NoError(err)
	err = Do(us, AfterQuery)
	ast.NoError(err)

	err = Do(u, BeforeRemove)
	ast.NoError(err)
	err = Do(u, AfterRemove)
	ast.NoError(err)
	err = Do(us, BeforeRemove)
	ast.NoError(err)
	err = Do(us, AfterRemove)
	ast.NoError(err)
}

var sliceBeforeQueryCount = 0

func TestSliceHook(t *testing.T) {
	sliceBeforeQueryCount = 0
	ast := require.New(t)

	u := &User{Name: "Lucas"}
	Do(u, BeforeQuery)

	uss := []*User{&User{Name: "Lucas"}, &User{Name: "Alices"}}
	Do(uss, BeforeQuery)

	Do(&uss, BeforeQuery)

	ast.Equal(5, sliceBeforeQueryCount)

}

func TestNilError(t *testing.T) {
	ast := require.New(t)

	err := Do(nil, BeforeUpsert)
	ast.NoError(err)

}
