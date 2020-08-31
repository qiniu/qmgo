package field

import (
	"time"
)

// DefaultFieldHook defines the interface to change default fields by hook
type DefaultFieldHook interface {
	DefaultUpdateAt()
	DefaultCreateAt()
}

// DefaultField defines the default fields to handle when operation happens
// import the DefaultField in document struct to make it working
type DefaultField struct {
	CreateAt time.Time `bson:"createAt"`
	UpdateAt time.Time `bson:"updateAt"`
}

// DefaultUpdateAt changes the default updateAt field
func (df *DefaultField) DefaultUpdateAt() {
	df.UpdateAt = time.Now().Local()
}

// DefaultCreateAt changes the default createAt field
func (df *DefaultField) DefaultCreateAt() {
	df.CreateAt = time.Now().Local()
}
