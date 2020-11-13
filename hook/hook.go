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
	"reflect"
)

type HookType string

const (
	BeforeInsert HookType = "beforeInsert"
	AfterInsert  HookType = "afterInsert"
	BeforeUpdate HookType = "beforeUpdate"
	AfterUpdate  HookType = "afterUpdate"
	BeforeQuery  HookType = "beforeQuery"
	AfterQuery   HookType = "afterQuery"
	BeforeRemove HookType = "beforeRemove"
	AfterRemove  HookType = "afterRemove"
	BeforeUpsert HookType = "beforeUpsert"
	AfterUpsert  HookType = "afterUpsert"
)

// hookHandler defines the relations between hook type and handler
var hookHandler = map[HookType]func(hook interface{}) error{
	BeforeInsert: beforeInsert,
	AfterInsert:  afterInsert,
	BeforeUpdate: beforeUpdate,
	AfterUpdate:  afterUpdate,
	BeforeQuery:  beforeQuery,
	AfterQuery:   afterQuery,
	BeforeRemove: beforeRemove,
	AfterRemove:  afterRemove,
	BeforeUpsert: beforeUpsert,
	AfterUpsert:  afterUpsert,
}

// Do call the specific method to handle hook based on hType
func Do(hook interface{}, hType HookType) error {
	to := reflect.TypeOf(hook)
	if to == nil {
		return nil
	}
	switch to.Kind() {
	case reflect.Slice:
		return sliceHandle(hook, hType)
	case reflect.Ptr:
		v := reflect.ValueOf(hook).Elem()
		switch v.Kind() {
		case reflect.Slice:
			return sliceHandle(v.Interface(), hType)
		default:
			return hookHandler[hType](hook)
		}
	default:
		return hookHandler[hType](hook)
	}
	return nil
}

// sliceHandle handles the slice hooks
func sliceHandle(hook interface{}, hType HookType) error {
	// []interface{}{UserType{}...}
	if h, ok := hook.([]interface{}); ok {
		for _, v := range h {
			if err := hookHandler[hType](v); err != nil {
				return err
			}
		}
		return nil
	}
	// []UserType{}
	s := reflect.ValueOf(hook)
	for i := 0; i < s.Len(); i++ {
		if err := hookHandler[hType](s.Index(i).Interface()); err != nil {
			return err
		}
	}
	return nil
}

// InsertHook defines the insert hook interface
type BeforeInsertHook interface {
	BeforeInsert() error
}
type AfterInsertHook interface {
	AfterInsert() error
}

// beforeInsert calls custom BeforeInsert
func beforeInsert(hook interface{}) error {
	if ih, ok := hook.(BeforeInsertHook); ok {
		return ih.BeforeInsert()
	}
	return nil
}

// afterInsert calls custom AfterInsert
func afterInsert(hook interface{}) error {
	if ih, ok := hook.(AfterInsertHook); ok {
		return ih.AfterInsert()
	}
	return nil
}

// UpdateHook defines the Update hook interface
type BeforeUpdateHook interface {
	BeforeUpdate() error
}
type AfterUpdateHook interface {
	AfterUpdate() error
}

// beforeUpdate calls custom BeforeUpdate
func beforeUpdate(hook interface{}) error {
	if ih, ok := hook.(BeforeUpdateHook); ok {
		return ih.BeforeUpdate()
	}
	return nil
}

// afterUpdate calls custom AfterUpdate
func afterUpdate(hook interface{}) error {
	if ih, ok := hook.(AfterUpdateHook); ok {
		return ih.AfterUpdate()
	}
	return nil
}

// QueryHook defines the query hook interface
type BeforeQueryHook interface {
	BeforeQuery() error
}
type AfterQueryHook interface {
	AfterQuery() error
}

// beforeQuery calls custom BeforeQuery
func beforeQuery(hook interface{}) error {
	if ih, ok := hook.(BeforeQueryHook); ok {
		return ih.BeforeQuery()
	}
	return nil
}

// afterQuery calls custom AfterQuery
func afterQuery(doc interface{}) error {
	if ih, ok := doc.(AfterQueryHook); ok {
		return ih.AfterQuery()
	}
	return nil
}

// RemoveHook defines the remove hook interface
type BeforeRemoveHook interface {
	BeforeRemove() error
}
type AfterRemoveHook interface {
	AfterRemove() error
}

// beforeRemove calls custom BeforeRemove
func beforeRemove(hook interface{}) error {
	if ih, ok := hook.(BeforeRemoveHook); ok {
		return ih.BeforeRemove()
	}
	return nil
}

// afterRemove calls custom AfterRemove
func afterRemove(hook interface{}) error {
	if ih, ok := hook.(AfterRemoveHook); ok {
		return ih.AfterRemove()
	}
	return nil
}

// UpsertHook defines the upsert hook interface
type BeforeUpsertHook interface {
	BeforeUpsert() error
}
type AfterUpsertHook interface {
	AfterUpsert() error
}

// beforeUpsert calls custom BeforeUpsert
func beforeUpsert(hook interface{}) error {
	if ih, ok := hook.(BeforeUpsertHook); ok {
		return ih.BeforeUpsert()
	}
	return nil
}

// afterUpsert calls custom AfterUpsert
func afterUpsert(hook interface{}) error {
	if ih, ok := hook.(AfterUpsertHook); ok {
		return ih.AfterUpsert()
	}
	return nil
}
