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
	SetUpdateAt(filedName string) CustomFieldsBuilder
	SetCreateAt(filedName string) CustomFieldsBuilder
	SetId(filedName string) CustomFieldsBuilder
}

// NewCustom creates new Builder which is used to set the custom fields
func NewCustom() CustomFieldsBuilder {
	return &CustomFields{}
}

// SetUpdateAt set the custom UpdateAt field
func (c *CustomFields) SetUpdateAt(filedName string) CustomFieldsBuilder {
	c.updateAt = filedName
	return c
}

// SetCreateAt set the custom CreateAt field
func (c *CustomFields) SetCreateAt(filedName string) CustomFieldsBuilder {
	c.createAt = filedName
	return c
}

// SetId set the custom Id field
func (c *CustomFields) SetId(filedName string) CustomFieldsBuilder {
	c.id = filedName
	return c
}

// CustomCreateTime changes the custom create time
func (c CustomFields) CustomCreateTime(doc interface{}) CustomFields {
	if c.createAt == "" {
		return c
	}
	fieldName := c.createAt
	setTime(doc, fieldName)
	return c
}

// CustomUpdateTime changes the custom update time
func (c CustomFields) CustomUpdateTime(doc interface{}) CustomFields {
	if c.updateAt == "" {
		return c
	}
	fieldName := c.updateAt
	setTime(doc, fieldName)
	return c
}

// CustomUpdateTime changes the custom update time
func (c CustomFields) CustomId(doc interface{}) CustomFields {
	if c.id == "" {
		return c
	}
	fieldName := c.id
	setId(doc, fieldName)
	return c
}

// setTime changes the custom time fields
func setTime(doc interface{}, fieldName string) {
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
			ca.Set(reflect.ValueOf(tt))
		case int64:
			ca.SetInt(tt.Unix())
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
			ca.Set(reflect.ValueOf(primitive.NewObjectID()))
		case string:
			ca.SetString(primitive.NewObjectID().String())
		default:
			fmt.Println("unsupported type to setId", a)
		}
	}
}
