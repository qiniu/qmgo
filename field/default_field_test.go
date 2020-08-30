package field

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDefaultField(t *testing.T) {
	ast := require.New(t)

	df := &DefaultField{}
	df.DefaultCreateAt()
	df.DefaultUpdateAt()
	ast.NotEqual(time.Time{}, df.UpdateAt)
	ast.NotEqual(time.Time{}, df.CreateAt)
}
