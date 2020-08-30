package options

import "github.com/qiniu/qmgo/hook"

type InsertOneOptions struct {
	InsertHook hook.InsertHook
}
type InsertManyOptions struct {
	InsertHook interface{}
}
