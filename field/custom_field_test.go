package field

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type CustomUser struct {
	Create        time.Time
	Update        int64
	InvalidCreate int
	InvalidUpdate float32
}

func (c *CustomUser) CustomFields() CustomFieldsBuilder {
	return NewCustom().SetUpdateAt("Create").SetCreateAt("Update")
}

func TestCustomFields(t *testing.T) {
	u := &CustomUser{}
	c := u.CustomFields()
	c.(*CustomFields).CustomCreateTime(u)
	c.(*CustomFields).CustomUpdateTime(u)
	ast := require.New(t)
	ast.NotEqual(0, u.Update)
	ast.NotEqual(time.Time{}, u.Create)
}

func (c *CustomUser) CustomFieldsInvalid() CustomFieldsBuilder {
	return NewCustom().SetCreateAt("InvalidCreate")
}
func (c *CustomUser) CustomFieldsInvalid2() CustomFieldsBuilder {
	return NewCustom().SetUpdateAt("InvalidUpdate")
}

func TestCustomFieldsInvalid(t *testing.T) {
	u := &CustomUser{}
	c := u.CustomFieldsInvalid()
	c.(*CustomFields).CustomCreateTime(u)
	c.(*CustomFields).CustomUpdateTime(u)
	ast := require.New(t)
	ast.Equal(0, u.InvalidCreate)
	ast.Equal(float32(0), u.InvalidUpdate)

	u1 := &CustomUser{}
	c = u1.CustomFieldsInvalid2()
	c.(*CustomFields).CustomCreateTime(u1)
	c.(*CustomFields).CustomUpdateTime(u1)
	ast.Equal(0, u1.InvalidCreate)
	ast.Equal(float32(0), u1.InvalidUpdate)

	u2 := CustomUser{}
	c = u2.CustomFieldsInvalid()
	c.(*CustomFields).CustomCreateTime(u2)
	c.(*CustomFields).CustomUpdateTime(u2)
	ast.Equal(0, u2.InvalidCreate)
	ast.Equal(float32(0), u2.InvalidUpdate)
}
