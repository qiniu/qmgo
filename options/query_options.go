package options

import (
	"github.com/qiniu/qmgo/hook"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FindOptions struct {
	QueryHook                hook.QueryHook
	FindOneAndDeleteOptions  options.FindOneAndDeleteOptions
	FindOneAndUpdateOptions  options.FindOneAndUpdateOptions
	FindOneAndReplaceOptions options.FindOneAndReplaceOptions
}
