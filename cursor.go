package qmgo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

// Cursor struct define
type Cursor struct {
	ctx    context.Context
	cursor *mongo.Cursor
	err    error
}

// Next gets the next document for this cursor. It returns true if there were no errors and the cursor has not been
// exhausted.
func (c *Cursor) Next(result interface{}) bool {
	if c.err != nil {
		return false
	}
	var err error
	if c.cursor.Next(c.ctx) {
		err = c.cursor.Decode(result)
		if err == nil {
			return true
		}
	}
	return false
}

// All iterates the cursor and decodes each document into results. The results parameter must be a pointer to a slice.
// recommend to use All() in struct Query or Aggregate
func (c *Cursor) All(results interface{}) error {
	if c.err != nil {
		return c.err
	}
	return c.cursor.All(c.ctx, results)
}

// ID returns the ID of this cursor, or 0 if the cursor has been closed or exhausted.
//func (c *Cursor) ID() int64 {
//	if c.err != nil {
//		return 0
//	}
//	return c.cursor.ID()
//}

// Close closes this cursor. Next and TryNext must not be called after Close has been called.
// When the cursor object is no longer in use, it should be actively closed
func (c *Cursor) Close() error {
	if c.err != nil {
		return c.err
	}
	return c.cursor.Close(c.ctx)
}

// Err return the last error of Cursor, if no error occurs, return nil
func (c *Cursor) Err() error {
	if c.err != nil {
		return c.err
	}
	return c.cursor.Err()
}
