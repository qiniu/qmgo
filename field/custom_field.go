package field

import (
	"fmt"
	"reflect"
	"time"
)

// CustomFields defines struct of supported custom fields
type CustomFields struct {
	createAt string
	updateAt string
}

// CustomFieldsHook defines the interface, CustomFields return custom field user want to change
type CustomFieldsHook interface {
	CustomFields() CustomFieldsBuilder
}

// CustomFieldsBuilder defines the interface which user use to set custom fields
type CustomFieldsBuilder interface {
	SetUpdateAt(FiledName string) CustomFieldsBuilder
	SetCreateAt(FiledName string) CustomFieldsBuilder
}

// New creates new Builder which is used to set the custom fields
func NewCustom() CustomFieldsBuilder {
	return &CustomFields{}
}

// SetUpdateAt set the custom UpdateAt field
func (c *CustomFields) SetUpdateAt(filedName string) CustomFieldsBuilder {
	c.updateAt = filedName
	return c
}

// SetCreateAt set the custom UpdateAt field
func (c *CustomFields) SetCreateAt(filedName string) CustomFieldsBuilder {
	c.createAt = filedName
	return c
}

// CustomCreateTime changes the custom create time
func (c CustomFields) CustomCreateTime(doc interface{}) {
	if c.createAt == "" {
		return
	}
	fieldName := c.createAt
	updateTime(doc, fieldName)
}

// CustomUpdateTime changes the custom update time
func (c CustomFields) CustomUpdateTime(doc interface{}) {
	if c.updateAt == "" {
		return
	}
	fieldName := c.updateAt
	updateTime(doc, fieldName)
}

// updateTime changes the time fields
func updateTime(doc interface{}, fieldName string) {
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
			fmt.Println("unsupported type to set", a)
		}
	}
}
