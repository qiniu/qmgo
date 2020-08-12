package qmgo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

// Now return Millisecond current time
func Now() time.Time {
	return time.Unix(0, time.Now().UnixNano()/1e6*1e6)
}

// NewObjectID generates a new ObjectID.
// Watch out: the way it generates objectID is different from mgo
func NewObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

// SplitSortField handle sort symbol: "+"/"-" in front of field
// if "+"， return sort as 1
// if "-"， return sort as -1
func SplitSortField(field string) (key string, sort int32) {
	sort = 1
	key = field

	if len(field) != 0 {
		switch field[0] {
		case '+':
			key = strings.TrimPrefix(field, "+")
			sort = 1
		case '-':
			key = strings.TrimPrefix(field, "-")
			sort = -1
		}
	}

	return key, sort
}
