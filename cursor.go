package mongox

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// Cursor
type Cursor struct {
	ctx    context.Context
	cursor *mongo.Cursor
}

// Next
// 获取游标下一条文档
func (c *Cursor) Next(result interface{}) bool {
	var err error

	if c.cursor.Next(c.ctx) {
		err = c.cursor.Decode(result)
		if err == nil {
			return true
		}
	}

	return false
}

// Close
// 关闭 Cursor，关闭后，Next 不可执行
// 当游标对象不再使用时，应主动关闭
func (c *Cursor) Close() error {
	return c.cursor.Close(c.ctx)
}

// Err
// 返回 Cursor 的最后一次 error，若无错误发生，则返回nil
func (c *Cursor) Err() error {
	return c.cursor.Err()
}
