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
	ast.NoError(Do(ctx, "success", operator.BeforeInsert))

	// valid register
	Register(callbackTest)
	ast.NoError(Do(ctx, "success", operator.BeforeInsert))
	ast.Error(Do(ctx, "failure", operator.BeforeUpsert))
	ast.NoError(Do(ctx, "failure", operator.BeforeUpdate, "success"))
}

func callbackTest(ctx context.Context, doc interface{}, opType operator.OpType, opts ...interface{}) error {
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
