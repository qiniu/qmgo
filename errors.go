package qmgo

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	// ErrQueryNotSlicePointer return if result argument is not a pointer to a slice
	ErrQueryNotSlicePointer = errors.New("result argument must be a pointer to a slice")
	// ErrQueryNotSliceType return if result argument is not slice address
	ErrQueryNotSliceType = errors.New("result argument must be a slice address")
	// ErrQueryResultTypeInconsistent return if result type is not equal mongodb value type
	ErrQueryResultTypeInconsistent = errors.New("result type is not equal mongodb value type")
	// ErrQueryResultValCanNotChange return if the value of result can not be changed
	ErrQueryResultValCanNotChange = errors.New("the value of result can not be changed")
	// ErrNoSuchDocuments return if no document found
	ErrNoSuchDocuments = errors.New(mongo.ErrNoDocuments.Error())
)
