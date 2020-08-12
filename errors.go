package qmgo

import (
	"errors"
	"strings"

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

// IsErrNoDocuments check if err is no documents, both mongo-go-driver error and qmgo custom error
func IsErrNoDocuments(err error) bool {
	if err == mongo.ErrNoDocuments || err == ErrNoSuchDocuments {
		return true
	}
	return false
}

// IsDup check if err is mongo E11000 (duplicate err)。
func IsDup(err error) bool {
	return strings.Contains(err.Error(), "E11000")
}
