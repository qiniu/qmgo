package qmgo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDatabase(t *testing.T) {
	ast := require.New(t)

	var sTimeout int64 = 500000
	var cTimeout int64 = 3000
	var maxPoolSize uint64 = 3000
	var minPoolSize uint64 = 0
	collName := "testopen"
	dbName := "qmgotest"

	cfg := Config{
		Uri:              "mongodb://localhost:27017",
		Database:         dbName,
		Coll:             collName,
		ConnectTimeoutMS: &cTimeout,
		SocketTimeoutMS:  &sTimeout,
		MaxPoolSize:      &maxPoolSize,
		MinPoolSize:      &minPoolSize,
	}

	c, err := NewClient(context.Background(), &cfg)
	ast.NoError(err)
	cli := c.Database(cfg.Database)
	ast.Nil(err)
	ast.Equal(dbName, cli.GetDatabaseName())
	coll := cli.Collection(collName)
	ast.Equal(collName, coll.GetCollectionName())
	cli.Collection(collName).DropCollection(context.Background())
	cli.DropDatabase(context.Background())

}
