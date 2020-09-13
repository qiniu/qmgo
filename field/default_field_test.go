package field

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDefaultField(t *testing.T) {
	ast := require.New(t)

	df := &DefaultField{}
	df.DefaultCreateAt()
	df.DefaultUpdateAt()
	df.DefaultId()
	ast.NotEqual(time.Time{}, df.UpdateAt)
	ast.NotEqual(time.Time{}, df.CreateAt)
	ast.NotEqual(primitive.NilObjectID, df.Id)
}
