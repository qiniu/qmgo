package field

import (
	"reflect"
)

type FieldType string

const (
	BeforeInsert FieldType = "beforeInsert"
	BeforeUpdate FieldType = "beforeUpdate"
)

// filedHandler defines the relations between field type and handler
var fieldHandler = map[FieldType]func(doc interface{}) error{
	BeforeInsert: beforeInsert,
	BeforeUpdate: beforeUpdate,
}

// Do call the specific method to handle field based on fType
func Do(doc interface{}, fType FieldType) error {
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
func beforeInsert(doc interface{}) error {
	if ih, ok := doc.(DefaultFieldHook); ok {
		ih.DefaultCreateAt()
		ih.DefaultUpdateAt()
		ih.DefaultId()
	}
	if ih, ok := doc.(CustomFieldsHook); ok {
		fields := ih.CustomFields()
		fields.(*CustomFields).CustomCreateTime(doc)
		fields.(*CustomFields).CustomUpdateTime(doc)
		fields.(*CustomFields).CustomId(doc)
	}

	return nil
}

// beforeUpdate handles field before update
func beforeUpdate(hook interface{}) error {
	if ih, ok := hook.(DefaultFieldHook); ok {
		ih.DefaultUpdateAt()
	}
	if ih, ok := hook.(CustomFieldsHook); ok {
		fields := ih.CustomFields()
		fields.(*CustomFields).CustomUpdateTime(hook)
	}
	return nil
}
