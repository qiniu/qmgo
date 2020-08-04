package qmgo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// Cursor
type Cursor struct {
	ctx    context.Context
	cursor *mongo.Cursor
}

// Next gets the next document for this cursor. It returns true if there were no errors and the cursor has not been
// exhausted.
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

// Close closes this cursor. Next and TryNext must not be called after Close has been called.
// When the cursor object is no longer in use, it should be actively closed
func (c *Cursor) Close() error {
	return c.cursor.Close(c.ctx)
}

// Err return the last error of Cursor, if no error occurs, return nil
func (c *Cursor) Err() error {
	return c.cursor.Err()
}
