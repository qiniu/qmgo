/*
 Copyright 2020 The Qmgo Authors.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

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
