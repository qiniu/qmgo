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
	"context"
	"github.com/qiniu/qmgo/operator"
	"reflect"
)

// hookHandler defines the relations between hook type and handler
var hookHandler = map[operator.OpType]func(hook interface{}, ctx context.Context) error{
	operator.BeforeInsert:  beforeInsert,
	operator.AfterInsert:   afterInsert,
	operator.BeforeUpdate:  beforeUpdate,
	operator.AfterUpdate:   afterUpdate,
	operator.BeforeQuery:   beforeQuery,
	operator.AfterQuery:    afterQuery,
	operator.BeforeRemove:  beforeRemove,
	operator.AfterRemove:   afterRemove,
	operator.BeforeUpsert:  beforeUpsert,
	operator.AfterUpsert:   afterUpsert,
	operator.BeforeReplace: beforeUpdate,
	operator.AfterReplace:  afterUpdate,
}

//
//func init() {
//	middleware.Register(Do)
//}

// Do call the specific method to handle hook based on hType
// If opts has valid value, use it instead of original hook
func Do(hook interface{}, opType operator.OpType, ctx context.Context, opts ...interface{}) error {
	if len(opts) > 0 {
		hook = opts[0]
	}

	to := reflect.TypeOf(hook)
	if to == nil {
		return nil
	}
	switch to.Kind() {
	case reflect.Slice:
		return sliceHandle(hook, opType, ctx)
	case reflect.Ptr:
		v := reflect.ValueOf(hook).Elem()
		switch v.Kind() {
		case reflect.Slice:
			return sliceHandle(v.Interface(), opType, ctx)
		default:
			return do(hook, opType, ctx)
		}
	default:
		return do(hook, opType, ctx)
	}
}

// sliceHandle handles the slice hooks
func sliceHandle(hook interface{}, opType operator.OpType, ctx context.Context) error {
	// []interface{}{UserType{}...}
	if h, ok := hook.([]interface{}); ok {
		for _, v := range h {
			if err := do(v, opType, ctx); err != nil {
				return err
			}
		}
		return nil
	}
	// []UserType{}
	s := reflect.ValueOf(hook)
	for i := 0; i < s.Len(); i++ {
		if err := do(s.Index(i).Interface(), opType, ctx); err != nil {
			return err
		}
	}
	return nil
}

// BeforeInsertHook InsertHook defines the insert hook interface
type BeforeInsertHook interface {
	BeforeInsert(ctx context.Context) error
}
type AfterInsertHook interface {
	AfterInsert(ctx context.Context) error
}

// beforeInsert calls custom BeforeInsert
func beforeInsert(hook interface{}, ctx context.Context) error {
	if ih, ok := hook.(BeforeInsertHook); ok {
		return ih.BeforeInsert(ctx)
	}
	return nil
}

// afterInsert calls custom AfterInsert
func afterInsert(hook interface{}, ctx context.Context) error {
	if ih, ok := hook.(AfterInsertHook); ok {
		return ih.AfterInsert(ctx)
	}
	return nil
}

// BeforeUpdateHook defines the Update hook interface
type BeforeUpdateHook interface {
	BeforeUpdate(ctx context.Context) error
}
type AfterUpdateHook interface {
	AfterUpdate(ctx context.Context) error
}

// beforeUpdate calls custom BeforeUpdate
func beforeUpdate(hook interface{}, ctx context.Context) error {
	if ih, ok := hook.(BeforeUpdateHook); ok {
		return ih.BeforeUpdate(ctx)
	}
	return nil
}

// afterUpdate calls custom AfterUpdate
func afterUpdate(hook interface{}, ctx context.Context) error {
	if ih, ok := hook.(AfterUpdateHook); ok {
		return ih.AfterUpdate(ctx)
	}
	return nil
}

// BeforeQueryHook QueryHook defines the query hook interface
type BeforeQueryHook interface {
	BeforeQuery(ctx context.Context) error
}
type AfterQueryHook interface {
	AfterQuery(ctx context.Context) error
}

// beforeQuery calls custom BeforeQuery
func beforeQuery(hook interface{}, ctx context.Context) error {
	if ih, ok := hook.(BeforeQueryHook); ok {
		return ih.BeforeQuery(ctx)
	}
	return nil
}

// afterQuery calls custom AfterQuery
func afterQuery(doc interface{}, ctx context.Context) error {
	if ih, ok := doc.(AfterQueryHook); ok {
		return ih.AfterQuery(ctx)
	}
	return nil
}

// BeforeRemoveHook RemoveHook defines the remove hook interface
type BeforeRemoveHook interface {
	BeforeRemove(ctx context.Context) error
}
type AfterRemoveHook interface {
	AfterRemove(ctx context.Context) error
}

// beforeRemove calls custom BeforeRemove
func beforeRemove(hook interface{}, ctx context.Context) error {
	if ih, ok := hook.(BeforeRemoveHook); ok {
		return ih.BeforeRemove(ctx)
	}
	return nil
}

// afterRemove calls custom AfterRemove
func afterRemove(hook interface{}, ctx context.Context) error {
	if ih, ok := hook.(AfterRemoveHook); ok {
		return ih.AfterRemove(ctx)
	}
	return nil
}

// BeforeUpsertHook UpsertHook defines the upsert hook interface
type BeforeUpsertHook interface {
	BeforeUpsert(ctx context.Context) error
}
type AfterUpsertHook interface {
	AfterUpsert(ctx context.Context) error
}

// beforeUpsert calls custom BeforeUpsert
func beforeUpsert(hook interface{}, ctx context.Context) error {
	if ih, ok := hook.(BeforeUpsertHook); ok {
		return ih.BeforeUpsert(ctx)
	}
	return nil
}

// afterUpsert calls custom AfterUpsert
func afterUpsert(hook interface{}, ctx context.Context) error {
	if ih, ok := hook.(AfterUpsertHook); ok {
		return ih.AfterUpsert(ctx)
	}
	return nil
}

// do check if opType is supported and call hookHandler
func do(hook interface{}, opType operator.OpType, ctx context.Context) error {
	if f, ok := hookHandler[opType]; !ok {
		return nil
	} else {
		return f(hook, ctx)
	}
}
