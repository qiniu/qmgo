package hook

import (
	"github.com/qiniu/qmgo/field"
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
}

// Do call the specific method to handle hook based on hType
func Do(hook interface{}, hType HookType) error {
	switch reflect.TypeOf(hook).Kind() {
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
type InsertHook interface {
	BeforeInsert() error
	AfterInsert() error
}

// beforeInsert calls custom BeforeInsert
func beforeInsert(hook interface{}) error {
	if ih, ok := hook.(field.DefaultFieldHook); ok {
		ih.DefaultCreateAt()
		ih.DefaultUpdateAt()
	}
	if ih, ok := hook.(field.CustomFieldsHook); ok {
		fields := ih.CustomFields()
		fields.(*field.CustomFields).CustomCreateTime(hook)
		fields.(*field.CustomFields).CustomUpdateTime(hook)
	}
	if ih, ok := hook.(InsertHook); ok {
		return ih.BeforeInsert()
	}
	return nil
}

// afterInsert calls custom AfterInsert
func afterInsert(hook interface{}) error {
	if ih, ok := hook.(InsertHook); ok {
		return ih.AfterInsert()
	}
	return nil
}

// UpdateHook defines the Update hook interface
type UpdateHook interface {
	BeforeUpdate() error
	AfterUpdate() error
}

// beforeUpdate calls custom BeforeUpdate
func beforeUpdate(hook interface{}) error {
	if ih, ok := hook.(field.DefaultFieldHook); ok {
		ih.DefaultUpdateAt()
	}
	if ih, ok := hook.(field.CustomFieldsHook); ok {
		fields := ih.CustomFields()
		fields.(*field.CustomFields).CustomUpdateTime(hook)
	}
	if ih, ok := hook.(UpdateHook); ok {
		return ih.BeforeUpdate()
	}
	return nil
}

// afterUpdate calls custom AfterUpdate
func afterUpdate(hook interface{}) error {
	if ih, ok := hook.(UpdateHook); ok {
		return ih.AfterUpdate()
	}
	return nil
}

// QueryHook defines the query hook interface
type QueryHook interface {
	BeforeQuery() error
	AfterQuery() error
}

// beforeQuery calls custom BeforeQuery
func beforeQuery(hook interface{}) error {
	if ih, ok := hook.(QueryHook); ok {
		return ih.BeforeQuery()
	}
	return nil
}

// afterQuery calls custom AfterQuery
func afterQuery(doc interface{}) error {
	if ih, ok := doc.(QueryHook); ok {
		return ih.AfterQuery()
	}
	return nil
}

// RemoveHook defines the remove hook interface
type RemoveHook interface {
	BeforeRemove() error
	AfterRemove() error
}

// beforeRemove calls custom BeforeRemove
func beforeRemove(hook interface{}) error {
	if ih, ok := hook.(RemoveHook); ok {
		return ih.BeforeRemove()
	}
	return nil
}

// afterRemove calls custom AfterRemove
func afterRemove(hook interface{}) error {
	if ih, ok := hook.(RemoveHook); ok {
		return ih.AfterRemove()
	}
	return nil
}
