package qmgo

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNow(t *testing.T) {
	t1 := time.Unix(0, time.Now().UnixNano()/1e6*1e6)
	t2 := Now()
	fmt.Println(t1, t2)
}

func TestNewObjectID(t *testing.T) {
	objId := NewObjectID()
	objId.Hex()
}

func TestCompareVersions(t *testing.T) {
	ast := require.New(t)
	i, err := CompareVersions("4.4.0", "3.0")
	ast.NoError(err)
	ast.True(i > 0)
	i, err = CompareVersions("3.0.1", "3.0")
	ast.NoError(err)
	ast.True(i == 0)
	i, err = CompareVersions("3.1.5", "4.0")
	ast.NoError(err)
	ast.True(i < 0)
}
