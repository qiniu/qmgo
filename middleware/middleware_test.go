package middleware

import (
	"context"
	"fmt"
	"github.com/qiniu/qmgo/operator"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMiddleware(t *testing.T) {
	ast := require.New(t)
	ctx := context.Background()
	// not register
	ast.NoError(Do("success", operator.BeforeInsert, ctx))

	// valid register
	Register(callbackTest)
	ast.NoError(Do("success", operator.BeforeInsert, ctx))
	ast.Error(Do("failure", operator.BeforeUpsert, ctx))
	ast.NoError(Do("failure", operator.BeforeUpdate, ctx, "success"))
}

func callbackTest(doc interface{}, opType operator.OpType, ctx context.Context, opts ...interface{}) error {
	if doc.(string) == "success" && opType == operator.BeforeInsert {
		return nil
	}
	if len(opts) > 0 && opts[0].(string) == "success" {
		return nil
	}
	if doc.(string) == "failure" && opType == operator.BeforeUpsert {
		return fmt.Errorf("this is error")
	}
	return nil
}
