package middleware

import (
	"github.com/qiniu/qmgo/field"
	"github.com/qiniu/qmgo/hook"
	"github.com/qiniu/qmgo/operator"
)

type callback func(doc interface{}, opType operator.OpType, opts ...interface{}) error

var middlewareCallback = []callback{
	hook.Do,
	field.Do,
}

// Register register callback into middleware, field and hook callback register when init for the order
func Register(cb callback) {
	middlewareCallback = append(middlewareCallback, cb)
}

// Do call every registers
// The doc is always the document to operate
// The opts maybe valid, e.g. the hook pass by opts which doesn't implement by document struct
func Do(doc interface{}, opType operator.OpType, opts ...interface{}) error {
	for _, cb := range middlewareCallback {
		if err := cb(doc, opType, opts...); err != nil {
			return err
		}
	}
	return nil
}
