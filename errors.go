package qmgo

import "errors"

var (
	ERR_QUERY_NOT_SLICE_POINTER         = errors.New("result argument must be a pointer to a slice")
	ERR_QUERY_NOT_SLICE_TYPE            = errors.New("result argument must be a slice address")
	ERR_QUERY_RESULT_TYPE_INCONSISTEN   = errors.New("result type is not equal mongodb value type")
	ERR_QUERY_RESULT_VAL_CAN_NOT_CHANGE = errors.New("the value of result can not be changed")
	ERR_NO_SUCH_RECORD                  = errors.New("no such record")
)
