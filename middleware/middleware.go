package middleware

import (
	"github.com/qiniu/qmgo/field"
	"github.com/qiniu/qmgo/hook"
	"github.com/qiniu/qmgo/operator"
)

// callback define the callback function type
type callback func(doc interface{}, opType operator.OpType, opts ...interface{}) error

// middlewareCallback the register callback slice, the hook and field callback register when init for the order
var middlewareCallback = []callback{
	hook.Do,
	field.Do,
}

// Register register callback into middleware,
func Register(cb callback) {
	middlewareCallback = append(middlewareCallback, cb)
}

// Do call every registers
// The doc is always the document to operate
func Do(doc interface{}, opType operator.OpType, opts ...interface{}) error {
	for _, cb := range middlewareCallback {
		if err := cb(doc, opType, opts...); err != nil {
			return err
		}
	}
	return nil
}
