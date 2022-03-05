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
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"time"
)

// CustomFields defines struct of supported custom fields
type CustomFields struct {
	createAt string
	updateAt string
	id       string
}

// CustomFieldsHook defines the interface, CustomFields return custom field user want to change
type CustomFieldsHook interface {
	CustomFields() CustomFieldsBuilder
}

// CustomFieldsBuilder defines the interface which user use to set custom fields
type CustomFieldsBuilder interface {
	SetUpdateAt(fieldName string) CustomFieldsBuilder
	SetCreateAt(fieldName string) CustomFieldsBuilder
	SetId(fieldName string) CustomFieldsBuilder
}

// NewCustom creates new Builder which is used to set the custom fields
func NewCustom() CustomFieldsBuilder {
	return &CustomFields{}
}

// SetUpdateAt set the custom UpdateAt field
func (c *CustomFields) SetUpdateAt(fieldName string) CustomFieldsBuilder {
	c.updateAt = fieldName
	return c
}

// SetCreateAt set the custom CreateAt field
func (c *CustomFields) SetCreateAt(fieldName string) CustomFieldsBuilder {
	c.createAt = fieldName
	return c
}

// SetId set the custom Id field
func (c *CustomFields) SetId(fieldName string) CustomFieldsBuilder {
	c.id = fieldName
	return c
}

// CustomCreateTime changes the custom create time
func (c CustomFields) CustomCreateTime(doc interface{}) {
	if c.createAt == "" {
		return
	}
	fieldName := c.createAt
	setTime(doc, fieldName, false)
	return
}

// CustomUpdateTime changes the custom update time
func (c CustomFields) CustomUpdateTime(doc interface{}) {
	if c.updateAt == "" {
		return
	}
	fieldName := c.updateAt
	setTime(doc, fieldName, true)
	return
}

// CustomUpdateTime changes the custom update time
func (c CustomFields) CustomId(doc interface{}) {
	if c.id == "" {
		return
	}
	fieldName := c.id
	setId(doc, fieldName)
	return
}

// setTime changes the custom time fields
// The overWrite defines if change value when the filed has valid value
func setTime(doc interface{}, fieldName string, overWrite bool) {
	if reflect.Ptr != reflect.TypeOf(doc).Kind() {
		fmt.Println("not a point type")
		return
	}
	e := reflect.ValueOf(doc).Elem()
	ca := e.FieldByName(fieldName)
	if ca.CanSet() {
		tt := time.Now()
		switch a := ca.Interface().(type) {
		case time.Time:
			if ca.Interface().(time.Time).IsZero() {
				ca.Set(reflect.ValueOf(tt))
			} else if overWrite {
				ca.Set(reflect.ValueOf(tt))
			}
		case int64:
			if ca.Interface().(int64) == 0 {
				ca.SetInt(tt.Unix())
			} else if overWrite {
				ca.SetInt(tt.Unix())
			}
		default:
			fmt.Println("unsupported type to setTime", a)
		}
	}
}

// setId changes the custom Id fields
func setId(doc interface{}, fieldName string) {
	if reflect.Ptr != reflect.TypeOf(doc).Kind() {
		fmt.Println("not a point type")
		return
	}
	e := reflect.ValueOf(doc).Elem()
	ca := e.FieldByName(fieldName)
	if ca.CanSet() {
		switch a := ca.Interface().(type) {
		case primitive.ObjectID:
			if ca.Interface().(primitive.ObjectID).IsZero() {
				ca.Set(reflect.ValueOf(primitive.NewObjectID()))
			}
		case string:
			if ca.String() == "" {
				ca.SetString(primitive.NewObjectID().Hex())
			}
		default:
			fmt.Println("unsupported type to setId", a)
		}
	}
}
