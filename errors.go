package qmgo

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrQueryNotSlicePointer        = errors.New("result argument must be a pointer to a slice")
	ErrQueryNotSliceType           = errors.New("result argument must be a slice address")
	ErrQueryResultTypeInconsistent = errors.New("result type is not equal mongodb value type")
	ErrQueryResultValCanNotChange  = errors.New("the value of result can not be changed")
	ErrNoSuchDocuments             = errors.New(mongo.ErrNoDocuments.Error())
)
