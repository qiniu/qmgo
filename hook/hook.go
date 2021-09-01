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
var hookHandler = map[operator.OpType]func(ctx context.Context, hook interface{}) error{
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
func Do(ctx context.Context, hook interface{}, opType operator.OpType, opts ...interface{}) error {
	if len(opts) > 0 {
		hook = opts[0]
	}

	to := reflect.TypeOf(hook)
	if to == nil {
		return nil
	}
	switch to.Kind() {
	case reflect.Slice:
		return sliceHandle(ctx, hook, opType)
	case reflect.Ptr:
		v := reflect.ValueOf(hook).Elem()
		switch v.Kind() {
		case reflect.Slice:
			return sliceHandle(ctx, v.Interface(), opType)
		default:
			return do(ctx, hook, opType)
		}
	default:
		return do(ctx, hook, opType)
	}
}

// sliceHandle handles the slice hooks
func sliceHandle(ctx context.Context, hook interface{}, opType operator.OpType) error {
	// []interface{}{UserType{}...}
	if h, ok := hook.([]interface{}); ok {
		for _, v := range h {
			if err := do(ctx, v, opType); err != nil {
				return err
			}
		}
		return nil
	}
	// []UserType{}
	s := reflect.ValueOf(hook)
	for i := 0; i < s.Len(); i++ {
		if err := do(ctx, s.Index(i).Interface(), opType); err != nil {
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
func beforeInsert(ctx context.Context, hook interface{}) error {
	if ih, ok := hook.(BeforeInsertHook); ok {
		return ih.BeforeInsert(ctx)
	}
	return nil
}

// afterInsert calls custom AfterInsert
func afterInsert(ctx context.Context, hook interface{}) error {
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
func beforeUpdate(ctx context.Context, hook interface{}) error {
	if ih, ok := hook.(BeforeUpdateHook); ok {
		return ih.BeforeUpdate(ctx)
	}
	return nil
}

// afterUpdate calls custom AfterUpdate
func afterUpdate(ctx context.Context, hook interface{}) error {
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
func beforeQuery(ctx context.Context, hook interface{}) error {
	if ih, ok := hook.(BeforeQueryHook); ok {
		return ih.BeforeQuery(ctx)
	}
	return nil
}

// afterQuery calls custom AfterQuery
func afterQuery(ctx context.Context, hook interface{}) error {
	if ih, ok := hook.(AfterQueryHook); ok {
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
func beforeRemove(ctx context.Context, hook interface{}) error {
	if ih, ok := hook.(BeforeRemoveHook); ok {
		return ih.BeforeRemove(ctx)
	}
	return nil
}

// afterRemove calls custom AfterRemove
func afterRemove(ctx context.Context, hook interface{}) error {
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
func beforeUpsert(ctx context.Context, hook interface{}) error {
	if ih, ok := hook.(BeforeUpsertHook); ok {
		return ih.BeforeUpsert(ctx)
	}
	return nil
}

// afterUpsert calls custom AfterUpsert
func afterUpsert(ctx context.Context, hook interface{}) error {
	if ih, ok := hook.(AfterUpsertHook); ok {
		return ih.AfterUpsert(ctx)
	}
	return nil
}

// do check if opType is supported and call hookHandler
func do(ctx context.Context, hook interface{}, opType operator.OpType) error {
	if f, ok := hookHandler[opType]; !ok {
		return nil
	} else {
		return f(ctx, hook)
	}
}
