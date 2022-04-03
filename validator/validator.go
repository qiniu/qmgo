package validator

import (
	"context"
	"reflect"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/qiniu/qmgo/operator"
)

// use a single instance of Validate, it caches struct info
var validate = validator.New()

// SetValidate let validate can use custom rules
func SetValidate(v *validator.Validate) {
	validate = v
}

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
func Do(ctx context.Context, doc interface{}, opType operator.OpType, opts ...interface{}) error {
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
	default:
		return do(doc)
	}
}

// sliceHandle handles the slice docs
func sliceHandle(docs interface{}, opType operator.OpType) error {
	// []interface{}{UserType{}...}
	if h, ok := docs.([]interface{}); ok {
		for _, v := range h {
			if err := do(v); err != nil {
				return err
			}
		}
		return nil
	}
	// []UserType{}
	s := reflect.ValueOf(docs)
	for i := 0; i < s.Len(); i++ {
		if err := do(s.Index(i).Interface()); err != nil {

			return err
		}
	}
	return nil
}

// do check if opType is supported and call fieldHandler
func do(doc interface{}) error {
	if !validatorStruct(doc) {
		return nil
	}
	return validate.Struct(doc)
}

// validatorStruct check if kind of doc is validator supported struct
// same implement as validator
func validatorStruct(doc interface{}) bool {
	val := reflect.ValueOf(doc)
	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct || val.Type() == reflect.TypeOf(time.Time{}) {
		return false
	}
	return true
}
