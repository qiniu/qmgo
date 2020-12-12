package validator

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/qiniu/qmgo/operator"
)

// use a single instance of Validate, it caches struct info
var validate = validator.New()

// validatorNeeded checks if the validator is needed to opType
func validatorNeeded(opType operator.OpType) bool {
	switch opType {
	case operator.BeforeInsert, operator.BeforeUpsert, operator.BeforeReplace:
		return true
	}
	return false
}

// Do calls validator check
// Don't use opts here
func Do(doc interface{}, opType operator.OpType, opts ...interface{}) error {
	if !validatorNeeded(opType) {
		return nil
	}
	to := reflect.TypeOf(doc)
	if to == nil {
		return nil
	}
	switch reflect.TypeOf(doc).Kind() {
	case reflect.Slice:
		return sliceHandle(doc, opType)
	case reflect.Ptr:
		v := reflect.ValueOf(doc).Elem()
		switch v.Kind() {
		case reflect.Slice:
			return sliceHandle(v.Interface(), opType)
		default:
			return do(doc)
		}
	}
	//fmt.Println("not support type")
	return nil
}

// sliceHandle handles the slice docs
func sliceHandle(docs interface{}, opType operator.OpType) error {
	// []interface{}{UserType{}...}
	if h, ok := docs.([]interface{}); ok {
		for _, v := range h {
			if reflect.TypeOf(v).Kind() == reflect.Slice {
				if err := Do(v, opType); err != nil {
					return err
				}
			} else {
				if err := do(v); err != nil {
					return err
				}
			}
		}
		return nil
	}
	// []UserType{}
	s := reflect.ValueOf(docs)
	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Kind() == reflect.Slice {
			if err := Do(s.Index(i).Interface(), opType); err != nil {
				return err
			}
		} else {
			if err := do(s.Index(i).Interface()); err != nil {

				return err
			}
		}
	}
	return nil
}

// do check if opType is supported and call fieldHandler
func do(doc interface{}) error {
	//fmt.Println(reflect.TypeOf(doc).Name(), reflect.TypeOf(doc).Kind())
	return validate.Struct(doc)
}
