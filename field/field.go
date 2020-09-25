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
	"reflect"
	"time"
)

var nilTime time.Time

type FieldType string

const (
	BeforeInsert FieldType = "beforeInsert"
	BeforeUpdate FieldType = "beforeUpdate"
	BeforeUpsert FieldType = "beforeUpsert"
)

// filedHandler defines the relations between field type and handler
var fieldHandler = map[FieldType]func(doc interface{}) error{
	BeforeInsert: beforeInsert,
	BeforeUpdate: beforeUpdate,
	BeforeUpsert: beforeUpsert,
}

// Do call the specific method to handle field based on fType
func Do(doc interface{}, fType FieldType) error {
	to := reflect.TypeOf(doc)
	if to == nil {
		return nil
	}
	switch reflect.TypeOf(doc).Kind() {
	case reflect.Slice:
		return sliceHandle(doc, fType)
	case reflect.Ptr:
		v := reflect.ValueOf(doc).Elem()
		switch v.Kind() {
		case reflect.Slice:
			return sliceHandle(v.Interface(), fType)
		default:
			return fieldHandler[fType](doc)
		}
	}
	//fmt.Println("not support type")
	return nil
}

// sliceHandle handles the slice docs
func sliceHandle(docs interface{}, fType FieldType) error {
	// []interface{}{UserType{}...}
	if h, ok := docs.([]interface{}); ok {
		for _, v := range h {
			if err := fieldHandler[fType](v); err != nil {
				return err
			}
		}
		return nil
	}
	// []UserType{}
	s := reflect.ValueOf(docs)
	for i := 0; i < s.Len(); i++ {
		if err := fieldHandler[fType](s.Index(i).Interface()); err != nil {
			return err
		}
	}
	return nil
}

// beforeInsert handles field before insert
// If value of field createAt is valid in doc, upsert doesn't change it
// If value of field id is valid in doc, upsert doesn't change it
// Change the value of field updateAt anyway
func beforeInsert(doc interface{}) error {
	if ih, ok := doc.(DefaultFieldHook); ok {
		ih.DefaultId()
		ih.DefaultCreateAt()
		ih.DefaultUpdateAt()
	}
	if ih, ok := doc.(CustomFieldsHook); ok {
		fields := ih.CustomFields()
		fields.(*CustomFields).CustomId(doc)
		fields.(*CustomFields).CustomCreateTime(doc)
		fields.(*CustomFields).CustomUpdateTime(doc)
	}
	return nil
}

// beforeUpdate handles field before update
func beforeUpdate(doc interface{}) error {
	if ih, ok := doc.(DefaultFieldHook); ok {
		ih.DefaultUpdateAt()
	}
	if ih, ok := doc.(CustomFieldsHook); ok {
		fields := ih.CustomFields()
		fields.(*CustomFields).CustomUpdateTime(doc)
	}
	return nil
}

// beforeUpsert handles field before upsert
// If value of field createAt is valid in doc, upsert doesn't change it
// If value of field id is valid in doc, upsert doesn't change it
// Change the value of field updateAt anyway
func beforeUpsert(doc interface{}) error {
	if ih, ok := doc.(DefaultFieldHook); ok {
		ih.DefaultId()
		ih.DefaultCreateAt()
		ih.DefaultUpdateAt()
	}
	if ih, ok := doc.(CustomFieldsHook); ok {
		fields := ih.CustomFields()
		fields.(*CustomFields).CustomId(doc)
		fields.(*CustomFields).CustomCreateTime(doc)
		fields.(*CustomFields).CustomUpdateTime(doc)
	}
	return nil
}
