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
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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
