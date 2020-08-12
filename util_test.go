package qmgo

import (
	"fmt"
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
